package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/url"
	"stew/router"
	"stew/routes/v1/gateway"
	"stew/types"
	globalUtils "stew/utils"
	"strconv"
	"strings"
	"testing"
)

type infoEntry struct {
	uuid    string
	name    string
	version string
}

func addPlayerInfo(t *testing.T, expectStatus int, uuid string, name string, version string) {
	resp, err := http.PostForm(fmt.Sprintf("http://%s:%d%s",
		router.ListenAddr, router.ListenPort, gateway.RouteGroup+gateway.PlayerInfoPath),
		url.Values{
			"uuid":    []string{uuid},
			"name":    []string{name},
			"version": []string{version},
		})
	require.NoError(t, err)
	require.Equal(t, expectStatus, resp.StatusCode)
	defer resp.Body.Close()
}

func addPlayerInfoInvalid(t *testing.T, expectStatus int, contentType string, reader io.Reader) {
	resp, err := http.Post(fmt.Sprintf("http://%s:%d%s",
		router.ListenAddr, router.ListenPort, gateway.RouteGroup+gateway.PlayerInfoPath), contentType, reader)
	require.NoError(t, err)
	require.Equal(t, expectStatus, resp.StatusCode)
	defer resp.Body.Close()
}

func getPlayerInfo(t *testing.T, expectStatus int, uuid string) *types.PlayerInfoResponse {
	resp, err := http.Get(fmt.Sprintf("http://%s:%d%s?uuid=%s",
		router.ListenAddr, router.ListenPort, gateway.RouteGroup+gateway.PlayerInfoPath, uuid))
	require.NoError(t, err)
	require.Equal(t, expectStatus, resp.StatusCode)
	defer resp.Body.Close()
	player := &types.PlayerInfoResponse{}
	if expectStatus == http.StatusOK {
		err = json.NewDecoder(resp.Body).Decode(player)
		require.NoError(t, err)
		if !player.UUID.Valid {
			player = nil
		} else {
			require.True(t, strings.EqualFold(globalUtils.PGUUIDToString(player.UUID), uuid))
		}
	} else {
		player = nil
	}
	return player
}

func getPlayerInfoInvalid(t *testing.T, expectStatus int, field string, value string) {
	resp, err := http.Get(fmt.Sprintf("http://%s:%d%s?%s=%s",
		router.ListenAddr, router.ListenPort, gateway.RouteGroup+gateway.PlayerInfoPath, field, value))
	require.NoError(t, err)
	require.Equal(t, expectStatus, resp.StatusCode)
	defer resp.Body.Close()
}

