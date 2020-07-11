package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Promo struct {
	Label       string `yaml:"label"`
	Description string `yaml:"description"`
	Path        string `yaml:"path"`
	Icon        string `yaml:"icon"`
}

type Linkbox struct {
	Name     string `yaml:"name"`
	Contents []struct {
		Label string `yaml:"label"`
		Path  string `yaml:"path"`
	} `yaml:"contents"`
}

type Footer struct {
	Footer []struct {
		Promos    []Promo   `yaml:"promos,omitempty"`
		Linkboxes []Linkbox `yaml:"linkboxes,omitempty"`
		Banner    string    `yaml:"banner,omitempty"`
	} `yaml:"footer"`
}

func ParseFooter(filepath string) (string, []Promo, []Linkbox, error) {
	footerContent, err := ioutil.ReadFile(flagSitePath + filepath)
	if err != nil {
		return "", nil, nil, err
	}
	footer := Footer{}
	promos := []Promo{}
	linkboxes := []Linkbox{}
	banner := ""
	err = yaml.Unmarshal(footerContent, &footer)
	for _, f := range footer.Footer {
		promos = append(promos, f.Promos...)
		linkboxes = append(linkboxes, f.Linkboxes...)
		if f.Banner != "" {
			banner = f.Banner
		}
	}
	return banner, promos, linkboxes, nil
}
