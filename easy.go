package lorca

import (
	"fmt"
	"github.com/injoyai/conv"
	"io/ioutil"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

type Config struct {
	Width   int //宽
	Height  int //高
	Html    string
	Dir     string
	Options []string
}

func (this *Config) new() *Config {
	if this.Width == 0 {
		this.Width = 480
	}
	if this.Height == 0 {
		this.Height = 320
	}
	if strings.HasPrefix(this.Html, "./") || strings.HasPrefix(this.Html, "/") {
		bs, err := ioutil.ReadFile(this.Html)
		if err == nil {
			this.Html = string(bs)
		}
	}
	if len(this.Html) == 0 {
		this.Html = `
	<html>
		<head><title>Hello</title></head>
		<body><h1>Hello, world!</h1></body>
	</html>`
	}
	return this
}

func Run(cfg *Config, fn func(APP) error) error {
	cfg.new()
	ui, err := New("data:text/html,"+url.PathEscape(cfg.Html), cfg.Dir, cfg.Width, cfg.Height, cfg.Options...)
	if err != nil {
		return err
	}
	defer ui.Close()
	if err = fn(&app{ui}); err != nil {
		return err
	}
	sign := make(chan os.Signal)
	signal.Notify(sign, os.Interrupt, os.Kill, syscall.SIGTERM)
	select {
	case <-sign:
	case <-ui.Done():
	}
	return nil
}

type APP interface {
	UI
	GetValueByID(id string) string
	GetVarByID(id string) *conv.Var
	SetValueByID(id string, value interface{})
	SetInnerByID(id string, value interface{})
}

type app struct {
	UI
}

func (this *app) GetValueByID(id string) string {
	js := fmt.Sprintf("document.getElementById('%s').value", id)
	return this.Eval(js).String()
}

func (this *app) GetVarByID(id string) *conv.Var {
	return conv.New(this.GetValueByID(id))
}

func (this *app) SetValueByID(id string, value interface{}) {
	js := fmt.Sprintf("document.getElementById('%s').value='%s'", id, conv.String(value))
	this.Eval(js)
}

func (this *app) SetInnerByID(id string, value interface{}) {
	js := fmt.Sprintf("document.getElementById('%s').innerHTML='%s'", id, conv.String(value))
	this.Eval(js)
}
