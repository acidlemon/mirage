package main

import (
	"sync"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/acidlemon/mirage/reverseproxy"
	"github.com/acidlemon/mirage/webapi"
	_ "github.com/acidlemon/mirage/docker"

//	"github.com/acidlemon/go-dumper"
)


func main() {
	fmt.Println("Launch succeeded!")

	foreignHost := "127.0.0.1"
	ports := []int{ 8080, 8443 }

	mirage := NewMirage()

	var wg sync.WaitGroup
	for _, v := range ports {
		wg.Add(1)
		go func(port int) {
			defer wg.Done()
			laddr := fmt.Sprintf("%s:%d", foreignHost, port)
			listener, err := net.Listen("tcp", laddr)
			if err != nil {
				fmt.Errorf("cannot listen %s", laddr)
				return
			}

			// TODO: SSL Support
			http.Serve(listener, mirage)
		}(v)
	}

	wg.Wait()
}


type Mirage struct {
	WebApi *webapi.WebApi
	notFound http.Handler
}

func NewMirage() *Mirage {
	mirage := &Mirage{
		WebApi: webapi.NewWebApi("localhost"),
		notFound: http.NotFoundHandler(),
	}
	reverseproxy.SetHostSuffix(".example.net")

	return mirage
}

func (app *Mirage) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	host := strings.ToLower(strings.Split(req.Host, ":")[0])

	switch {
	case app.isDockerHost(host):
		fmt.Println("this is revproxy host")
		reverseproxy.ServeHTTP(w, req)

	case app.isWebApiHost(host):
		fmt.Println("this is webapi host")
		app.WebApi.ServeHTTP(w, req)

	default:
		// return 404
		app.notFound.ServeHTTP(w, req)
	}

}

func (app *Mirage) isDockerHost(host string) bool {
	if strings.HasSuffix(host, reverseproxy.DefaultReverseProxy.HostSuffix) {
		// TODO search docker name

		return true
	} 

	return false
}

func (app *Mirage) isWebApiHost(host string) bool {
	return isSameHost(app.WebApi.Host, host)
}

func isSameHost(s1 string, s2 string) bool {
	lower1 := strings.Trim(strings.ToLower(s1), " ")
	lower2 := strings.Trim(strings.ToLower(s2), " ")

	return lower1 == lower2
}


