package main

import (
	"net/http"
	"fmt"

	"github.com/acidlemon/rocket"
)


type WebApi struct {
	rocket.WebApp
	cfg *Config
}


func NewWebApi(cfg *Config) *WebApi {
	app := &WebApi{}
	app.Init()
	app.cfg = cfg

	view := &rocket.View{
		BasicTemplates: []string{"html/layout.html"},
	}

	app.AddRoute("/", app.List, view)
	app.AddRoute("/launcher", app.Launcher, view)
	app.AddRoute("/launch", app.Launch, view)
	app.AddRoute("/terminate", app.Terminate, view)
	app.AddRoute("/api/list", app.ApiList, view)
	app.AddRoute("/api/launch", app.ApiLaunch, view)
	app.AddRoute("/api/terminate", app.ApiTerminate, view)

	app.BuildRouter()

	return app
}

func (api *WebApi) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	api.Handler(w, req)
}

func (api *WebApi) List(c rocket.CtxData) {
	info, err := app.Docker.List()
	errStr := ""
	if err != nil {
		errStr = err.Error()
	}
	value := rocket.RenderVars {
		"info" : info,
		"error": errStr,
	}

	c.Render("html/list.html", value)
}

func (api *WebApi) Launcher(c rocket.CtxData) {
	c.Render("html/launcher.html", rocket.RenderVars{
		"DefaultImage": api.cfg.Docker.DefaultImage,
	})
}

func (api *WebApi) Launch(c rocket.CtxData) {
	result := api.launch(c)
	if result["result"] == "ok" {
		c.Redirect("/")
	} else {
		c.RenderJSON(result)
	}
}

func (api *WebApi) Terminate(c rocket.CtxData) {
	result := api.terminate(c)
	if result["result"] == "ok" {
		c.Redirect("/")
	} else {
		c.RenderJSON(result)
	}
}

func (api *WebApi) ApiList(c rocket.CtxData) {
	info, err := app.Docker.List()
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

func (api *WebApi) ApiLaunch(c rocket.CtxData) {
	result := api.launch(c)

	c.RenderJSON(result)
}

func (api *WebApi) ApiTerminate(c rocket.CtxData) {
	result := api.terminate(c)

	c.RenderJSON(result)
}

func (api *WebApi) launch(c rocket.CtxData) rocket.RenderVars {
	if c.Req().Method != "POST" {
		c.Res().StatusCode = http.StatusMethodNotAllowed
		c.RenderText("you must use POST")
		return rocket.RenderVars{}
	}

	subdomain, _ := c.ParamSingle("subdomain")
	branch   , _ := c.ParamSingle("branch")
	image    , _ := c.ParamSingle("image")

	status := "ok"

	if subdomain == "" || branch == "" || image == "" {
		status = fmt.Sprintf("parameter required: subdomain=%s, branch=%s, image=%s",
			subdomain, branch, image)
	} else {
		err := app.Docker.Launch(subdomain, branch, image)
		if err != nil {
			status = err.Error()
		}
	}

	result := rocket.RenderVars {
		"result": status,
	}

	return result
}

func (api *WebApi) terminate(c rocket.CtxData) rocket.RenderVars {
	if c.Req().Method != "POST" {
		c.Res().StatusCode = http.StatusMethodNotAllowed
		c.RenderText("you must use POST")
		return rocket.RenderVars{}
	}

	status := "ok"

	subdomain, _ := c.ParamSingle("subdomain")
	if subdomain == "" {
		status = fmt.Sprintf("parameter required: subdomain")
	} else {
		err := app.Docker.Terminate(subdomain)

		if err != nil {
			status = err.Error()
		}
	}

	result := rocket.RenderVars {
		"result": status,
	}

	return result
}

