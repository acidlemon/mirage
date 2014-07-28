package reverseproxy

import (
	"net/http"
	"net/http/httputil"
	"strings"
	"fmt"
	"net/url"

//	"github.com/acidlemon/go-dumper"
)

var DefaultReverseProxy = NewReverseProxy()

func NewReverseProxy() *ReverseProxy{
	return &ReverseProxy{
		domainMap: map[string]ProxyInformation{},
	}
}


type ReverseProxy struct {
	HostSuffix string
	domainMap map[string]ProxyInformation
}


func ServeHTTP(w http.ResponseWriter, req *http.Request) {
	DefaultReverseProxy.ServeHTTP(w, req)
}
func (r *ReverseProxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	subdomain := strings.ToLower(strings.Split(req.Host, ".")[0])

	if _, ok := r.domainMap[subdomain]; !ok {
		fmt.Println("subdomain not found: ", subdomain)
		http.NotFoundHandler().ServeHTTP(w, req)
		return
	}

	fmt.Println("found subdomain:", subdomain)

	r.domainMap[subdomain].proxyHandler.ServeHTTP(w, req)
}


type ProxyInformation struct {
	IPAddress string
	proxyHandler http.Handler
}

func AddSubdomain(subdomain string, ipaddress string) {
	fmt.Println("AddSubdomain")
	DefaultReverseProxy.AddSubdomain(subdomain, ipaddress)
}

func RemoveSubdomain(subdomain string) {
	DefaultReverseProxy.RemoveSubdomain(subdomain)
}

func SetHostSuffix(suffix string) {
	DefaultReverseProxy.HostSuffix = suffix
}

func (r *ReverseProxy) AddSubdomain(subdomain string, ipaddress string) {
	// create reverse proxy
	destUrlString := fmt.Sprintf("http://%s:%d", ipaddress, 5000)
	destUrl, _ := url.Parse(destUrlString)
	handler := httputil.NewSingleHostReverseProxy(destUrl)

	fmt.Println("add subdomain: ", subdomain)

	// add to map
	r.domainMap[subdomain] = ProxyInformation{
		IPAddress: ipaddress,
		proxyHandler: handler,
	}
}

func (r *ReverseProxy) RemoveSubdomain(subdomain string) {
	delete(r.domainMap, subdomain)
}



