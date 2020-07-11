package main

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	Username    string   `json:"username"`
	Displayname string   `json:"displayname"`
	Service     string   `json:"service"`
	Groups      []string `json:"groups"`
	jwt.StandardClaims
}

var SubAccValid = regexp.MustCompile(`^[a-zA-Z.0-9\-_]+$`).MatchString

// username++impersonation+subacc@domain
// then displayname
func getUsername(req *http.Request) (string, string, error) {
	username, displayname := getMainUsername(req)
	var impersonateToken string
	if impersonateToken = req.Header.Get(_configuration.ImpersonateTokenHeader); impersonateToken != "" {
		impUsername, _ := getUsernameFromJWT(impersonateToken, username)
		if impUsername != "anonymous@anonymous."+_configuration.CorpDomain {
			s := strings.Split(username, "@")
			username = s[0] + "++" + strings.Split(strings.Split(impUsername, "@")[0], "+")[0] + "@" + s[1]
		}
	}
	subacc := req.Header.Get(_configuration.SubAccHeader)
	if subacc != "" {
		if !(SubAccValid(subacc) && len(subacc) < 10) {
			return "", "", errors.New("invalid subacc")
		}
		s := strings.Split(username, "@")
		username = s[0] + "+" + subacc + "@" + s[1]
	}
	return username, displayname, nil
}

func getGroups(username string) (ret []string) {
	r, err := http.Get(_configuration.GAIAEndpoint + "/api/getgroups?ldap=" + username)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer r.Body.Close()
	json.NewDecoder(r.Body).Decode(&ret)
	return
}

func getEmailFromRDN(name *pkix.Name) (email string, err error) {
	// most likely the last one so we loop in reverse
	for i := len(name.Names) - 1; i >= 0; i -= 1 {
		t := name.Names[i].Type
		if t[0] == 1 && t[1] == 2 && t[2] == 840 && t[3] == 113549 && t[4] == 1 && t[5] == 9 && t[6] == 1 {
			var ok bool
			email, ok = name.Names[i].Value.(string)
			if !ok {
				return "", errors.New("email not string")
			}
			return
		}
	}
	return "", errors.New("can't find email in cert")
}

// return handle, real name, and error
func verifyCert(cert *x509.Certificate) (string, string, error) {
	// log.Println("verifyCert")
	opts := x509.VerifyOptions{
		Roots: authCA,
	}
	_, err := cert.Verify(opts)
	if err != nil {
		log.Println(err)
		return "", "", err
	}
	email, err := getEmailFromRDN(&cert.Subject)
	if err != nil {
		return "", "", err
	}
	return strings.Split(email, "@")[0], cert.Subject.CommonName, nil
}

func getMainUsername(req *http.Request) (string, string) {
	certs := req.TLS.PeerCertificates
	// log.Println(certs)
	if len(certs) > 0 {
		name, displayname, err := verifyCert(certs[0])
		if err == nil {
			return name + "@" + _configuration.CorpDomain, displayname
		}
	}
	c, err := req.Cookie(_configuration.AuthCookieKey)
	var tknStr string
	if err != nil {
		if tknStr = req.Header.Get(_configuration.InternalJWTHeader); tknStr == "" {
			sn, err := getServiceNameFromIP(strings.Split(req.RemoteAddr, ":")[0])
			if err != nil {
				return "anonymous@anonymous." + _configuration.CorpDomain, "anonymous"
			}
			return sn + "@services." + _configuration.CorpDomain, sn
		}
	} else {
		tknStr = c.Value
	}
	return getUsernameFromJWT(tknStr, _configuration.ServiceName+"@services."+_configuration.CorpDomain)
}

func getUsernameFromJWT(tknStr, service string) (string, string) {
	claims := &Claims{}

	p := jwt.Parser{ValidMethods: []string{jwt.SigningMethodRS256.Name}}
	tkn, err := p.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return _configuration.VerifyKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			log.Println("Signature Invalid")
			return "anonymous@anonymous." + _configuration.CorpDomain, "anonymous"
		}
		log.Println("JWT Error")
		log.Println(err.Error())
		return "anonymous@anonymous." + _configuration.CorpDomain, "anonymous"
	}

	if !tkn.Valid {
		log.Println("JWT Invalid")
		return "anonymous@anonymous." + _configuration.CorpDomain, "anonymous"
	}

	if claims.Service != service {
		log.Printf("JWT not correct service - %v vs %v", claims.Service, service)
		return "anonymous@anonymous." + _configuration.CorpDomain, "anonymous"
	}

	return claims.Username, claims.Displayname
}
