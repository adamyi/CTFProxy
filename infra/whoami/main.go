package main

import (
	"crypto/rsa"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Configuration struct {
	ListenAddress string
	VerifyKey     *rsa.PublicKey
}

type Claims struct {
	Username    string   `json:"username"`
	Displayname string   `json:"displayname"`
	Service     string   `json:"service"`
	Groups      []string `json:"groups"`
	jwt.StandardClaims
}

var _configuration = Configuration{}
var ctf_domain string

func readConfig() {
	var publicKeyPath string
	flag.StringVar(&_configuration.ListenAddress, "listen", "0.0.0.0:80", "http listen address")
	flag.StringVar(&publicKeyPath, "jwt_public_key", "", "Path to JWT public key")
	flag.Parse()
	JwtPubKey, err := ioutil.ReadFile(publicKeyPath)
	if err != nil {
		panic(err)
	}
	_configuration.VerifyKey, err = jwt.ParseRSAPublicKeyFromPEM(JwtPubKey)
	if err != nil {
		panic(err)
	}
}

func initRZRsp(rsp http.ResponseWriter) {
	rsp.Header().Add("Server", "whoami")
	rsp.Header().Add("Content-Type", "text/plain")
}

func handleRZ(rsp http.ResponseWriter, req *http.Request) {
	initRZRsp(rsp)

	tknStr := req.Header.Get("X-CTFProxy-JWT")

	claims := &Claims{}

	p := jwt.Parser{ValidMethods: []string{jwt.SigningMethodRS256.Name}}
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

	if strings.HasPrefix(claims.Username, "anonymous@") {
		rsp.Write([]byte("You have not logged in! Please visit https://login." + ctf_domain))
		return
	}

	fmt.Fprintf(rsp, "Hello %s! You are authenticated as %s.", claims.Displayname, claims.Username)

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
