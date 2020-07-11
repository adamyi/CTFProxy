package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
)

var (
	wsUpgrader = websocket.Upgrader{CheckOrigin: wsCheckOrigin}
	wsDialer   = websocket.DefaultDialer
)

func wsCheckOrigin(r *http.Request) bool {
	o := r.Header.Get("Origin")
	h := r.Host
	if o == "" || h == "" {
		log.Print("Websocket missing origin and/or host")
		return false
	}
	ou, err := url.Parse(o)
	if err != nil {
		log.Printf("Couldn't parse url: %v", err)
		return false
	}
	if ou.Host != h && ou.Host != "cli-relay.corp."+_configuration.CorpDomain {
		log.Print("Origin doesn't match host")
		return false
	}
	return true
}

// adapted from https://github.com/koding/websocketproxy/blob/master/websocketproxy.go
func handleWs(ctx context.Context, rsp http.ResponseWriter, req *http.Request, jwttoken string, levelShift bool) {
	requestHeader := http.Header{}
	requestHeader.Set("Host", req.Host)
	requestHeader.Set(_configuration.InternalJWTHeader, jwttoken)
	if origin := req.Header.Get("Origin"); origin != "" {
		requestHeader.Add("Origin", origin)
	}
	for _, prot := range req.Header[http.CanonicalHeaderKey("Sec-WebSocket-Protocol")] {
		requestHeader.Add("Sec-WebSocket-Protocol", prot)
	}
	for _, cookie := range req.Header[http.CanonicalHeaderKey("Cookie")] {
		requestHeader.Add("Cookie", cookie)
	}
	backendURL := *req.URL
	backendURL.Host = req.Host
	backendURL.Scheme = "ws"

	if levelShift {
		backendURL.Scheme = "wss"
		requestHeader.Add("X-CTFProxy-LevelShift", "1")
	}

	// dial backend
	connBackend, resp, err := wsDialer.DialContext(ctx, backendURL.String(), requestHeader)
	if err != nil {
		log.Printf("couldn't dial to remote backend url %s %s", backendURL.String(), err)
		if resp != nil {
			// If the WebSocket handshake fails, ErrBadHandshake is returned
			// along with a non-nil *http.Response so that callers can handle
			// redirects, authentication, etcetera.
			if err := copyResponse(rsp, resp); err != nil {
				log.Printf("couldn't write response after failed remote backend handshake: %s", err)
			}
		} else {
			http.Error(rsp, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
		}
		return
	}
	defer connBackend.Close()

	// Only pass those headers to the upgrader.
	upgradeHeader := http.Header{}
	if hdr := resp.Header.Get("Sec-Websocket-Protocol"); hdr != "" {
		upgradeHeader.Set("Sec-Websocket-Protocol", hdr)
	}
	if hdr := resp.Header.Get("Set-Cookie"); hdr != "" {
		upgradeHeader.Set("Set-Cookie", hdr)
	}

	connPub, err := wsUpgrader.Upgrade(rsp, req, upgradeHeader)
	if err != nil {
		log.Printf("couldn't upgrade %s", err)
		return
	}
	defer connPub.Close()

	errClient := make(chan error, 1)
	errBackend := make(chan error, 1)
	replicateWebsocketConn := func(dst, src *websocket.Conn, errc chan error) {
		for {
			msgType, msg, err := src.ReadMessage()
			if err != nil {
				m := websocket.FormatCloseMessage(websocket.CloseNormalClosure, fmt.Sprintf("%v", err))
				if e, ok := err.(*websocket.CloseError); ok {
					if e.Code != websocket.CloseNoStatusReceived {
						m = websocket.FormatCloseMessage(e.Code, e.Text)
					}
				}
				errc <- err
				dst.WriteMessage(websocket.CloseMessage, m)
				break
			}
			err = dst.WriteMessage(msgType, msg)
			if err != nil {
				errc <- err
				break
			}
		}
	}

	go replicateWebsocketConn(connPub, connBackend, errClient)
	go replicateWebsocketConn(connBackend, connPub, errBackend)

	var message string
	select {
	case err = <-errClient:
		message = "Error when copying from backend to client: %v"
	case err = <-errBackend:
		message = "Error when copying from client to backend: %v"

	}
	if e, ok := err.(*websocket.CloseError); !ok || e.Code == websocket.CloseAbnormalClosure {
		log.Printf(message, err)
	}

}
