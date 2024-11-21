package gateway

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"stew/database"
	"stew/logging"
	"stew/types"
)

func addIpInfo(ipString string, c *gin.Context) {
	ctx, cancel := database.SetTimeout(3)
	defer cancel()

	_, err := database.Pool.Exec(ctx, "SELECT stew_player_stats.add_ip_info($1);", ipString)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		logging.AppLogger.WithError(err).Error("Error adding ip info!!!")
		return
	}

	c.Status(http.StatusNoContent)
}

func getIpInfo(ipString string, c *gin.Context) {
	ctx, cancel := database.SetTimeout(3)
	defer cancel()

	exec, err := database.Pool.Query(ctx, "SELECT * FROM stew_player_stats.get_ip_info($1);", ipString)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		logging.AppLogger.WithError(err).Error("Error getting ip info!!!")
		return
	}
	defer exec.Close()

	res := make([]types.IpInfoResponse, 0)
	for exec.Next() {
		var i types.IpInfoResponse
		err = exec.Scan(&i.Id, nil)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			logging.AppLogger.WithError(err).Error("Error forging ip info response!!!")
			return
		}
		res = append(res, i)
	}
	c.JSON(http.StatusOK, res)
}

const IpInfoPath = "/ip"
