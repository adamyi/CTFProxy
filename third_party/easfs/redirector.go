package main

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Redirects struct {
	Redirects []struct {
		From string `yaml:"from"`
		To   string `yaml:"to"`
	} `yaml:"redirects"`
}

func ParseRedirects(filepath string) (Redirects, error) {
	redContent, err := ioutil.ReadFile(flagSitePath + filepath + "/_redirects.yaml")
	red := Redirects{}
	if err != nil {
		return red, err
	}
	err = yaml.Unmarshal(redContent, &red)
	return red, err
}

func GetRedirect(filepath string) (string, error) {
	// TODO: recursive check
	red, err := ParseRedirects("")
	if err != nil {
		return "", err
	}
	for _, r := range red.Redirects {
		if r.From == filepath {
			return r.To, nil
		}
	}
	return "", fmt.Errorf("no redirection found")
}
