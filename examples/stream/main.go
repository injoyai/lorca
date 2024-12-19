package main

import "github.com/injoyai/lorca"

func main() {
	lorca.Run(&lorca.Config{
		Index: "./examples/stream/index.html",
	})
}
