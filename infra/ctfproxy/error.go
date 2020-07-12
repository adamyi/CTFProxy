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

type UPError struct {
	Title         string
	Description   string
	PublicDebug   string
	InternalDebug string
	Debug         template.HTML
	Code          int
}

var errorTemplate *template.Template

func init() {
	errorTemplate = template.Must(template.New("error").Parse(templates.Data["error.html"]))
}

func NewUPError(code int, title, description, publicDebug, internalDebug string) *UPError {
	ret := &UPError{Code: code, Title: title, Description: description, PublicDebug: publicDebug, InternalDebug: internalDebug}
	ret.InternalDebug += "\n\n===CTFProxy Stack Trace===\n" + string(debug.Stack())
	return ret
}

func returnError(err *UPError, rsp http.ResponseWriter) {
	if err.Code >= 300 && err.Code < 400 {
		rsp.Header().Add("Location", err.Description)
	}
	rsp.WriteHeader(err.Code)
	dbgstr := err.PublicDebug
	if rsp.Header().Get("X-CTFProxy-I-Debug") == "1" {
		dbgstr += "\n===Internal Debug Info because you're in ctfproxy-debug@===\n" + err.InternalDebug
	}
	err.Debug = template.HTML(strings.Replace(html.EscapeString(dbgstr), "\n", "<br>\n", -1))
	errorTemplate.Execute(rsp, err)
}

func WrapHandlerWithRecovery(wrappedHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				log.Println("[PANIC RECOVERY TRIGGERED] something terrible happened.")
				log.Println(err)
				returnError(NewUPError(http.StatusInternalServerError, "Server Internal Error", "Something went terribly wrong...", "", fmt.Sprintf("===panic recovery===\n%v", err)), w)
			}
		}()
		wrappedHandler.ServeHTTP(w, r)
	})
}

func handleUpstreamUPError(rw http.ResponseWriter, resp *http.Response) {
	defer resp.Body.Close()
	var e UPError
	err := json.NewDecoder(resp.Body).Decode(&e)
	if err != nil {
		e = *NewUPError(http.StatusInternalServerError, "Server Internal Error", "Something went wrong with the service, please try again later", "", "upstream service returned ctfproxy/uperror but failed to decode UPError json "+err.Error())
	}
	returnError(&e, rw)
}
