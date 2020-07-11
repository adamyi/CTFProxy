package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/adamyi/CTFProxy/third_party/autocertcache"

	"golang.org/x/crypto/acme/autocert"
)

func redirectSSL(rsp http.ResponseWriter, req *http.Request) {
	target := "https://" + req.Host + req.URL.Path
	if len(req.URL.RawQuery) > 0 {
		target += "?" + req.URL.RawQuery
	}
	http.Redirect(rsp, req, target,
		http.StatusTemporaryRedirect)
}

var authCA *x509.CertPool
var certManager *autocert.Manager

func buildSSLServer() *http.Server {
	authCA = x509.NewCertPool()

	cacert, err := ioutil.ReadFile(_configuration.MTLSCA)
	cfg := &tls.Config{}
	if err != nil && len(cacert) > 0 {
		log.Println("Error reading mTLS CA or CA file empty. mTLS Auth disabled.")
		log.Println(err)
	} else {
		log.Println("mTLS auth enabled")
		authCA.AppendCertsFromPEM(cacert)
		cfg = &tls.Config{
			ClientAuth: tls.RequestClientCert,
			ClientCAs:  authCA,
		}
	}

	cfg.Certificates = make([]tls.Certificate, len(_configuration.SSLCertificates))
	for i, cert := range _configuration.SSLCertificates {
		cfg.Certificates[i], err = tls.LoadX509KeyPair(cert, _configuration.SSLPrivateKeys[i])
		if err != nil {
			log.Fatal(err)
		}
	}

	// i'm lazy
	cfg.BuildNameToCertificate()
	ntc := cfg.NameToCertificate
	cfg.NameToCertificate = nil

	hostPolicy := func(ctx context.Context, host string) error {
		sn, err := getServiceNameFromDomain(host)
		if err != nil {
			return err
		}
		if _, ok := _configuration.AccessPolicies.ConfigOrNil().(map[string]string)[sn]; ok {
			return nil
		}
		return errors.New("non-existing host")
	}

	certManager = &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: hostPolicy,
		Cache:      autocertcache.NewGoogleCloudStorageCache(gcsClient, _configuration.CertBucket),
	}

	cfg.GetCertificate = func(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
		name := strings.ToLower(hello.ServerName)
		if cert, ok := ntc[name]; ok {
			return cert, nil
		}
		if len(name) > 0 {
			labels := strings.Split(name, ".")
			labels[0] = "*"
			wildcardName := strings.Join(labels, ".")
			if cert, ok := ntc[wildcardName]; ok {
				return cert, nil
			}
		}
		return certManager.GetCertificate(hello)
	}

	server := &http.Server{
		Addr:    _configuration.SSLListenAddress,
		Handler: WrapHandlerWithLogging(WrapHandlerWithRecovery(http.HandlerFunc(handleUP))),
		//Handler:   http.HandlerFunc(handleUP),
		TLSConfig: cfg,
	}
	return server
}
