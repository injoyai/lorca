package main

import "github.com/injoyai/lorca"

func main() {
	lorca.Run(&lorca.Config{
		Source: "http://www.baidu.com",
	})
}
