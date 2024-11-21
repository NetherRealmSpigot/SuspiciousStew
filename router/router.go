package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net"
	"stew/types"
)

func Serve(conf types.APIConfig) net.Listener {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", conf.ListenAddress, conf.ListenPort))
	if err != nil {
		panic(err)
	}
	err = Router.RunListener(listener)
	if err != nil {
		panic(err)
	}
	ListenAddr = listener.Addr().(*net.TCPAddr).IP.String()
	ListenPort = listener.Addr().(*net.TCPAddr).Port
	return listener
}

var Router *gin.Engine
var ginMode = "debug"

func LoadRouter(conf types.APIConfig) {
	if Router == nil {
		gin.SetMode(ginMode)
		Router = gin.New()
	}
	Router.Use(addRequestIdHeader)
	Router.Use(routerLogger, gin.Recovery())
}
