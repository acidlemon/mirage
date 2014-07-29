package mirage

import (
	"sync"
	"fmt"
	"net"
	"net/http"
	"strings"

//	"github.com/acidlemon/go-dumper"
)

var app *Mirage

type Config struct {
	ForeignHost string
	WebApiHost string
	ReverseProxyHostSuffix string
	ListenPorts map[int]int // map[ListenPort] = ProxyPassPort
	ListenSSLPorts map[int]int
	DockerEndpoint string
}

func NewConfig() *Config {
	// default config
	cfg := &Config{
		ForeignHost: "127.0.0.1",
		WebApiHost: "localhost",
		ReverseProxyHostSuffix: ".example.net",
		ListenPorts: map[int]int{ 8080: 8080 },
		ListenSSLPorts: map[int]int{},
		DockerEndpoint: "unix:///var/run/docker.sock",
	}

	return cfg
}


type Mirage struct {
	Config *Config
	WebApi *WebApi
	ReverseProxy *ReverseProxy
	Docker *Docker
	notFound http.Handler
}

func Setup(cfg *Config) {
	m := &Mirage{
		Config: cfg,
		WebApi: NewWebApi(cfg),
		ReverseProxy: NewReverseProxy(cfg),
		Docker: NewDocker(cfg),
		notFound: http.NotFoundHandler(),
	}

	infolist, err := m.Docker.List()
	if err != nil {
		fmt.Println("cannot initialize reverse proxy: ", err.Error())
	}

	for _, info := range infolist {
		app.ReverseProxy.AddSubdomain(info.SubDomain, info.IPAddress)
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
			laddr := fmt.Sprintf("%s:%d", app.Config.ForeignHost, port)
			listener, err := net.Listen("tcp", laddr)
			if err != nil {
				fmt.Println("cannot listen %s", laddr)
				return
			}
			http.Serve(listener, app)
		}(k)
	}

	// TODO SSL Support

	wg.Wait()
}

func (m *Mirage) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	host := strings.ToLower(strings.Split(req.Host, ":")[0])

	switch {
	case m.isDockerHost(host):
		m.ReverseProxy.ServeHTTP(w, req)

	case m.isWebApiHost(host):
		m.WebApi.ServeHTTP(w, req)

	default:
		// return 404
		m.notFound.ServeHTTP(w, req)
	}

}

func (m *Mirage) isDockerHost(host string) bool {
	if strings.HasSuffix(host, m.Config.ReverseProxyHostSuffix) {
		// TODO search docker name

		return true
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


