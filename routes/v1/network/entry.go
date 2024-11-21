package network

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"stew/router"
	"stew/types"
)

const RouteGroup = router.V1RootRouteGroup + "/network"

var Routes = []types.APIRoute{
	{"", http.MethodGet, []gin.HandlerFunc{
		func(ctx *gin.Context) {
			ctx.Status(http.StatusNotFound)
		},
	}},
}
