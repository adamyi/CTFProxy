package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

type uplogEntry struct {
	ClientIP    string `json:"clientIP"`
	Host        string `json:"host"`
	RequestTime string `json:"requestTime"`
	Latency     int64  `json:"latency"`
	Request     struct {
		URI     string      `json:"uri"`
		Method  string      `json:"method"`
		Version string      `json:"version"`
		Body    string      `json:"body"`
		Header  http.Header `json:"headers"`
	} `json:"request"`
	Response struct {
		StatusCode int         `json:"statusCode"`
		Body       string      `json:"body"`
		Header     http.Header `json:"headers"`
	} `json:"response"`
}

func WrapHandlerWithLogging(wrappedHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// log.Printf("%v - %v %v%v", req.RemoteAddr, req.Method, req.Host, req.RequestURI)
		// don't log k8s health check for ctfproxy
		if req.Host == "ctfproxyz."+_configuration.CorpDomain && req.RequestURI == "/healthz" {
			wrappedHandler.ServeHTTP(w, req)
			return
		}
		// don't log kibana/elasticsearch requests
		if req.Host == "kibana."+_configuration.CorpDomain || req.Host == "elasticsearch."+_configuration.CorpDomain {
			wrappedHandler.ServeHTTP(w, req)
			return
		}
		entry := uplogEntry{}
		t := time.Now()
		entry.RequestTime = t.UTC().Format(time.RFC3339)
		entry.ClientIP = strings.Split(req.RemoteAddr, ":")[0] // we don't trust any proxy except ourselves
		entry.Host = req.Host
		entry.Request.URI = req.RequestURI
		entry.Request.Method = req.Method
		entry.Request.Version = fmt.Sprintf("%d.%d", req.ProtoMajor, req.ProtoMinor)
		entry.Request.Header = req.Header
		// http.MaxBytesReader might be better but let's just use io.LimitedReader since we are doing the wrapped logger.
		limitedReader := &io.LimitedReader{R: req.Body, N: 10485760}
		reqcontent, lrerr := ioutil.ReadAll(limitedReader)
		req.Body = ioutil.NopCloser(bytes.NewReader(reqcontent))
		entry.Request.Body = string(reqcontent)
		buf := new(bytes.Buffer)
		lrw := newUplogResponseWriter(buf, w, &entry)
		if lrerr != nil {
			NewCPError(http.StatusBadRequest, "You issued a malformed request", "Entity Too Large", "", "").Write(lrw, req)
		} else if limitedReader.N < 1 {
			NewCPError(http.StatusRequestEntityTooLarge, "You issued a malformed request", "Entity Too Large", "", "").Write(lrw, req)
		} else if len(reqcontent) > 0 && req.Method == "GET" {
			NewCPError(http.StatusBadRequest, "You issued a malformed request", "No Body for GET", "", "").Write(lrw, req)
		} else {
			wrappedHandler.ServeHTTP(lrw, req)
		}
		entry.Response.Body = buf.String()
		// entry.Response.Header = lrw.Header()
		entry.Latency = (time.Now().UnixNano() - t.UnixNano()) / (int64(time.Millisecond) / int64(time.Nanosecond))
		log.Printf("%v %v - %v - %v %v%v [%vms]\n", entry.ClientIP, entry.Response.Header.Get("X-CTFProxy-I-User"), entry.Response.StatusCode, entry.Request.Method, entry.Host, entry.Request.URI, entry.Latency)
		go func() {
			estr, _ := json.Marshal(entry)
			rsp, err := http.Post("http://elasticsearch."+_configuration.ResolvingDomain+"/ctfproxy/_doc", "application/json", bytes.NewBuffer(estr))
			if err == nil {
				rsp.Body.Close()
			} else {
				log.Printf("Error logging: %v", err.Error())
			}
		}()
	})
}

type uplogResponseWriter struct {
	file        io.Writer
	resp        http.ResponseWriter
	multi       io.Writer
	entry       *uplogEntry
	wroteHeader bool
}

func newUplogResponseWriter(file io.Writer, resp http.ResponseWriter, entry *uplogEntry) http.ResponseWriter {
	multi := io.MultiWriter(file, resp)
	return &uplogResponseWriter{
		file:        file,
		resp:        resp,
		multi:       multi,
		entry:       entry,
		wroteHeader: false,
	}
}

// implement http.ResponseWriter
// https://golang.org/pkg/net/http/#ResponseWriter
func (w *uplogResponseWriter) Header() http.Header {
	return w.resp.Header()
}

func (w *uplogResponseWriter) Write(b []byte) (int, error) {
	return w.multi.Write(b)
}

func (w *uplogResponseWriter) WriteHeader(i int) {
	if !w.wroteHeader {
		w.wroteHeader = true
		w.entry.Response.StatusCode = i
		w.entry.Response.Header = w.resp.Header().Clone()
		if w.resp.Header().Get("X-CTFProxy-I-Debug") == "" {
			for k := range w.resp.Header() {
				if strings.HasPrefix(strings.ToLower(k), "x-ctfproxy-i-") {
					w.resp.Header().Del(k)
				}
			}
		}
		w.resp.WriteHeader(i)
	}
}

// websocket needs this
func (w *uplogResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hj, ok := w.resp.(http.Hijacker); ok {
		return hj.Hijack()
	}
	return nil, nil, errors.New("Error in hijacker")
}
