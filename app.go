package lorca

import (
	"errors"
	"fmt"
	"github.com/injoyai/conv"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

type Option func(app APP) error

func WithIndex(index string) Option {
	return func(app APP) error {
		return app.SwitchPage(index)
	}
}

type Config struct {
	Width   int             //宽
	Height  int             //高
	Options []string        //可选参数
	Index   string          //首页,可以是网址,本地html文件,或者html内容
	Pages   map[string]Page //多页面
}

func (this *Config) init() *Config {
	if this.Width == 0 {
		this.Width = 480
	}
	if this.Height == 0 {
		this.Height = 320
	}
	this.Index = dealSource(this.Index)
	return this
}

func Run(cfg *Config, op ...Option) error {
	cfg.init()
	//"data:text/html,"+ url.PathEscape(this.Source)
	ui, err := New(cfg.Index, "", cfg.Width, cfg.Height, cfg.Options...)
	if err != nil {
		return err
	}
	defer ui.Close()
	_app := &app{UI: ui, pages: cfg.Pages}
	for _, v := range op {
		if v != nil {
			if err = v(_app); err != nil {
				return err
			}
		}
	}
	if len(op) == 0 && len(cfg.Index) == 0 && len(_app.pages) > 0 && _app.pages["index"] != nil {
		_app.SwitchPage("index")
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
	Alert(msg interface{})
	SetPage(name string, source string, op ...Option)
	SwitchPage(name string) error
	GetByName(name string, nature string) string
	GetByID(id string, nature string) string
	GetVarByID(id string, nature string) *conv.Var
	GetValueByID(id string) string
	SetValueByID(id string, value interface{})
	SetInnerByID(id string, value interface{})
}

type app struct {
	UI
	pages map[string]Page
}

func (this *app) Alert(msg interface{}) {
	this.Eval("alert('" + conv.String(msg) + "')")
}

func (this *app) SetPage(name string, source string, op ...Option) {
	if this.pages == nil {
		this.pages = make(map[string]Page)
	}
	this.pages[name] = NewPage(source, op...)
}

func (this *app) SwitchPage(name string) error {
	if this.pages == nil {
		return errors.New("page not found")
	}
	p, ok := this.pages[name]
	if ok && p != nil {
		return p.Switch(this)
	}
	return errors.New("page not found")
}

func (this *app) GetByID(id string, nature string) string {
	js := fmt.Sprintf("document.getElementById('%s').%s", id, nature)
	return this.Eval(js).String()
}

func (this *app) GetByName(name string, nature string) string {
	js := fmt.Sprintf("document.getElementByName('%s').%s", name, nature)
	return this.Eval(js).String()
}

func (this *app) GetVarByID(id string, nature string) *conv.Var {
	js := fmt.Sprintf("document.getElementById('%s').%s", id, nature)
	return conv.New(this.Eval(js).String())
}

func (this *app) GetValueByID(id string) string {
	js := fmt.Sprintf("document.getElementById('%s').value", id)
	return this.Eval(js).String()
}

func (this *app) SetValueByID(id string, value interface{}) {
	js := fmt.Sprintf("document.getElementById('%s').value='%s'", id, conv.String(value))
	this.Eval(js)
}

func (this *app) SetInnerByID(id string, value interface{}) {
	js := fmt.Sprintf("document.getElementById('%s').innerHTML='%s'", id, conv.String(value))
	this.Eval(js)
}

func dealSource(s string) string {
	if len(s) == 0 {
		//空白内容

	} else if strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://") {
		//是个网址

	} else if stat, err := os.Stat(s); stat != nil && !os.IsNotExist(err) {
		if !stat.IsDir() {
			bs, _ := os.ReadFile(s)
			s = "data:text/html," + url.PathEscape(string(bs))
		} else {
			//文件夹
		}

	} else {
		s = "data:text/html," + url.PathEscape(s)

	}
	return s
}
