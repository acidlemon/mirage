package webapi

import (
	"net/http"
	"fmt"

	"github.com/acidlemon/mirage/docker"

	"github.com/acidlemon/rocket"
)


type WebApi struct {
	rocket.WebApp
	Host string
}


func NewWebApi(host string) *WebApi {
	app := &WebApi{}
	app.Init()
	app.Host = host

	app.AddRoute("/", app.List, &rocket.View{})
	app.AddRoute("/api/list", app.ApiList, &rocket.View{})
	app.AddRoute("/api/launch", app.ApiLaunch, &rocket.View{})
	app.AddRoute("/api/terminate", app.ApiTerminate, &rocket.View{})

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

func (app *WebApi) ApiList(c rocket.CtxData) {
	info, err := docker.List()
	var status interface{}
	if err != nil {
		status = err.Error()
	} else {
		status = info
	}

	result := rocket.RenderVars {
		"result": status,
	}

	c.RenderJSON(result)
}

func (app *WebApi) ApiLaunch(c rocket.CtxData) {
	if c.Req().Method != "POST" {
		c.Res().StatusCode = http.StatusMethodNotAllowed
		c.RenderText("you must use POST")
		return
	}

	subdomain, _ := c.ParamSingle("subdomain")
	branch   , _ := c.ParamSingle("branch")
	image    , _ := c.ParamSingle("image")

	status := "ok"

	if subdomain == "" || branch == "" || image == "" {
		status = fmt.Sprintf("parameter required: subdomain=%s, branch=%s, image=%s",
			subdomain, branch, image)
	} else {
		err := docker.Launch(subdomain, branch, image)
		if err != nil {
			status = err.Error()
		}
	}

	result := rocket.RenderVars {
		"result": status,
	}

	c.RenderJSON(result)
}

func (app *WebApi) ApiTerminate(c rocket.CtxData) {
	if c.Req().Method != "POST" {
		c.Res().StatusCode = http.StatusMethodNotAllowed
		c.RenderText("you must use POST")
		return
	}

	status := "ok"

	subdomain, _ := c.ParamSingle("subdomain")
	if subdomain == "" {
		status = fmt.Sprintf("parameter required: subdomain")
	} else {
		err := docker.Terminate(subdomain)

		if err != nil {
			status = err.Error()
		}
	}

	result := rocket.RenderVars {
		"result": status,
	}

	c.RenderJSON(result)
}


