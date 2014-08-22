package mirage

import (
	"sync"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sort"

//	"github.com/acidlemon/go-dumper"
)

var app *Mirage

type Config struct {
	ForeignAddress string
	WebApiHost string
	ReverseProxyHostSuffix string
	ListenPorts map[int]int // map[ListenPort] = ProxyPassPort
	ListenSSLPorts map[int]int
	DockerEndpoint string
	DefaultImage string
}

func NewConfig() *Config {
	// default config
	cfg := &Config{
		ForeignAddress: "127.0.0.1",
		WebApiHost: "localhost",
		ReverseProxyHostSuffix: ".example.net",
		ListenPorts: map[int]int{ 8080: 8080 },
		ListenSSLPorts: map[int]int{},
		DockerEndpoint: "unix:///var/run/docker.sock",
		DefaultImage: "",
	}

	return cfg
}


type Mirage struct {
	Config *Config
	WebApi *WebApi
	ReverseProxy *ReverseProxy
	Docker *Docker
}

func Setup(cfg *Config) {
	m := &Mirage{
		Config: cfg,
		WebApi: NewWebApi(cfg),
		ReverseProxy: NewReverseProxy(cfg),
		Docker: NewDocker(cfg),
	}

	infolist, err := m.Docker.List()
	if err != nil {
		fmt.Println("cannot initialize reverse proxy: ", err.Error())
	}

	for _, info := range infolist {
		m.ReverseProxy.AddSubdomain(info.SubDomain, info.IPAddress)
	}

	app = m
}

func Run() {
	// launch server
	var wg sync.WaitGroup
	for k, _ := range app.Config.ListenPorts {
		wg.Add(1)
		go func(port int) {
			defer wg.Done()
			laddr := fmt.Sprintf("%s:%d", app.Config.ForeignAddress, port)
			listener, err := net.Listen("tcp", laddr)
			if err != nil {
				fmt.Println("cannot listen %s", laddr)
				return
			}

			mux := http.NewServeMux()
			mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
				app.ServeHTTPWithPort(w, req, port)
			})

			http.Serve(listener, mux)
		}(k)
	}

	// TODO SSL Support

	wg.Wait()
}

func (m *Mirage) ServeHTTPWithPort(w http.ResponseWriter, req *http.Request, port int) {
	host := strings.ToLower(strings.Split(req.Host, ":")[0])

	switch {
	case m.isWebApiHost(host):
		m.WebApi.ServeHTTP(w, req)

	case m.isDockerHost(host):
		m.ReverseProxy.ServeHTTPWithPort(w, req, port)

	default:
		// return 404
		http.NotFound(w, req)
	}

}

func (m *Mirage) isDockerHost(host string) bool {
	if strings.HasSuffix(host, m.Config.ReverseProxyHostSuffix) {
		ms := NewMirageStorage()
		subdomainList, err := ms.GetSubdomainList()
		if err != nil {
			return false
		}

		subdomain := strings.ToLower(strings.Split(host, ".")[0])
		index := sort.StringSlice(subdomainList).Search(subdomain)
		if index < len(subdomainList) && subdomainList[index] == subdomain {
			// found
			return true
		}

		return false
	} 

	return false
}

func (m *Mirage) isWebApiHost(host string) bool {
	return isSameHost(m.Config.WebApiHost, host)
}

func isSameHost(s1 string, s2 string) bool {
	lower1 := strings.Trim(strings.ToLower(s1), " ")
	lower2 := strings.Trim(strings.ToLower(s2), " ")

	return lower1 == lower2
}


