package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/flosch/pongo2"
	"github.com/gosimple/slug"
)

func Slugify(in *pongo2.Value, param *pongo2.Value) (out *pongo2.Value, err *pongo2.Error) {
	if !in.IsString() {
		return nil, &pongo2.Error{
			OrigError: fmt.Errorf("only strings should be sent to the slugify filter"),
		}
	}

	s := in.String()
	s = slug.Make(s)

	return pongo2.AsValue(s), nil
}

func IsDir(filepath string) bool {
	fi, err := os.Stat(filepath)
	if err != nil {
		return false
	}
	return fi.IsDir()
}

func GetInclude(include []byte) []byte {
	// {% include "_shared/latest_articles.html" %}
	// fmt.Println("src/content/" + string(include[12:len(include)-4]))
	inc, _ := ioutil.ReadFile(flagSitePath + string(include[12:len(include)-4]))
	return inc
}

func RenderContent(content []byte) []byte {
	// TODO: includecode & htmlescape
	// fmt.Println(string(content))
	content = regexp.MustCompile(`(?ms){%[ ]?include .+%}`).ReplaceAllFunc(content, GetInclude)

	return content
}
