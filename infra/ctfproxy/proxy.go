package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/pprof"
	"net/url"
	"runtime"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/adamyi/CTFProxy/third_party/eddsa"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

var _limiter *limiter.Limiter
var zmux *http.ServeMux

func initRateLimit() {
	rate, err := limiter.NewRateFromFormatted(_configuration.RateLimit)
	if err != nil {
		panic(err)
	}

	store := memory.NewStore()
	_limiter = limiter.New(store, rate, limiter.WithTrustForwardHeader(true))
}

func init() {
	zmux = http.NewServeMux()
	zmux.HandleFunc("/ipz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%v", r.RemoteAddr)
	})
	zmux.HandleFunc("/testpanic", func(w http.ResponseWriter, r *http.Request) {
		panic(r)
	})
	zmux.HandleFunc("/statusz", StatusHandler)
	zmux.HandleFunc("/debug/gc", func(w http.ResponseWriter, r *http.Request) {
		runtime.GC()
		w.Write([]byte("done"))
	})
	zmux.HandleFunc("/debug/pprof/", pprof.Index)
	zmux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	zmux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	zmux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	zmux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	zmux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		NewCPError(http.StatusNotFound, "404 Not Found", "The requested URL does not exist.", "", "ZDomain 404").Write(w)
	})
}

