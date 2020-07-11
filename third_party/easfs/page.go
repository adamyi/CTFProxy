package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func GetIndex(w http.ResponseWriter, path string) error {
	extensions := [3]string{"_index.yaml", "index.md", "index.html"}
	for _, ext := range extensions {
		fileLocation := flagSitePath + path + ext
		// fmt.Printf("checking %s\n", fileLocation)
		content, err := ioutil.ReadFile(fileLocation)
		if err == nil {
			if ext == "index.md" {
				return ParseMD(w, content, path)
			} else if ext == "_index.yaml" {
				return ParseYAML(w, content, path)
			}
			w.Header().Set("Content-Type", "text/html")
			w.Write(content)
			return nil
		}
	}
	return fmt.Errorf("file not found")
}

func GetPage(w http.ResponseWriter, path string) error {
	extensions := [5]string{".md", ".html", ".json", ""}
	for _, ext := range extensions {
		fileLocation := flagSitePath + path + ext
		// fmt.Printf("checking %s\n", fileLocation)
		content, err := ioutil.ReadFile(fileLocation)
		if err == nil {
			if ext == ".md" {
				w.Header().Set("Content-Type", "text/html")
				return ParseMD(w, content, path)
			} else if ext == ".html" {
				// return ParseHTML(content)
			}
			if strings.HasSuffix(fileLocation, ".html") {
				w.Header().Set("Content-Type", "text/html")
			} else if strings.HasSuffix(fileLocation, ".js") {
				w.Header().Set("Content-Type", "application/javascript")
			} else if strings.HasSuffix(fileLocation, ".css") {
				w.Header().Set("Content-Type", "text/css")
			} else if strings.HasSuffix(fileLocation, ".json") {
				w.Header().Set("Content-Type", "application/json")
			} else if strings.HasSuffix(fileLocation, ".svg") {
				w.Header().Set("Content-Type", "image/svg+xml")
			}
			w.Write(content)
			return nil
		}
	}
	return fmt.Errorf("file not found")
}
