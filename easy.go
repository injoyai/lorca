package lorca

import (
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

func Run(cfg *Config, fn func(UI) error) error {
	cfg.new()
	ui, err := New("data:text/html,"+url.PathEscape(cfg.Html), cfg.Dir, cfg.Width, cfg.Height, cfg.Options...)
	if err != nil {
		return err
	}
	defer ui.Close()
	if err = fn(ui); err != nil {
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
