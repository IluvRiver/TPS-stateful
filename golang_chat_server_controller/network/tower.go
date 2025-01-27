package network

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type tower struct {
	server *Server
}

func registerTowerAPI(server *Server) {
	t := &tower{server: server}

	t.server.engine.GET("/server-list", t.serverList)
}

func (t *tower) serverList(c *gin.Context) {
	response(c, http.StatusOK, t.server.service.GetAvgServerList())
}
