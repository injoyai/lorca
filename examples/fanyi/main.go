package main

import (
	"github.com/injoyai/lorca"
)

func main() {
	lorca.Run(&lorca.Config{
		Source: "./examples/fanyi/index.html",
	})
}
