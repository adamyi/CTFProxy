package main

import (
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	Username string `json:"username"`
	Service  string `json:"service"`
	jwt.StandardClaims
}

func getUsername(tknStr string) (string, error) {
	claims := &Claims{}

	p := jwt.Parser{ValidMethods: []string{jwt.SigningMethodRS256.Name}}
	tkn, err := p.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return _configuration.VerifyKey, nil
	})

	if err != nil {
		return "", err
	}

	if !tkn.Valid {
		return "", fmt.Errorf("JWT Invalid")
	}
	return strings.Split(claims.Username, "@")[0], nil
}
