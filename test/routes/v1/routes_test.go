package v1

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"stew/config"
	"stew/database"
	"stew/logging"
	"stew/router"
	"stew/routes"
	"stew/routes/v1/gateway"
	"stew/routes/v1/network"
	"stew/utils"
	"testing"
)

func TestMain(m *testing.M) {
	logging.LoadLogger()
	logging.AppLogger.Level = logrus.DebugLevel
	utils.Info()
	logging.AppLogger.Info("Starting testing...")

	logging.AppLogger.Info("Creating config for testing")
	godotenv.Load(path.Join("..", "..", "..", ".env"))
	dbConf, apiConf := config.LoadConfig()

	logging.AppLogger.Info("Loading database")
	db := database.LoadDatabase(dbConf)
	database.ConnectDatabase(db)

	for _, sc := range []string{"player_stats.sql", "accounts.sql"} {
		logging.AppLogger.Info(fmt.Sprintf("Running SQL script %s", sc))
		fbytes, err := os.ReadFile(path.Join("..", "..", "..", "sql", sc))
		if err != nil {
			panic(err)
		}
		sqlScript := string(fbytes)
		ctx, cancel := database.SetTimeout(3)
		defer cancel()
		_, err = db.Exec(ctx, sqlScript)
		if err != nil {
			panic(err)
		}
	}

	logging.AppLogger.Info("Loading router")
	router.LoadRouter(apiConf)
	routes.LoadRoutes(apiConf)

	logging.AppLogger.Info("Starting server")
	ts := httptest.NewUnstartedServer(router.Router)
	ts.Config.Addr = "127.0.0.1:0"
	router.ListenAddr = "127.0.0.1"
	ts.Start()
	port := ts.Listener.Addr().(*net.TCPAddr).Port
	router.ListenPort = port

	logging.AppLogger.Info("Running tests")
	m.Run()

	logging.AppLogger.Info("Shutting down!")
	defer ts.Listener.Close()
	defer db.Close()
}

func testRequestIdHeader(t *testing.T, resp *http.Response) {
	head := resp.Header.Get(router.RequestIdHeader)
	require.NotEmpty(t, head)
	require.Len(t, head, 24)
	require.Regexp(t, "^[0-9a-zA-Z]+$", head)
	logging.AppLogger.Debug(fmt.Sprintf("%s %s", resp.Request.URL.String(), head))
}

func TestRequestIdHeader(t *testing.T) {
	paths := []string{
		router.V1RootRouteGroup,
		gateway.RouteGroup,
		network.RouteGroup,
	}

	for _, p := range paths {
		for _, methodFunc := range []func(string) (resp *http.Response, err error){http.Get, http.Head} {
			url := fmt.Sprintf("http://%s:%d%s", router.ListenAddr, router.ListenPort, p)
			t.Run(fmt.Sprintf("Request id header %s", url), func(tt *testing.T) {
				resp, err := methodFunc(url)
				require.NoError(tt, err)
				testRequestIdHeader(tt, resp)
			})
		}
	}
}
