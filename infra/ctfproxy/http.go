package main

import (
	"github.com/google/uuid"
	"io"
	"net/http"
	"path"
	"strings"
)

func initUPRsp(rsp http.ResponseWriter) {
	rsp.Header().Add("Server", "CTFProxy/"+GitVersion)
	rsp.Header().Add("X-CTFProxy-Trace-Context", uuid.New().String())
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			if k == "Server" {
				dst.Set(k, v)
			} else {
				dst.Add(k, v)
			}
		}
	}
}

func copyResponse(rw http.ResponseWriter, resp *http.Response) error {
	copyHeader(rw.Header(), resp.Header)
	rw.WriteHeader(resp.StatusCode)
	defer resp.Body.Close()

	_, err := io.Copy(rw, resp.Body)
	return err
}

// cleanPath returns the canonical path for p, eliminating . and .. elements.
func cleanPath(p string) string {
	if p == "" {
		return "/"
	}
	if p[0] != '/' {
		p = "/" + p
	}
	np := path.Clean(p)
	// path.Clean removes trailing slash except for root;
	// put the trailing slash back if necessary.
	if p[len(p)-1] == '/' && np != "/" {
		// Fast path for common case of p being the string we want:
		if len(p) == len(np)+1 && strings.HasPrefix(p, np) {
			np = p
		} else {
			np += "/"
		}
	}
	return np
}