func updatePlayerInfo(t *testing.T, expectStatus int, uuid string, name string, version string) {
	req, err := http.NewRequest(http.MethodPatch,
		fmt.Sprintf("http://%s:%d%s?uuid=%s",
			router.ListenAddr, router.ListenPort, gateway.RouteGroup+gateway.PlayerInfoPath, uuid),
		bytes.NewBufferString(url.Values{
			"name":    []string{name},
			"version": []string{version},
		}.Encode()))
	require.NoError(t, err)
	req.Header.Set("content-type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	require.Equal(t, expectStatus, resp.StatusCode)
	defer resp.Body.Close()
}

func updatePlayerInfoInvalid(t *testing.T, expectStatus int, uuidString string, contentType string, reader io.Reader) {
	req, err := http.NewRequest(http.MethodPatch,
		fmt.Sprintf("http://%s:%d%s?uuid=%s",
			router.ListenAddr, router.ListenPort, gateway.RouteGroup+gateway.PlayerInfoPath, uuidString), reader)
	require.NoError(t, err)
	req.Header.Set("content-type", contentType)
	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	require.Equal(t, expectStatus, resp.StatusCode)
	defer resp.Body.Close()
}

var validEntries = []infoEntry{
	{"D48ac2c8-0D91-4E90-9BF4-DC15E98C4B90", "Anime_Ban", "47"},
	{"7659CEDB-C9C1-4F28-B966-19823FD8666B", "P_l_a_y_O_s_u", "107"},
	{"dfca7025-ff60-47e3-a747-aba3a625e5ef", "M_i_c_k_y_", "210"},
	{"a99a0ff2-0060-4bfe-b84a-60f8efe288b2", "material", "315"},
	{"6db0eb9b-04f9-4777-bcaf-94bf53f42056", "outcube_d", "335"},
	{"0a463e19-7544-4d0d-9991-d3983214b565", "Winter_Slime3377", "393"},
	{"273d5a90-6b3d-4c6c-a69e-3259a00eb392", "frog_hat", "477"},
	{"22500b81-e889-4367-b83c-24c52914e2de", "kaci_nass", "573"},
}

var invalidEntriesValidUUID = []infoEntry{
	{"d48ac2c8-0d91-4e90-9bf4-dc15e98c4b90", "Anime_Ban", "GGGGGGGGGGGGGGGGggg42900f\x20\x20\x00\x00"},
	{"d48ac2c8-0d91-4e90-9bf4-dc15e98c4b90", "Anime_Ban", "46"},
	{"d48ac2c8-0d91-4e90-9bf4-dc15e98c4b90", "Anime-Ban", "48"},
	{"d48ac2c8-0d91-4e90-9bf4-dc15e98c4b90", "Anime_BanAnime_BanAnime_Ban", "47"},
	{"d48ac2c8-0d91-4e90-9bf4-dc15e98c4b90", "AnG!#v0s9d", "47"},
	{"d48ac2c8-0d91-4e90-9bf4-dc15e98c4b90", "AnG!#v0s9dAnG!#v0s9dAnG!#v0s9d\x20\x20\x20\x20\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00", "-99999999999999999999999999999999999999999999999999999999999999999999"},
	{"d48ac2c8-0d91-4e90-9bf4-dc15e98c4b90", "An--31", "47"},
	{"d48ac2c8-0d91-4e90-9bf4-dc15e98c4b90", "Anime-Ban", "770"},
	{"d48ac2c8-0d91-4e90-9bf4-dc15e98c4b90", "Anime-Ban", "481"},
	{"d48ac2c8-0d91-4e90-9bf4-dc15e98c4b90", "Anime-Ban", "999999999999999999999999999999999999999999999999999999999999999"},
}

var invalidEntriesInvalidUUID = []infoEntry{
	{"000000000000000000000000000000000000", "Anime-Ban", "46"},
	{"000000000000000000000000000000000000", "Anime-Ban", "47"},
	{"00000000000000000000000000000000", "Anime-Ban", "47"},
	{"00000000-0000-0000-0000-000000000000", "Anime-Ban", "47"},
	{"d48ac2c8-0d91-3e90-9bf4-dc15e98c4b90", "Anime-Ban", "47"},
	{"d48ac2c8-0d91-4e90-7bf4-dc15e98c4b90", "Anime-Ban", "47"},
	{"d48ac2c8-0d91-4e90-7bf4-dc15e98c4b90", "Anime-Ban", "481"},
	{"22500b81-e889-4367-b83c-d", "kacin-as", "573"},
	{"e889-4367-b83c-24c52914e2de", "kacin-as", "573"},
	{"", "", ""},
	{" WHERE 1=1 --", "WHERE 1=1 --", " WHERE 1=1 --"},
	{" ;;", "WHERE 1=1 --", ";;;;"},
}

var validUpdateEntries = []infoEntry{
	{"D48ac2c8-0D91-4E90-9BF4-DC15E98C4B90", "Anime_Ban", "47"},
	{"D48ac2c8-0D91-4E90-9BF4-DC15E98C4B90", "FUCKER", "107"},
	{"D48ac2c8-0D91-4E90-9BF4-DC15E98C4B90", "Anime_Ban", "573"},
}

var invalidUpdateEntriesValidUUID = []infoEntry{
	{"D48ac2c8-0D91-4E90-9BF4-DC15E98C4B90", "Anime-Ban", "999999999999999999999999999"},
	{"D48ac2c8-0D91-4E90-9BF4-DC15E98C4B90", "FdG4290t1i0249GGEGR031CKER", "107"},
	{"D48ac2c8-0D91-4E90-9BF4-DC15E98C4B90", "!!1309f%51", "-99999999999999999999999999999999999999999999999999900000000000000000000000009999"},
}

var invalidUpdateEntriesInvalidUUID = []infoEntry{
	{"D48ac2c8-0D91-3E90-9BF4-DC15E98C4B90", "Anime-Ban", "48"},
	{"D48ac2c8-0D91-4E90-7BF4-DC15E98C4B90", "FUCKER", "107"},
}

func TestAddPlayerInfo(t *testing.T) {
	for _, ei := range validEntries {
		t.Run(fmt.Sprintf("Add player info valid entry %s %s %s", ei.uuid, ei.name, ei.version), func(tt *testing.T) {
			addPlayerInfo(tt, http.StatusNoContent, ei.uuid, ei.name, ei.version)
		})
	}

	for _, ei := range invalidEntriesValidUUID {
		t.Run(fmt.Sprintf("Add player info invalid entry valid uuid %s %s %s", ei.uuid, ei.name, ei.version), func(tt *testing.T) {
			addPlayerInfo(tt, http.StatusBadRequest, ei.uuid, ei.name, ei.version)
		})
	}

	for _, ei := range invalidEntriesInvalidUUID {
		t.Run(fmt.Sprintf("Add player info invalid entry invalid uuid %s %s %s", ei.uuid, ei.name, ei.version), func(tt *testing.T) {
			addPlayerInfo(tt, http.StatusBadRequest, ei.uuid, ei.name, ei.version)
		})
	}
	for _, ent := range []struct {
		contentType string
		reader      io.Reader
	}{
		{contentType: "application/json", reader: nil},
		{contentType: "application/json", reader: bytes.NewBuffer([]byte("[]"))},
		{contentType: "application/json", reader: bytes.NewBuffer([]byte("{}"))},
		{contentType: "application/x-www-form-urlencoded", reader: bytes.NewBuffer([]byte("{}"))},
		{contentType: "application/x-www-form-urlencoded", reader: bytes.NewBuffer([]byte(""))},
		{contentType: "application/x-www-form-urlencoded", reader: bytes.NewBuffer([]byte("WHERE 1=1;--"))},
	} {
		t.Run("Add player info invalid request", func(tt *testing.T) {
			addPlayerInfoInvalid(tt, http.StatusBadRequest, ent.contentType, ent.reader)
		})
	}
}

func TestGetPlayerInfo(t *testing.T) {
	for _, u := range validEntries {
		t.Run(fmt.Sprintf("Get player info valid entry %s %s %s", u.uuid, u.name, u.version), func(tt *testing.T) {
			player := getPlayerInfo(tt, http.StatusOK, u.uuid)
			require.NotNil(tt, player)
			require.True(tt, strings.EqualFold(player.Name, u.name))
			v := strconv.Itoa(player.Version)
			us := globalUtils.PGUUIDToString(player.UUID)
			require.True(tt, strings.EqualFold(v, u.version))
			require.True(tt, strings.EqualFold(us, u.uuid))
		})
	}
	for _, u := range invalidEntriesValidUUID {
		t.Run(fmt.Sprintf("Get player info invalid entry valid uuid %s %s %s", u.uuid, u.name, u.version), func(tt *testing.T) {
			player := getPlayerInfo(tt, http.StatusOK, u.uuid)
			require.NotNil(tt, player)
			us := globalUtils.PGUUIDToString(player.UUID)
			require.True(tt, strings.EqualFold(us, u.uuid))
		})
	}
	for _, u := range invalidEntriesInvalidUUID {
		t.Run(fmt.Sprintf("Get player info invalid entry invalid uuid %s %s %s", u.uuid, u.name, u.version), func(tt *testing.T) {
			player := getPlayerInfo(tt, http.StatusBadRequest, u.uuid)
			require.Nil(tt, player)
		})
	}
	for _, ent := range []struct {
		field string
		value string
	}{
		{"", ""},
		{"uuid", ""},
		{"uuidSHIT", "as9di0!$\x20"},
		{"uuidSHIT", "d48ac2c8-0d91-4e90-9bf4-dc15e98c4b90"},
		{"fuck", "123"},
		{"123", "456"},
		{"123", ""},
		{"#F10-d!!20g?$g", ""},
		{"#F10-d!!20g?$g", "s0d1#R?WF?`ff\\t4"},
	} {
		t.Run(fmt.Sprintf("Get player info invalid request"), func(tt *testing.T) {
			getPlayerInfoInvalid(tt, http.StatusBadRequest, ent.field, ent.value)
		})
	}
}

func TestUpdatePlayerInfo(t *testing.T) {
	for _, u := range validUpdateEntries {
		t.Run(fmt.Sprintf("Update player info valid entry %s %s %s", u.uuid, u.name, u.version), func(tt *testing.T) {
			updatePlayerInfo(tt, http.StatusNoContent, u.uuid, u.name, u.version)
			player := getPlayerInfo(tt, http.StatusOK, u.uuid)
			require.NotNil(tt, player)
			require.True(tt, strings.EqualFold(player.Name, u.name))
			v := strconv.Itoa(player.Version)
			us := globalUtils.PGUUIDToString(player.UUID)
			require.True(tt, strings.EqualFold(v, u.version))
			require.True(tt, strings.EqualFold(us, u.uuid))
		})
	}

	for _, u := range invalidUpdateEntriesValidUUID {
		t.Run(fmt.Sprintf("Update player info invalid entry valid uuid %s %s %s", u.uuid, u.name, u.version), func(tt *testing.T) {
			updatePlayerInfo(tt, http.StatusBadRequest, u.uuid, u.name, u.version)
		})
	}

	for _, u := range invalidUpdateEntriesInvalidUUID {
		t.Run(fmt.Sprintf("Update player info invalid entry invalid uuid %s %s %s", u.uuid, u.name, u.version), func(tt *testing.T) {
			updatePlayerInfo(tt, http.StatusBadRequest, u.uuid, u.name, u.version)
		})
	}

	for _, ent := range []struct {
		uuid  string
		field string
		value string
	}{
		{"D48ac2c8-0D91-2E90-9BF4-DC15E98C4B90", "", ""},
		{"D48ac2c8-0D91-2E90-9BF4-DC15E98C4B90", "uuid", ""},
		{"D48ac2c8-0D91-2E90-9BF4-DC15E98C4B90", "uuidSHIT", "as9di0!$\x20"},
		{"D48ac2c8-0D91-2E90-9BF4-DC15E98C4B90", "uuidSHIT", "d48ac2c8-0d91-4e90-9bf4-dc15e98c4b90"},
		{"D48ac2c8-0D91-2E90-9BF4-DC15E98C4B90", "fuck", "123"},
		{"D48ac2c8-0D91-4E90-9BF4-DC15E98C4B90", "123", "456"},
		{"D48ac2c8-0D91-4E90-9BF4-DC15E98C4B90", "123", ""},
		{"D48ac2c8-0D91-4E90-9BF4-DC15E98C4B90", "#F10-d!!20g?$g", ""},
		{"D48ac2c8-0D91-4E90-9BF4-DC15E98C4B90", "#F10-d!!20g?$g", "s0d1#R?WF?`ff\\t4"},
	} {
		t.Run(fmt.Sprintf("Update player info invalid request"), func(tt *testing.T) {
			updatePlayerInfoInvalid(tt, http.StatusBadRequest, ent.uuid, "application/x-www-form-urlencoded", bytes.NewBufferString(fmt.Sprintf("%s=%s", ent.field, ent.value)))
		})
	}
}
