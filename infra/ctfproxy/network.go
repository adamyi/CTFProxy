package main

import (
	"context"
	"errors"
	"github.com/gorilla/websocket"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	dialer = &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: false,
	}
)

func upDialContext(ctx context.Context, network, address string) (net.Conn, error) {
	if v := ctx.Value("up_real_addr"); v != nil {
		address = v.(string)
	}
	// fmt.Println("updialContext", address)
	return dialer.DialContext(ctx, network, address)
}

func getServiceNameFromDomain(domain string) (string, error) {
	host := strings.ToLower(domain)
	if !strings.HasSuffix(host, "."+_configuration.CorpDomain) {
		return "", errors.New("not corp domain")
	}

	sn := strings.ReplaceAll(host[:len(host)-len(_configuration.CorpDomain)-1], ".", "-dot-")

	if tsn, ok := _configuration.ServiceAliases.ConfigOrNil().(map[string]string)[sn]; ok {
		return tsn, nil
	}
	return sn, nil
}

func getRealAddr(host string) (string, error) {
	sn, err := getServiceNameFromDomain(host)
	if err != nil {
		return "", err
	}
	host = sn + "." + _configuration.ResolvingDomain

	// hack for GAE, which is no longer needed
	// if strings.HasSuffix(host, ".apps.geegle.org") {
	// 	host = "apps.geegle.org"
	// }

	ips, err := net.LookupIP(host)
	if err != nil {
		return "", err
	}

	for _, ip := range ips {
		if !isDockerIP(ip) {
			return "", errors.New("not internal ip")
		}
	}

	return ips[0].String() + ":80", nil
}

func getL2Addr(player string) (string, error) {
	if os.Getenv("CTFPROXY_CLUSTER") == "all" {
		return "", errors.New("levelshift disabled on all server")
	}
	if os.Getenv("CTFPROXY_CLUSTER") != "master" {
		player = "master"
	}
	host := player + ".prod." + _configuration.CorpDomain
	ips, err := net.LookupIP(host)
	// fmt.Println("getL2", ips, err, host)
	if err != nil || len(ips) == 0 {
		return "", errors.New("not valid internal")
	}

	return ips[0].String() + ":443", nil
}

func getNetworkContext(req *http.Request, username string) (context.Context, bool, error) {
	addr, err := getRealAddr(req.Host)
	if err == nil {
		return context.WithValue(context.Background(), "up_real_addr", addr), false, nil
	}

	if req.Header.Get("X-CTFProxy-LevelShift") == "1" {
		return context.Background(), false, errors.New("domain not present in two-level UP infra")
	}

	players := strings.Split(strings.Split(username, "@")[0], "+")
	hp := req.Header.Get("X-CTFProxy-Player")
	if hp != "" {
		players = append(players, hp)
	}
	// fmt.Println(players)
	for _, player := range players {
		addr, err = getL2Addr(player)
		if err == nil {
			return context.WithValue(context.Background(), "up_real_addr", addr), true, nil
		}
	}
	return context.Background(), false, errors.New("not found anywhere")
}

func init() {
	http.DefaultTransport.(*http.Transport).DialContext = upDialContext
	websocket.DefaultDialer.NetDialContext = upDialContext
}
