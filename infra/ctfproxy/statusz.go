// Based off github.com/youtube/doorman/blob/master/go/status/status.go

// Copyright 2019 Adam Yi
// Copyright 2018 Serge Bazanski
// Copyright 2016 Google, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"html"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/user"
	"strconv"
	"sync"
	"time"
)

var (
	// binaryName  = filepath.Base(os.Args[0])
	binaryHash  string
	hostname    string
	username    string
	serverStart = time.Now()

	lock     sync.Mutex
	sections []section
	tmpl     = template.Must(reparse(nil))
	funcs    = make(template.FuncMap)

	GitCommit      = "unknown"
	GitVersion     = "unknown"
	Builder        = "unknown"
	BuildTimestamp = "0"
)

type section struct {
	Banner   string
	Fragment string
	F        func() interface{}
}

var statusHTML = `<!DOCTYPE html>
<html>
<head>
<title>Status for {{.BinaryName}}</title>
<style>
body {
background: #fff;
}
h1 {
font-family: sans-serif;
clear: both;
width: 100%;
text-align: center;
font-size: 120%;
padding-top: 0.3em;
padding-bottom: 0.3em;
background: #eeeeff;
margin-top: 1em;
}
.lefthand {
float: left;
width: 80%;
}
.righthand {
text-align: right;
}
.table td,th {
background-color: #eee;
padding: 0.2em 0.4em 0.2em 0.4em;
}
.table th {
background-color: #c0c0c0;
}
.table {
background-color: #fff;
border-spacing: 0.2em;
margin-left: auto;
margin-right: auto;
}
</style>
</head>
<h1>Status for {{.BinaryName}}</h1>
<div>
<div class=lefthand>
Started: {{.StartTime}} -- up {{.Up}}<br>
Built on {{.BuildTime}}<br>
Built at {{.Builder}}<br>
Built from git checkout <a href="https://github.com/adamyi/CTFProxy/tree/{{.GitCommit}}">{{.GitVersion}}</a><br>
SHA256 {{.BinaryHash}}<br>
Access Policies last updated at {{.AccessPolicyTime}}<br>
Aliases last updated at {{.AliasTime}}<br>
</div>
<div class=righthand>
Running as {{.Username}} on {{.Hostname}}<br>
</div>
<div style="clear: both;"> </div>
</div>`

const accessPoliciesFragment = `
    <style type="text/css">
	</style>
    <div>
        <table class="table">
            <tr>
                <th>service</th>
		<th>policy</th>
            </tr>
	    {{ range $key, $value := . }}
		<tr>
			<td>{{ $key }}</td>
			<td><pre>{{ $value }}</pre></td>
		</tr>
	{{ end }}
        </table>
    </div>
`

func reparse(sections []section) (*template.Template, error) {
	var buf bytes.Buffer

	io.WriteString(&buf, `{{define "status"}}`)
	io.WriteString(&buf, statusHTML)

	for i, sec := range sections {
		fmt.Fprintf(&buf, "<h1>%s</h1>\n", html.EscapeString(sec.Banner))
		fmt.Fprintf(&buf, "{{$sec := index .Sections %d}}\n", i)
		fmt.Fprintf(&buf, `{{template "sec-%d" call $sec.F}}`+"\n", i)
	}
	fmt.Fprintf(&buf, `</html>`)
	io.WriteString(&buf, "{{end}}\n")

	for i, sec := range sections {
		fmt.Fprintf(&buf, `{{define "sec-%d"}}%s{{end}}\n`, i, sec.Fragment)
	}
	return template.New("").Funcs(funcs).Parse(buf.String())
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "ok")
}

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-Frame-Options", "SAMEORIGIN")
	w.Header().Set("Server", "ctfproxy")
	lock.Lock()
	defer lock.Unlock()

	buildTime := time.Unix(0, 0)
	if buildTimeNum, err := strconv.Atoi(BuildTimestamp); err == nil {
		buildTime = time.Unix(int64(buildTimeNum), 0)
	}

	data := struct {
		Sections         []section
		BinaryName       string
		BinaryHash       string
		GitVersion       string
		GitCommit        string
		Builder          string
		BuildTime        string
		Hostname         string
		Username         string
		StartTime        string
		AccessPolicyTime string
		AliasTime        string
		Up               string
	}{
		Sections:         sections,
		BinaryName:       "CTFProxy",
		BinaryHash:       binaryHash,
		GitVersion:       GitVersion,
		GitCommit:        GitCommit,
		Builder:          Builder,
		BuildTime:        fmt.Sprintf("%s (%d)", buildTime.Format(time.RFC1123), buildTime.Unix()),
		Hostname:         hostname,
		Username:         username,
		StartTime:        serverStart.Format(time.RFC1123),
		AccessPolicyTime: _configuration.AccessPolicies.LastUpdated().Format(time.RFC1123),
		AliasTime:        _configuration.ServiceAliases.LastUpdated().Format(time.RFC1123),
		Up:               time.Since(serverStart).String(),
	}

	if err := tmpl.ExecuteTemplate(w, "status", data); err != nil {
		log.Printf("servenv: couldn't execute template: %v\n", err)
	}
}

func init() {
	var err error
	hostname, err = os.Hostname()
	if err != nil {
		log.Fatalf("os.Hostname: %v", err)
	}

	user, err := user.Current()
	if err != nil {
		log.Fatalf("user.Current: %v", err)
	}
	username = fmt.Sprintf("%s (%s)", user.Username, user.Uid)

	exec, err := os.Executable()
	if err == nil {
		f, err := os.Open(exec)
		if err == nil {
			h := sha256.New()
			if _, err := io.Copy(h, f); err != nil {
				log.Fatalf("io.Copy: %v", err)
			}
			binaryHash = fmt.Sprintf("%x", h.Sum(nil))
		} else {
			log.Printf("Could not get SHA256 of binary: os.Open(%q): %v\n", exec, err)
			binaryHash = "could not read executable"
		}
	} else {
		log.Printf("Could not get SHA256 of binary: os.Executable(): %v\n", err)
		binaryHash = "could not get executable"
	}

	// http.HandleFunc("/statusz", StatusHandler)
	// http.HandleFunc("/healthz", HealthCheckHandler)
}

func init() {
	AddStatusPart("AccessPolicies", accessPoliciesFragment, func(ctx context.Context) interface{} {
		return _configuration.AccessPolicies.ConfigOrNil().(map[string]string)
	})
}

// AddStatusPart adds a new section to status. frag is used as a
// subtemplate of the template used to render /statusz, and will
// be executed using the value of invoking f at the time of the
// /statusz request. frag is parsed and executed with the
// html/template package. Functions registered with AddStatusFuncs
// may be used in the template.
func AddStatusPart(banner, frag string, f func(context.Context) interface{}) {
	lock.Lock()
	defer lock.Unlock()

	secs := append(sections, section{
		Banner:   banner,
		Fragment: frag,
		F:        func() interface{} { return f(context.Background()) },
	})

	var err error
	tmpl, err = reparse(secs)
	if err != nil {
		secs[len(secs)-1] = section{
			Banner:   banner,
			Fragment: "<code>bad status template: {{.}}</code>",
			F:        func() interface{} { return err },
		}
	}
	tmpl, _ = reparse(secs)
	sections = secs
}

// AddStatusSection registers a function that generates extra
// information for /statusz. If banner is not empty, it will be
// used as a header before the information. If more complex output
// than a simple string is required use AddStatusPart instead.
func AddStatusSection(banner string, f func(context.Context) string) {
	AddStatusPart(banner, `{{.}}`, func(ctx context.Context) interface{} { return f(ctx) })
}
