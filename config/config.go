package config

import (
	"stew/embeds"
	"stew/types"
	"strings"
)

func LoadConfig() (types.DatabaseConfig, types.APIConfig) {
	var db types.DatabaseConfig
	var api types.APIConfig
	db.SQLHost = readStr(key("SQL_HOST"), "127.0.0.1")
	db.SQLPort = readUInt16(key("SQL_PORT"), 5432)
	if db.SQLPort <= 1024 || db.SQLPort > 65535 {
		panic("Illegal port number.")
	}

	db.SQLUsername = readStr(key("SQL_USERNAME"), "postgres")
	db.SQLPassword = readStr(key("SQL_PASSWORD"), "")
	db.SQLDatabase = readStr(key("SQL_DATABASE"), strings.ToLower(embeds.Code))
	db.MinConns = readInt32(key("SQL_MIN_CONNS"), 1)
	db.MaxConns = readInt32(key("SQL_MAX_CONNS"), 10)
	if db.MaxConns < db.MinConns || db.MinConns <= 0 || db.MaxConns <= 0 {
		panic("Illegal number of min/max SQL connections.")
	}

	api.ListenAddress = readStr(key("LISTEN_ADDRESS"), "127.0.0.1")
	api.ListenPort = readUInt16(key("LISTEN_PORT"), 8080)
	if api.ListenPort <= 1024 || api.ListenPort > 65535 {
		panic("Illegal port number.")
	}

	return db, api
}
