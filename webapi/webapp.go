package webapi

import (
	"net/http"

	"github.com/acidlemon/rocket"
)


type WebApi struct {
	rocket.WebApp
	Host string
}


func NewWebApi() (WebApi) {
	app := WebApi{}
	app.Init()

	app.AddRoute("/", app.List, &rocket.View{})

	app.BuildRouter()

	return app
}

func (app *WebApi) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	
}

func (app *WebApi) List(c rocket.CtxData) {
	
}


