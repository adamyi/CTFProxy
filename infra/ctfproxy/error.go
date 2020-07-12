package main

import (
	"encoding/json"
	"fmt"
	"html"
	"html/template"
	"log"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/adamyi/CTFProxy/infra/ctfproxy/templates"
)

type CPError struct {
	Title         string
	Description   string
	PublicDebug   string
	InternalDebug string
	Type          string
	Debug         template.HTML
	Code          int
}

var errorTemplate *template.Template

func init() {
	errorTemplate = template.Must(template.New("error").Parse(templates.Data["error.html"]))
}

func (e *CPError) SetType(t string) {
	e.Type = t
}

func NewCPError(code int, title, description, publicDebug, internalDebug string) *CPError {
	ret := &CPError{Code: code, Title: title, Description: description, PublicDebug: publicDebug, InternalDebug: internalDebug}
	ret.InternalDebug += "\n\n===CTFProxy Stack Trace===\n" + string(debug.Stack())
	ret.Type = "cp"
	return ret
}

func (e CPError) Write(w http.ResponseWriter) {
	if e.Code >= 300 && e.Code < 400 {
		w.Header().Add("Location", e.Description)
	}
	w.WriteHeader(e.Code)
	dbgstr := e.PublicDebug
	if w.Header().Get("X-CTFProxy-I-Debug") == "1" {
		dbgstr += "\n===Internal Debug Info because you're in ctfproxy-debug@===\n" + e.InternalDebug
	}
	e.Debug = template.HTML(strings.Replace(html.EscapeString(dbgstr), "\n", "<br>\n", -1))
	errorTemplate.Execute(w, e)
}

func WrapHandlerWithRecovery(wrappedHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				log.Println("[PANIC RECOVERY TRIGGERED] something terrible happened.")
				log.Println(err)
				NewCPError(http.StatusInternalServerError, "Server Internal Error", "Something went terribly wrong...", "", fmt.Sprintf("===panic recovery===\n%v", err)).Write(w)
			}
		}()
		wrappedHandler.ServeHTTP(w, r)
	})
}

func handleUpstreamCPError(rw http.ResponseWriter, resp *http.Response) {
	defer resp.Body.Close()
	var e CPError
	err := json.NewDecoder(resp.Body).Decode(&e)
	if err != nil {
		e = *NewCPError(http.StatusInternalServerError, "Server Internal Error", "Something went wrong with the service, please try again later", "", "upstream service returned ctfproxy/uperror but failed to decode CPError json "+err.Error())
	}
	e.SetType("s")
	e.Write(rw)
}
