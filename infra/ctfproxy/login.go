package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
	"github.com/adamyi/CTFProxy/infra/ctfproxy/templates"
)

func verifyPassword(email, password string) bool {

	data, err := json.Marshal(struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{email, password})
	if err != nil {
		fmt.Println(err)
		return false
	}

	r, err := http.Post(_configuration.GAIAEndpoint+"/api/login", "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Println(err)
		return false
	}
	r.Body.Close()

	return r.StatusCode == 200
}

var loginTemplate *template.Template

func init() {
	loginTemplate = template.Must(template.New("login").Parse(templates.Data["login.html"]))
}

func handleLogin(rsp http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		loginTemplate.Execute(rsp, "")
	} else if req.Method == "POST" {
		req.ParseForm()
		username := strings.ToLower(req.Form.Get("username"))
		if !strings.HasSuffix(username, "@"+_configuration.CorpDomain) {
			username = username + "@" + _configuration.CorpDomain
		}
		password := req.Form.Get("password")

		if !verifyPassword(username, password) {
			loginTemplate.Execute(rsp, "Incorrect password")
			return
		}

		expirationTime := time.Now().Add(24 * 30 * time.Hour)
		pclaims := Claims{
			Username:    username,
			Displayname: strings.Split(username, "@")[0],
			Service:     _configuration.ServiceName + "@services." + _configuration.CorpDomain,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
			},
		}
		ptoken := jwt.NewWithClaims(jwt.SigningMethodRS256, pclaims)
		ptstr, err := ptoken.SignedString(_configuration.SignKey)
		if err != nil {
			returnError(NewUPError(http.StatusInternalServerError, "Internal Server Error", "Something went wrong while generating JWT", "", err.Error()), rsp)
			return
		}
		authcookie := &http.Cookie{Name: _configuration.AuthCookieKey, Value: ptstr, HttpOnly: true, Domain: _configuration.CorpDomain}
		http.SetCookie(rsp, authcookie)
		if req.URL.Query().Get("return_url") == "" {
			http.Redirect(rsp, req, "https://whoami."+_configuration.CorpDomain, http.StatusFound)
		} else {
			http.Redirect(rsp, req, req.URL.Query().Get("return_url"), http.StatusFound)
		}
	} else {
		returnError(NewUPError(http.StatusMethodNotAllowed, "Method Not Allowed", "Only GET and POST are supported", "", ""), rsp)
	}
}
