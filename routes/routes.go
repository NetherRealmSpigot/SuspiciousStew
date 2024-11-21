package routes

import (
	"github.com/gin-gonic/gin"
	"stew/router"
	"stew/routes/v1/gateway"
	"stew/routes/v1/network"
	"stew/types"
)

func loadRoutes(routeGroup *gin.RouterGroup, routes []types.APIRoute) {
	for _, route := range routes {
		routeGroup.Handle(route.Method, route.Path, route.Handlers...)
	}
}

// !!! INVOKE THIS AFTER LoadRouter !!!
func LoadRoutes(conf types.APIConfig) {
	loadRoutes(router.Router.Group(gateway.RouteGroup), gateway.Routes)
	loadRoutes(router.Router.Group(network.RouteGroup), network.Routes)
}
