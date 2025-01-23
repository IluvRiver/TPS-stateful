package main

import (
	"flag"
	"golang_chat_server_controller/cmd/app"
	"golang_chat_server_controller/config"
)

var pathFlag = flag.String("config", "./config.toml", "config set")

func main() {
	flag.Parse()
	c := config.NewConfig(*pathFlag)

	a := app.NewApp(c)
	a.Start()
}
