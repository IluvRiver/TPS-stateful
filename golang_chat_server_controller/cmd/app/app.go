package app

import "golang_chat_server_controller/config"

type App struct {
	cfg *config.Config
}

func NewApp(cfg *config.Config) *App {
	a := &App{cfg: cfg}

	return a
}
