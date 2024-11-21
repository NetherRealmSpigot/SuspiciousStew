package gateway

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"stew/router"
	"stew/routes/utils"
	"stew/types"
)

const RouteGroup = router.V1RootRouteGroup + "/gateway"

// uuid, name, version, id
func playerInfoValidator(ctx *gin.Context, path string, required ...bool) []string {
	func0 := utils.GetFormData

	if path == PlayerInfoPath && (ctx.Request.Method == http.MethodGet || ctx.Request.Method == http.MethodPatch) {
		func0 = utils.GetQueryData
	}
	if path == SessionPath && ctx.Request.Method == http.MethodGet {
		func0 = utils.GetQueryData
	}

	fields := []types.UnvalidatedField{
		{"uuid", func0, utils.ValidateUUID, true, true},
		{"name", utils.GetFormData, utils.ValidateIgn, true, true},
		{"version", utils.GetFormData, utils.ValidateVersion, true, true},
	}
	if path == SessionPath && ctx.Request.Method == http.MethodPost {
		fields = append(fields, types.UnvalidatedField{Name: "id", Getter: utils.GetFormData, Validator: utils.ValidateID, Required: true, AllowEmpty: true})
	}

	for i, _ := range fields {
		fields[i].AllowEmpty = !required[i]
		fields[i].Required = required[i]
	}

	if !utils.ValidateAllData(fields, ctx, false) {
		return nil
	}

	var res []string
	for _, field := range fields {
		res = append(res, field.Getter(field.Name, ctx))
	}
	return res
}

// playerUUID, ipId
func playerLoginInfoValidator(ctx *gin.Context) []string {
	fields := []types.UnvalidatedField{
		{"uuid", utils.GetFormData, utils.ValidateUUID, true, false},
		{"ipid", utils.GetFormData, utils.ValidateID, true, false},
	}

	if !utils.ValidateAllData(fields, ctx, false) {
		return nil
	}

	var res []string
	for _, field := range fields {
		res = append(res, field.Getter(field.Name, ctx))
	}
	return res
}

func ipInfoValidator(ctx *gin.Context) string {
	func0 := utils.GetFormData
	if ctx.Request.Method == http.MethodGet {
		func0 = utils.GetQueryData
	}

	fields := []types.UnvalidatedField{
		{"ip", func0, utils.ValidateIPv4, true, false},
	}

	if !utils.ValidateAllData(fields, ctx, false) {
		return ""
	}

	return fields[0].Getter(fields[0].Name, ctx)
}

var Routes = []types.APIRoute{
	{"", http.MethodGet, []gin.HandlerFunc{
		func(ctx *gin.Context) {
			ctx.Status(http.StatusNotFound)
		},
	}},
	{PlayerInfoPath, http.MethodGet, []gin.HandlerFunc{
		func(ctx *gin.Context) {
			res := playerInfoValidator(ctx, PlayerInfoPath, true, false, false)
			if res != nil {
				getPlayerInfo(res[0], ctx)
			}
		},
	}},
	{PlayerInfoPath, http.MethodPost, []gin.HandlerFunc{
		func(ctx *gin.Context) {
			res := playerInfoValidator(ctx, PlayerInfoPath, true, true, true)
			if res != nil {
				addPlayerInfo(res[0], res[1], res[2], ctx)
			}
		},
	}},
	{PlayerInfoPath, http.MethodPatch, []gin.HandlerFunc{
		func(ctx *gin.Context) {
			res := playerInfoValidator(ctx, PlayerInfoPath, true, true, true)
			if res != nil {
				updatePlayerInfo(res[0], res[1], res[2], ctx)
			}
		},
	}},
	{IpInfoPath, http.MethodGet, []gin.HandlerFunc{
		func(ctx *gin.Context) {
			res := ipInfoValidator(ctx)
			if res != "" {
				getIpInfo(res, ctx)
			}
		},
	}},
	{IpInfoPath, http.MethodPost, []gin.HandlerFunc{
		func(ctx *gin.Context) {
			res := ipInfoValidator(ctx)
			if res != "" {
				addIpInfo(res, ctx)
			}
		},
	}},
	{PlayerLoginPath, http.MethodPost, []gin.HandlerFunc{
		func(ctx *gin.Context) {
			res := playerLoginInfoValidator(ctx)
			if res != nil {
				handlePlayerLogin(res[0], res[1], ctx)
			}
		},
	}},
	{SessionPath, http.MethodGet, []gin.HandlerFunc{
		func(ctx *gin.Context) {
			res := playerInfoValidator(ctx, SessionPath, true, false, false)
			if res != nil {
				getSessionId(res[0], ctx)
			}
		},
	}},
	{SessionPath, http.MethodPost, []gin.HandlerFunc{
		func(ctx *gin.Context) {
			res := playerInfoValidator(ctx, SessionPath, false, false, false, true)
			if res != nil {
				updateLoginSession(res[3], ctx)
			}
		},
	}},
}
