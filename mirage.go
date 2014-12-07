package main

import (
	"fmt"
	"net"
	"net/http"
	"sort"
	"strings"
	"sync"
)

var app *Mirage

type Mirage struct {
	Config       *Config
	WebApi       *WebApi
	ReverseProxy *ReverseProxy
	Docker       *Docker
	Storage      *MirageStorage
}

func Setup(cfg *Config) {
	ms := NewMirageStorage()
	m := &Mirage{
		Config:       cfg,
		WebApi:       NewWebApi(cfg),
		ReverseProxy: NewReverseProxy(cfg),
		Docker:       NewDocker(cfg, ms),
		Storage:      ms,
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
	for _, v := range app.Config.Listen.HTTP {
		wg.Add(1)
		go func(port int) {
			defer wg.Done()
			laddr := fmt.Sprintf("%s:%d", app.Config.Listen.ForeignAddress, port)
			listener, err := net.Listen("tcp", laddr)
			if err != nil {
				fmt.Println("cannot listen %s", laddr)
				return
			}

			mux := http.NewServeMux()
			mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
				app.ServeHTTPWithPort(w, req, port)
			})

			fmt.Println("listen port:", port)
			http.Serve(listener, mux)
		}(v.ListenPort)
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
	if strings.HasSuffix(host, m.Config.Host.ReverseProxySuffix) {
		ms := m.Storage
		subdomainList, err := ms.GetSubdomainList()
		if err != nil {
			return false
		}

		subdomain := strings.ToLower(strings.Split(host, ".")[0])
		sortedList := sort.StringSlice(subdomainList)
		sortedList.Sort()
		index := sortedList.Search(subdomain)
		if index < len(sortedList) && sortedList[index] == subdomain {
			// found
			return true
		}

		return false
	}

	return false
}

func (m *Mirage) isWebApiHost(host string) bool {
	return isSameHost(m.Config.Host.WebApi, host)
}

func isSameHost(s1 string, s2 string) bool {
	lower1 := strings.Trim(strings.ToLower(s1), " ")
	lower2 := strings.Trim(strings.ToLower(s2), " ")

	return lower1 == lower2
}
