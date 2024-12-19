package main

import (
	_ "embed"
	"github.com/injoyai/lorca"
)

//go:embed login.html
var Login string

//go:embed home.html
var Home string

func main() {
	lorca.Run(&lorca.Config{
		Width:  800,
		Height: 600,
		Pages: map[string]lorca.Page{
			"index": lorca.NewPage(Login, func(app lorca.APP) error {
				return app.Bind("login", func() { app.SwitchPage("home") })
			}),
			"home": lorca.NewPage(Home, func(app lorca.APP) error {
				return app.Bind("back", func() { app.SwitchPage("index") })
			}),
		},
	})
}
