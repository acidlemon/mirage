package main

import (
	"net/http"
	"net/http/httputil"
	"strings"
	"fmt"
	"net/url"

//	"github.com/acidlemon/go-dumper"
)

type ReverseProxy struct {
	cfg *Config
	domainMap map[string]ProxyInformation
}

func NewReverseProxy(cfg *Config) *ReverseProxy{
	return &ReverseProxy{
		cfg: cfg,
		domainMap: map[string]ProxyInformation{},
	}
}

func (r *ReverseProxy) ServeHTTPWithPort(w http.ResponseWriter, req *http.Request, port int) {
	subdomain := strings.ToLower(strings.Split(req.Host, ".")[0])

	if _, ok := r.domainMap[subdomain]; !ok {
		fmt.Println("subdomain not found: ", subdomain)
		http.NotFound(w, req)
		return
	}

	if handler, ok := r.domainMap[subdomain].proxyHandlers[port]; ok {
		handler.ServeHTTP(w, req)
	} else {
		http.NotFound(w, req)
	}
}


type ProxyInformation struct {
	IPAddress string
	proxyHandlers map[int]http.Handler
}

func (r *ReverseProxy) AddSubdomain(subdomain string, ipaddress string) {
	handlers := make(map[int]http.Handler)

	// create reverse proxy
	for _, v := range r.cfg.Listen.HTTP {
		destUrlString := fmt.Sprintf("http://%s:%d", ipaddress, v.TargetPort)
		destUrl, _ := url.Parse(destUrlString)
		handler := httputil.NewSingleHostReverseProxy(destUrl)

		handlers[v.ListenPort] = handler
	}

	fmt.Println("add subdomain: ", subdomain)

	// add to map
	r.domainMap[subdomain] = ProxyInformation{
		IPAddress: ipaddress,
		proxyHandlers: handlers,
	}
}

func (r *ReverseProxy) RemoveSubdomain(subdomain string) {
	delete(r.domainMap, subdomain)
}



