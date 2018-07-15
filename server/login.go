package server

import (
	"encoding/json"
	"io/ioutil"

	"gopkg.in/macaron.v1"
)

type Login struct {
	Name     string
	Password string
}

func (app *App) Login(ctx *macaron.Context) {
	l := &Login{}
	data, err := ioutil.ReadAll(ctx.Req.Body().ReadCloser())
	if err != nil {
		ctx.Error(403, err.Error())
	}

	err = json.Unmarshal(data, l)
	if err != nil {
		ctx.Error(403, err.Error())
	}

	if l.Name == "halulu" && l.Password == "websocks" {
		session, _ := app.store.Get(ctx.Req.Request, "cookie")
		session.Values["authenticated"] = true
		session.Save(ctx.Req.Request, ctx)
		ctx.HTML(200, "ok")
		return
	}

	ctx.Error(403, "not halulu")
}
