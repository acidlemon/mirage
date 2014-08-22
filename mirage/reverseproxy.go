package mirage

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

	fmt.Println("host is", req.Host)
	if _, ok := r.domainMap[subdomain]; !ok {
		fmt.Println("subdomain not found: ", subdomain)
		http.NotFound(w, req)
		return
	}

	fmt.Println("found subdomain:", subdomain, "port:", port)

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
	for listen, target := range r.cfg.ListenPorts {
		destUrlString := fmt.Sprintf("http://%s:%d", ipaddress, target)
		destUrl, _ := url.Parse(destUrlString)
		handler := httputil.NewSingleHostReverseProxy(destUrl)

		handlers[listen] = handler
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



