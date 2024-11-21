package v1

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/url"
	"stew/router"
	"stew/routes/v1/gateway"
	"stew/types"
	"strconv"
	"strings"
	"testing"
)

type loginEntry struct {
	ipString string
	uuid     string
	name     string
	version  string
}

var validLoginEntries = []loginEntry{
	{"198.51.100.10", "7659CEDB-C9C1-4F28-B966-19823FD8666B", "R_enchaku", "47"},
	{"198.51.100.12", "a3af3bfa-b2ba-4acc-a5a3-693649ffbcbc", "Desirada", "47"},
	{"198.51.100.13", "0e96bfa8-0952-4087-92b6-8d8897f793cc", "n_ekoli", "107"},
}

var validAnotherLoginEntries = []loginEntry{
	{"198.51.100.10", "7659CEDB-C9C1-4F28-B966-19823FD8666B", "P_layOsu", "107"},
	{"198.51.100.12", "a3af3bfa-b2ba-4acc-a5a3-693649ffbcbc", "FoxyWinter", "107"},
	{"198.51.100.13", "0e96bfa8-0952-4087-92b6-8d8897f793cc", "N_ekolitv", "210"},
}

func handlePlayerLogin(t *testing.T, expectStatus int, uuid string, ipId string) {
	resp, err := http.PostForm(fmt.Sprintf("http://%s:%d%s",
		router.ListenAddr, router.ListenPort, gateway.RouteGroup+gateway.PlayerLoginPath),
		url.Values{
			"uuid": []string{uuid},
			"ipid": []string{ipId},
		})
	require.NoError(t, err)
	require.Equal(t, expectStatus, resp.StatusCode)
	defer resp.Body.Close()
}

func getSessionId(t *testing.T, expectStatus int, uuid string) *types.SessionIdResponse {
	resp, err := http.Get(fmt.Sprintf("http://%s:%d%s?uuid=%s",
		router.ListenAddr, router.ListenPort, gateway.RouteGroup+gateway.SessionPath, uuid))
	require.NoError(t, err)
	require.Equal(t, expectStatus, resp.StatusCode)
	defer resp.Body.Close()
	id := &types.SessionIdResponse{}
	if expectStatus == http.StatusOK {
		err = json.NewDecoder(resp.Body).Decode(&id)
		require.NoError(t, err)
		if id.Id <= 0 {
			id = nil
		}
	} else {
		id = nil
	}
	return id
}

func updateSession(t *testing.T, expectStatus int, sessionId string) {
	resp, err := http.PostForm(fmt.Sprintf("http://%s:%d%s",
		router.ListenAddr, router.ListenPort, gateway.RouteGroup+gateway.SessionPath),
		url.Values{
			"id": []string{sessionId},
		})
	require.NoError(t, err)
	require.Equal(t, expectStatus, resp.StatusCode)
	defer resp.Body.Close()
}

func TestPlayerLogin(t *testing.T) {
	for _, entries := range [][]loginEntry{validLoginEntries, validAnotherLoginEntries} {
		t.Run("Login/Logout", func(_ *testing.T) {
			for _, entry := range entries {
				sessionId := int64(0)
				t.Run(fmt.Sprintf("Simulate login process %s %s %s %s", entry.ipString, entry.uuid, entry.name, entry.version), func(tt *testing.T) {
					ipAddress := entry.ipString
					uuid := entry.uuid
					name := entry.name
					version := entry.version
					vnum, _ := strconv.Atoi(version)
					addOrUpdatePlayer := false

					ips := getIpInfo(tt, http.StatusOK, ipAddress)
					require.GreaterOrEqual(tt, len(ips), 0)
					if len(ips) <= 0 {
						addIpInfo(tt, http.StatusNoContent, ipAddress)
						ips = getIpInfo(tt, http.StatusOK, ipAddress)
						require.GreaterOrEqual(tt, len(ips), 1)
					}
					player := getPlayerInfo(tt, http.StatusOK, uuid)

					if player == nil {
						addPlayerInfo(tt, http.StatusNoContent, uuid, name, version)
						player = getPlayerInfo(tt, http.StatusOK, uuid)
						require.NotNil(tt, player)
					}

					if player.Version != vnum {
						addOrUpdatePlayer = true
					} else if !strings.EqualFold(player.Name, name) {
						addOrUpdatePlayer = true
					}

					if addOrUpdatePlayer {
						updatePlayerInfo(tt, http.StatusNoContent, uuid, name, version)
					}
					player = getPlayerInfo(tt, http.StatusOK, uuid)
					require.NotNil(tt, player)
					require.True(tt, strings.EqualFold(player.Name, name))
					require.Equal(tt, player.Version, vnum)

					handlePlayerLogin(tt, http.StatusNoContent, uuid, strconv.FormatInt(ips[len(ips)-1].Id, 10))
				})
				t.Run(fmt.Sprintf("Simulate logout process %s %s %s %s", entry.ipString, entry.uuid, entry.name, entry.version), func(tt *testing.T) {
					uuid := entry.uuid

					sessionObj := getSessionId(tt, http.StatusOK, uuid)
					require.NotNil(tt, sessionObj)
					sessionId = sessionObj.Id

					updateSession(tt, http.StatusNoContent, strconv.FormatInt(sessionId, 10))
				})
			}
		})
	}
}
