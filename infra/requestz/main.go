package main

import (
	"crypto/ed25519"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/adamyi/CTFProxy/third_party/eddsa"
)

type Configuration struct {
	ListenAddress string
	VerifyKey     *ed25519.PublicKey
}

type Claims struct {
	Username    string   `json:"username"`
	Displayname string   `json:"displayname"`
	Service     string   `json:"service"`
	Groups      []string `json:"groups"`
	jwt.StandardClaims
}

var _configuration = Configuration{}

func readConfig() {
	var publicKeyPath string
	flag.StringVar(&_configuration.ListenAddress, "listen", "0.0.0.0:80", "http listen address")
	flag.StringVar(&publicKeyPath, "jwt_public_key", "", "Path to JWT public key")
	flag.Parse()
	JwtPubKey, err := ioutil.ReadFile(publicKeyPath)
	if err != nil {
		panic(err)
	}
	_configuration.VerifyKey, err = eddsa.ParseEdPublicKeyFromPEM(JwtPubKey)
	if err != nil {
		panic(err)
	}
}

func initRZRsp(rsp http.ResponseWriter) {
	rsp.Header().Add("Server", "requestz")
	rsp.Header().Add("Content-Type", "text/plain")
}

func handleRZ(rsp http.ResponseWriter, req *http.Request) {
	initRZRsp(rsp)

	if req.URL.Query().Get("deb") == "on" {
		rsp.Header().Add("X-CTFProxy-I-Debug", "1")
	}

	rs, _ := httputil.DumpRequest(req, true)

	rsp.Write(rs)

	tknStr := req.Header.Get("X-CTFProxy-JWT")

	claims := &Claims{}

	p := jwt.Parser{ValidMethods: []string{eddsa.SigningMethodEdDSA.Alg()}}
	tkn, err := p.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return _configuration.VerifyKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			rsp.Write([]byte("JWT: signature invalid\n"))
			return
		}
		rsp.Write([]byte("JWT: error\n"))
		rsp.Write([]byte(err.Error()))
		return
	}

	if !tkn.Valid {
		rsp.Write([]byte("JWT: invalid\n"))
		return
	}

	s, _ := json.Marshal(claims)

	rsp.Write(s)

}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	initRZRsp(w)
	w.Write([]byte("ok"))
}

func main() {
	rand.Seed(time.Now().UnixNano())
	readConfig()
	http.HandleFunc("/healthz", HealthCheckHandler)
	http.HandleFunc("/", handleRZ)
	err := http.ListenAndServe(_configuration.ListenAddress, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
