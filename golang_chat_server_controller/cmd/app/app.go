package app

import (
	"golang_chat_server_controller/config"
	"golang_chat_server_controller/network"
	"golang_chat_server_controller/repository"
	"golang_chat_server_controller/service"
)

type App struct {
	cfg *config.Config
	//repostory
	repository *repository.Repository
	//service
	service *service.Service
	//network를 컨트롤할거임
	network *network.Server
}

func NewApp(cfg *config.Config) *App {
	a := &App{cfg: cfg}

	var err error
	if a.repository, err = repository.NewRepository(cfg); err != nil {
		panic(err)
	} else {
		a.service = service.NewService(a.repository)
		a.network = network.NewNetwork(a.service, cfg.Info.Port)
	}

	return a
}
func (a *App) Start() error {
	return a.network.Start()
}
