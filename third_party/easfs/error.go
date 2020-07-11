package main

import (
	"encoding/json"
	"net/http"
	"runtime/debug"
)

type UPError struct {
	Title         string
	Description   string
	PublicDebug   string
	InternalDebug string
	Code          int
}

func NewUPError(code int, title, description, publicDebug, internalDebug string) *UPError {
	ret := &UPError{Code: code, Title: title, Description: description, PublicDebug: publicDebug, InternalDebug: internalDebug}
	ret.InternalDebug += "\n\n===EASFS Stack Trace===\n" + string(debug.Stack())
	return ret
}

func ReturnError(rsp http.ResponseWriter, err *UPError) {
	rsp.Header().Set("Content-Type", "ctfproxy/error")
	rsp.WriteHeader(err.Code)
	json.NewEncoder(rsp).Encode(err)
}
