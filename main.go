package main

import (
	_ "embed"
	"github.com/joho/godotenv"
	"stew/config"
	"stew/database"
	"stew/embeds"
	"stew/logging"
	"stew/router"
	"stew/routes"
	"stew/utils"
)

//go:embed sql/player_stats.sql
var playerStatsSQLScript string

//go:embed sql/accounts.sql
var accountsSQLScript string

func main() {
	embeds.InitDBEmbed(playerStatsSQLScript, accountsSQLScript)

	logging.LoadLogger()

	utils.Info()

	godotenv.Load()
	logging.AppLogger.Info("Loading config")
	dbConf, apiConf := config.LoadConfig()

	logging.AppLogger.Info("Loading database pool")
	db := database.LoadDatabase(dbConf)
	database.ConnectDatabase(db)

	logging.AppLogger.Info("Loading router")
	router.LoadRouter(apiConf)
	routes.LoadRoutes(apiConf)

	logging.AppLogger.Info("Starting server")
	listener := router.Serve(apiConf)

	logging.AppLogger.Info("Shutting down!")
	defer listener.Close()
	defer db.Close()
}
