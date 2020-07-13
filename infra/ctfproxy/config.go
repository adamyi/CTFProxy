package main

import (
	"context"
	"crypto/ed25519"
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/adamyi/hotconfig"
	"github.com/adamyi/CTFProxy/third_party/eddsa"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

type flagArray []string

func (i *flagArray) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func (i *flagArray) String() string {
	return "this is a flag array"
}

type Configuration struct {
	ListenAddress          string
	SSLListenAddress       string
	SSLCertificates        flagArray
	SSLPrivateKeys         flagArray
	RawServices            flagArray
	RawServicesMap         map[string]bool
	MTLSCA                 string
	SignKey                *ed25519.PrivateKey
	VerifyKey              *ed25519.PublicKey
	CorpDomain             string
	ResolvingDomain        string
	ServiceName            string
	GAIAEndpoint           string
	LoginDomain            string
	ZDomain                string
	CLIRelayDomain         string
	RootServiceDomain      string
	AuthCookieKey          string
	SubAccHeader           string
	ImpersonateTokenHeader string
	InternalJWTHeader      string
	GCPProject             string
	CertBucket             string
	RateLimit              string
	SubAccLengthLimit      int
	ConfigFetchInterval    time.Duration
	AccessPolicies         *hotconfig.Config
	ServiceAliases         *hotconfig.Config
}

var _configuration = Configuration{}

var gcsClient *storage.Client

func newHotConfig(ctx context.Context, path string) (c *hotconfig.Config, err error) {
	if strings.HasPrefix(path, "gs://") {
		c, err = hotconfig.NewGCSConfig(ctx, gcsClient, path, hotconfig.DecoderFunc(
			func(r io.Reader) (interface{}, error) {
				var v map[string]string
				err := json.NewDecoder(r).Decode(&v)
				if err != nil {
					return nil, err
				}
				return v, nil
			}))
		if err != nil {
			return nil, err
		}
		go c.StartPeriodicUpdate(ctx, 5*time.Minute)
		return
	}
	ap, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var v map[string]string
	err = json.Unmarshal(ap, &v)
	if err != nil {
		return nil, err
	}
	c = hotconfig.NewConfig(ctx, hotconfig.FetcherFunc(
		func(ctx context.Context) (interface{}, error) {
			return v, nil
		}))
	return
}

func readConfig() {
	var err error
	ctx := context.Background()
	var publicKeyPath, privateKeyPath, accessPath, aliasPath, gcpcredPath, fetchInterval string
	flag.StringVar(&_configuration.ListenAddress, "listen", "0.0.0.0:80", "http listen address")
	flag.StringVar(&_configuration.SSLListenAddress, "ssl_listen", "0.0.0.0:443", "https listen address")
	flag.Var(&_configuration.SSLCertificates, "ssl_cert", "HTTPS certificate")
	flag.Var(&_configuration.SSLPrivateKeys, "ssl_key", "HTTPS certificate private key")
	flag.Var(&_configuration.RawServices, "raw_services", "Internal services to disable CTFProxy headers and URL cleaning")
	flag.StringVar(&_configuration.MTLSCA, "mtls_ca", "", "Path to MTLS Certificate Authority")
	flag.StringVar(&publicKeyPath, "jwt_public_key", "", "Path to JWT public key")
	flag.StringVar(&privateKeyPath, "jwt_private_key", "", "Path to JWT private key")
	flag.StringVar(&accessPath, "access_config_path", "", "Path to access policy config")
	flag.StringVar(&aliasPath, "alias_config_path", "", "Path to alias config")
	flag.StringVar(&_configuration.CorpDomain, "corp_domain", "ctfproxy.lhdev.xyz", "Corp Domain Name")
	flag.StringVar(&_configuration.LoginDomain, "login_subdomain", "login.corp", "Login subdomain")
	flag.StringVar(&_configuration.ZDomain, "zdomain", "ctfproxyz", "Z-Page domain")
	flag.StringVar(&_configuration.CLIRelayDomain, "clirelay_subdomain", "clirelay.corp", "CLI-Relay subdomain")
	flag.StringVar(&_configuration.RootServiceDomain, "root_service_domain", "search.corp", "Root Corp Domain redirecting service subdomain")
	flag.StringVar(&_configuration.ServiceName, "service_name", "ctfproxy", "Service Account Name for CTFProxy (useful for differentiating while running multiple instances, e.g. autopush, test, staging, preprod, prod)")
	// flag.StringVar(&_configuration.GAIAEndpoint, "gaia_api_endpoint", "http://gaia.corp.ctfproxy.lhdev.xyz", "GAIA API URL")
	flag.StringVar(&_configuration.AuthCookieKey, "auth_cookie_key", "ctfproxy_auth", "auth cookies name")
	flag.StringVar(&_configuration.SubAccHeader, "subacc_header", "X-CTFProxy-SubAcc", "http header for subaccount")
	flag.StringVar(&_configuration.ImpersonateTokenHeader, "impersonate_header", "X-CTFProxy-SubAcc-JWT", "http header for subaccount's jwt")
	flag.StringVar(&_configuration.InternalJWTHeader, "jwt_header", "X-CTFProxy-JWT", "http header for JWT sent to internal services")
	flag.StringVar(&_configuration.CertBucket, "cert_gcs_bucket", "", "GCS bucket for autocert cache")
	flag.StringVar(&_configuration.GCPProject, "gcp_project", "", "GCP Project")
	flag.StringVar(&_configuration.RateLimit, "rate_limit", "1000-M", "rate limit")
	flag.IntVar(&_configuration.SubAccLengthLimit, "subacc_length_limit", 10, "length limit for subacc (-1) to disable")
	flag.StringVar(&gcpcredPath, "gcp_service_account", "", "GCP service account json")
	flag.StringVar(&fetchInterval, "config_fetch_interval", "1m", "Config fetch interval")
	flag.Parse()
	flag.VisitAll(func(f *flag.Flag) {
		log.Printf("Runnng with flag %s: %s\n", f.Name, f.Value)
	})
	_configuration.ConfigFetchInterval, err = time.ParseDuration(fetchInterval)
	if err != nil {
		panic(err)
	}
	_configuration.ZDomain += "." + _configuration.CorpDomain
	_configuration.LoginDomain += "." + _configuration.CorpDomain
	_configuration.CLIRelayDomain += "." + _configuration.CorpDomain
	_configuration.RootServiceDomain += "." + _configuration.CorpDomain
	if len(_configuration.SSLCertificates) != len(_configuration.SSLPrivateKeys) {
		panic("len(ssl_cert) != len(ssl_key), check params")
	}
	gcpcreds, err := ioutil.ReadFile(gcpcredPath)
	if err != nil {
		panic(err)
	}
	gcscreds, err := google.CredentialsFromJSON(ctx, gcpcreds, storage.ScopeReadWrite)
	if err != nil {
		panic(err)
	}
	gcsClient, err = storage.NewClient(ctx, option.WithCredentials(gcscreds))
	if err != nil {
		panic(err)
	}
	JwtPubKey, err := ioutil.ReadFile(publicKeyPath)
	if err != nil {
		panic(err)
	}
	JwtKey, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		panic(err)
	}
	_configuration.SignKey, err = eddsa.ParseEdPrivateKeyFromPEM(JwtKey)
	if err != nil {
		panic(err)
	}
	_configuration.VerifyKey, err = eddsa.ParseEdPublicKeyFromPEM(JwtPubKey)
	if err != nil {
		panic(err)
	}
	_configuration.AccessPolicies, err = newHotConfig(ctx, accessPath)
	if err != nil {
		panic(err)
	}
	_configuration.ServiceAliases, err = newHotConfig(ctx, aliasPath)
	if err != nil {
		panic(err)
	}
	_configuration.ResolvingDomain = _configuration.CorpDomain
	if os.Getenv("CTFPROXY_CLUSTER") == "k8s" {
		_configuration.ResolvingDomain = "default.svc.cluster.local"
	}
	_configuration.GAIAEndpoint = "http://gaia." + _configuration.ResolvingDomain
	_configuration.RawServicesMap = make(map[string]bool)
	for _, service := range _configuration.RawServices {
		_configuration.RawServicesMap[service+"."+_configuration.CorpDomain] = true
	}
}
