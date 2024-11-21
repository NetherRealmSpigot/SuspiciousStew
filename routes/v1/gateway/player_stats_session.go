package gateway

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"stew/database"
	"stew/logging"
	"stew/types"
)

func updateLoginSession(id string, c *gin.Context) {
	ctx, cancel := database.SetTimeout(3)
	defer cancel()

	_, err := database.Pool.Exec(ctx, "SELECT stew_player_stats.update_login_session($1);", id)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		logging.AppLogger.WithError(err).Error("Error updating login session!!!")
		return
	}

	c.Status(http.StatusNoContent)
}

func getSessionId(uuid string, c *gin.Context) {
	ctx, cancel := database.SetTimeout(3)
	defer cancel()

	res := types.SessionIdResponse{}
	exec, err := database.Pool.Query(ctx, "SELECT * FROM stew_player_stats.get_session_id($1);", uuid)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		logging.AppLogger.WithError(err).Error("Error getting login session id!!!")
		return
	}
	defer exec.Close()

	if exec.Next() {
		err = exec.Scan(&res.Id, nil, nil, nil)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			logging.AppLogger.WithError(err).Error("Error forging login session id response!!!")
			return
		}
	}
	c.JSON(http.StatusOK, res)
}

const SessionPath = "/session"
