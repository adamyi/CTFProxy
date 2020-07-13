package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/adamyi/CTFProxy/third_party/eddsa"
)

type Claims struct {
	Username string `json:"username"`
	Service  string `json:"service"`
	jwt.StandardClaims
}

var slugre = regexp.MustCompile("[^a-z0-9]+")

func slugify(s string) string {
	return strings.Trim(slugre.ReplaceAllString(strings.ToLower(s), "_"), "_")
}

func getUsername(tknStr string) (string, string, error) {
	claims := &Claims{}

	p := jwt.Parser{ValidMethods: []string{eddsa.SigningMethodEdDSA.Alg()}}
	tkn, err := p.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return _configuration.VerifyKey, nil
	})

	if err != nil {
		return "", "", err
	}

	if !tkn.Valid {
		return "", "", fmt.Errorf("JWT Invalid")
	}
	username := strings.Replace(claims.Username, "++", "+", 1)
	userparts := strings.Split(username, "@")
	mainuser := strings.Split(userparts[0], "+")[0] + "@" + userparts[1]

	return slugify(username), slugify(mainuser), nil
}
