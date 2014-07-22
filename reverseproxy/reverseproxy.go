package reverseproxy

import (
	"net/http"
//	"net/http/httputil"
)


type ReverseProxy struct {
	HostSuffix string
}


func (r *ReverseProxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {


}



