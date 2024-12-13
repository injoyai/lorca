package lorca

import (
	"fmt"
	"github.com/injoyai/conv"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

type Config struct {
	Width   int      //宽
	Height  int      //高
	Source  string   //资源,可以是网址,本地html文件,或者html内容
	Options []string //可选参数
}

func (this *Config) init() *Config {
	if this.Width == 0 {
		this.Width = 480
	}
	if this.Height == 0 {
		this.Height = 320
	}

	if len(this.Source) == 0 {
		//空白内容

	} else if strings.HasPrefix(this.Source, "http://") || strings.HasPrefix(this.Source, "https://") {
		//是个网址

	} else if stat, err := os.Stat(this.Source); stat != nil && !os.IsNotExist(err) {
		if !stat.IsDir() {
			bs, _ := os.ReadFile(this.Source)
			this.Source = "data:text/html," + url.PathEscape(string(bs))
		} else {
			//文件夹
		}

	} else {
		this.Source = "data:text/html," + url.PathEscape(this.Source)

	}

	return this
}

func Run(cfg *Config, fn ...func(APP) error) error {
	cfg.init()
	//"data:text/html,"+ url.PathEscape(this.Source)
	ui, err := New(cfg.Source, "", cfg.Width, cfg.Height, cfg.Options...)
	if err != nil {
		return err
	}
	defer ui.Close()
	for _, v := range fn {
		if v != nil {
			if err = v(&app{ui}); err != nil {
				return err
			}
		}
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
	GetByID(id string, nature string) string
	GetVarByID(id string, nature string) *conv.Var
	GetValueByID(id string) string
	SetValueByID(id string, value interface{})
	SetInnerByID(id string, value interface{})
}

type app struct {
	UI
}

func (this *app) GetByID(id string, nature string) string {
	js := fmt.Sprintf("document.getElementById('%s').%s", id, nature)
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
