package main

import (
	"math/rand"
	"net/http"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	readConfig()
	initRateLimit()

	server := buildSSLServer()
	go http.ListenAndServe(_configuration.ListenAddress, certManager.HTTPHandler(http.HandlerFunc(redirectSSL)))
	server.ListenAndServeTLS("", "")
}
