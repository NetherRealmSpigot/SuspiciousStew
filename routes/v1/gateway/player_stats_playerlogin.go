package gateway

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"stew/database"
	"stew/logging"
)

func handlePlayerLogin(playerUUID string, ipId string, c *gin.Context) {
	ctx, cancel := database.SetTimeout(3)
	defer cancel()

	_, err := database.Pool.Exec(ctx, "SELECT stew_player_stats.handle_player_logins($1, $2);", playerUUID, ipId)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		logging.AppLogger.WithError(err).Error("Error handling player login!!!")
		return
	}

	c.Status(http.StatusNoContent)
}

const PlayerLoginPath = PlayerInfoPath + "/login"