func handleUP(rsp http.ResponseWriter, req *http.Request) {
	initUPRsp(rsp)

	if !_configuration.RawServicesMap[req.Host] {
		rpath := cleanPath(req.URL.Path)
		if req.URL.Path != rpath {
			url := *req.URL
			url.Path = rpath
			http.Redirect(rsp, req, rpath, http.StatusTemporaryRedirect)
			return
		}
	}

	if req.Host == _configuration.LoginDomain {
		handleLogin(rsp, req)
		return
	}

	if req.Host == _configuration.CorpDomain || req.Host == "www."+_configuration.CorpDomain {
		req.Host = _configuration.RootServiceDomain
	}

	username, displayname, err := getUsername(req)
	if err != nil {
		NewCPError(http.StatusBadRequest, "You issued a malformed request", "failed to get username", "", err.Error()).Write(rsp)
		return
	}
	remoteIP := strings.Split(req.RemoteAddr, ":")[0]
	rsp.Header().Add("X-CTFProxy-I-User", username)

	// only rate-limit external users for now...
	// change this if something is causing issues
	if !strings.HasSuffix(username, "@services."+_configuration.CorpDomain) {
		limit, errr := _limiter.Get(req.Context(), remoteIP+"_"+username)
		if errr != nil {
			NewCPError(http.StatusInternalServerError, "Internal Server Error", "Please notify course staff", "", errr.Error()).Write(rsp)
			return
		}

		if limit.Reached {
			NewCPError(420, "420 Enhance Your Calm", "Calm down!!! You're sending requests too fast", "", "rate limited").Write(rsp)
			return
		}
	}

	groups := getGroups(username)
	dbg := false
	for _, group := range groups {
		rsp.Header().Add("X-CTFProxy-I-Group", group)
		if group == "ctfproxy-debug@groups."+_configuration.CorpDomain {
			dbg = true
			rsp.Header().Add("X-CTFProxy-I-Debug", "1")
		}
	}

	if req.Host == _configuration.ZDomain {
		if req.URL.Path == "/healthz" {
			HealthCheckHandler(rsp, req)
			return
		}
		if !dbg {
			NewCPError(http.StatusForbidden, "403 Forbidden", "This is a debug-only URL", "", "").Write(rsp)
			return
		}
		zmux.ServeHTTP(rsp, req)
		return
	}

	ctx, levelShift, err := getNetworkContext(req, username)
	if err != nil {
		// fmt.Println("getNetowrkContext - ", levelShift, err)
		NewCPError(http.StatusBadRequest, "Could not resolve the IP address for host "+req.Host, "Your client has issued a malformed or illegal request.", "", "getNetworkContext: "+err.Error()).Write(rsp)
		return
	}
	if ra := ctx.Value("up_real_addr"); ra != nil {
		rsp.Header().Add("X-CTFProxy-I-Real-Addr", ra.(string))
	}

	full_url := req.Host + req.RequestURI
	// fmt.Println("getNetworkContext", levelShift, full_url)

	servicename, err := getServiceNameFromDomain(req.Host)
	if err != nil {
		NewCPError(http.StatusBadRequest, "Could not resolve the IP address for host "+req.Host, "Your client has issued a malformed or illegal request.", "", "getServiceNameFromDomain: "+err.Error()).Write(rsp)
		return
	}

	ae := hasAccess(servicename, username, groups, req)
	if ae != nil {
		if ae.Code == http.StatusForbidden && username == "anonymous@anonymous."+_configuration.CorpDomain {
			http.Redirect(rsp, req, "https://"+_configuration.LoginDomain+"/?return_url="+url.QueryEscape("https://"+full_url), http.StatusTemporaryRedirect)
		} else {
			ae.Write(rsp)
		}
		return
	}

	if levelShift {
		servicename = "ctfproxy"
	}

	servicename += "@services." + _configuration.CorpDomain

	ptstr := ""
	expirationTime := time.Now().Add(5 * time.Minute)
	pclaims := Claims{
		Username:    username,
		Displayname: displayname,
		Service:     servicename,
		Groups:      groups,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	ptoken := jwt.NewWithClaims(eddsa.SigningMethodEdDSA, pclaims)
	ptstr, err = ptoken.SignedString(_configuration.SignKey)
	if err != nil {
		NewCPError(http.StatusInternalServerError, "Internal Server Error", "Something went wrong while generating JWT", "", err.Error()).Write(rsp)
		return
	}

	if req.URL.Path == "/ws" {
		handleWs(ctx, rsp, req, ptstr, levelShift)
		return
	}

	bodyBytes, _ := ioutil.ReadAll(req.Body)
	// req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	scheme := "http://"
	if levelShift {
		scheme = "https://"
	}

	preq, err := http.NewRequestWithContext(ctx, req.Method, scheme+full_url, bytes.NewReader(bodyBytes))
	if err != nil {
		NewCPError(http.StatusBadGateway, "Bad Gateway", "Something went wrong connecting to internal service", "", "http.NewRequestWithContext ("+req.Method+" "+scheme+full_url+"): "+err.Error()).Write(rsp)
		return
	}

	if !_configuration.RawServicesMap[req.Host] {
		for name, value := range req.Header {
			val := value[0]
			ln := strings.ToLower(name)
			if ln == "cookie" {
				cookies := strings.Split(val, ";")
				l := len(cookies)
				for i, cookie := range cookies {
					if strings.HasPrefix(strings.TrimLeft(strings.ToLower(cookie), " "), _configuration.AuthCookieKey) {
						cookies[i] = cookies[l-1]
						l -= 1
					}
				}
				if l > 0 {
					val = strings.TrimLeft(strings.Join(cookies[:l], ";"), " ")
				} else {
					val = ""
				}
			} else if ln == strings.ToLower(_configuration.SubAccHeader) || ln == strings.ToLower(_configuration.ImpersonateTokenHeader) {
				continue
			}
			if val != "" {
				preq.Header.Set(name, val)
			}
		}

		preq.Header.Set(_configuration.InternalJWTHeader, ptstr)
		if levelShift {
			preq.Header.Set("X-CTFProxy-LevelShift", "1")
		}
		preq.Header.Set("X-CTFProxy-Remote-Addr", remoteIP)

	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	presp, err := client.Do(preq)
	if err != nil {
		NewCPError(http.StatusBadGateway, "Bad Gateway", "Something went wrong connecting to internal service", "", "client.Do("+req.Method+" "+scheme+full_url+"): "+err.Error()).Write(rsp)
		return
	}
	if presp.Header.Get("Content-Type") == "ctfproxy/error" {
		handleUpstreamCPError(rsp, presp)
		return
	}
	copyResponse(rsp, presp)
}
