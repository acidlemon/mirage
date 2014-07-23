package webapi

import (
	"net/http"

	"github.com/acidlemon/rocket"
)


type WebApi struct {
	rocket.WebApp
	Host string
}


func NewWebApi(host string) (WebApi) {
	app := WebApi{}
	app.Init()
	app.Host = host

	app.AddRoute("/", app.List, &rocket.View{})

	app.BuildRouter()

	return app
}

func (app *WebApi) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	app.Handler(w, req)
}

func (app *WebApi) List(c rocket.CtxData) {
	c.Res().StatusCode = http.StatusOK
	value := rocket.RenderVars {
		"test" : "powawa",
	}

	c.Render("webapi/list.html", value)
}


