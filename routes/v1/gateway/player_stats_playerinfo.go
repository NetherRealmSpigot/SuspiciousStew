package gateway

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"stew/database"
	"stew/logging"
	"stew/types"
)

func addPlayerInfo(uuid string, name string, version string, c *gin.Context) {
	ctx, cancel := database.SetTimeout(3)
	defer cancel()

	_, err := database.Pool.Exec(ctx, "SELECT stew_player_stats.add_player_info($1, $2, $3);", uuid, name, version)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		logging.AppLogger.WithError(err).Error("Error adding player info!!!")
		return
	}

	c.Status(http.StatusNoContent)
}

func getPlayerInfo(uuid string, c *gin.Context) {
	ctx, cancel := database.SetTimeout(3)
	defer cancel()

	res := types.PlayerInfoResponse{}
	exec, err := database.Pool.Query(ctx, "SELECT * FROM stew_player_stats.get_player_info($1);", uuid)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		logging.AppLogger.WithError(err).Error("Error getting player info!!!")
		return
	}
	defer exec.Close()

	if exec.Next() {
		err = exec.Scan(&res.UUID, &res.Name, &res.Version)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			logging.AppLogger.WithError(err).Error("Error forging player info response!!!")
			return
		}
	}
	c.JSON(http.StatusOK, res)
}

func updatePlayerInfo(uuid string, name string, version string, c *gin.Context) {
	ctx, cancel := database.SetTimeout(3)
	defer cancel()

	_, err := database.Pool.Exec(ctx, "SELECT stew_player_stats.update_player_info($1, $2, $3);", uuid, name, version)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		logging.AppLogger.WithError(err).Error("Error updating player info!!!")
		return
	}

	c.Status(http.StatusNoContent)
}

const PlayerInfoPath = "/player"
