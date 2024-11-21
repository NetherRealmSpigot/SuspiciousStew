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
	"testing"
)

func addIpInfo(t *testing.T, expectStatus int, ipString string) {
	resp, err := http.PostForm(fmt.Sprintf("http://%s:%d%s",
		router.ListenAddr, router.ListenPort, gateway.RouteGroup+gateway.IpInfoPath),
		url.Values{
			"ip": []string{ipString},
		})
	require.NoError(t, err)
	require.Equal(t, expectStatus, resp.StatusCode)
	defer resp.Body.Close()
}

func addIpInfoInvalid(t *testing.T, expectStatus int, contentType string, reader io.Reader) {
	resp, err := http.Post(fmt.Sprintf("http://%s:%d%s",
		router.ListenAddr, router.ListenPort, gateway.RouteGroup+gateway.IpInfoPath), contentType, reader)
	require.NoError(t, err)
	require.Equal(t, expectStatus, resp.StatusCode)
	defer resp.Body.Close()
}

var validIP = []string{
	"127.0.0.1", "127.0.0.2", "172.16.0.1", "172.33.14.65", "192.168.0.1",
	"192.168.5.2", "10.0.0.1", "10.69.1.2", "1.0.0.4", "1.2.3.4",
	"56.78.90.12", "192.168.255.254", "223.255.255.254",
	"192.168.255.0", "192.168.255.255",
}

var invalidIP = []string{
	"0.0.0.0", "224.0.0.0", "224.0.0.1", "238.255.255.254", "238.255.255.255",
	"239.0.0.0", "239.0.0.1", "255.255.255.254", "255.255.255.255",
	"192.168.255.256",
	"99999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999",
	"AnG!#v0s9dAnG!#v0s9dAnG!#v0s9d\x20\x20\x20\x20\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00",
	" WHERE 1=1 --", ";;", "",
}

func TestAddIpInfo(t *testing.T) {
	for _, ip := range validIP {
		t.Run(fmt.Sprintf("Add ip info valid entry %s", ip), func(tt *testing.T) {
			addIpInfo(tt, http.StatusNoContent, ip)
		})
	}

	for _, ip := range invalidIP {
		t.Run(fmt.Sprintf("Add ip info invalid entry %s", ip), func(tt *testing.T) {
			addIpInfo(tt, http.StatusBadRequest, ip)
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
		t.Run("Add ip info invalid request", func(tt *testing.T) {
			addIpInfoInvalid(tt, http.StatusBadRequest, ent.contentType, ent.reader)
		})
	}
}

func getIpInfo(t *testing.T, expectStatus int, ipString string) []types.IpInfoResponse {
	resp, err := http.Get(fmt.Sprintf("http://%s:%d%s?ip=%s",
		router.ListenAddr, router.ListenPort, gateway.RouteGroup+gateway.IpInfoPath, ipString))
	require.NoError(t, err)
	require.Equal(t, expectStatus, resp.StatusCode)
	defer resp.Body.Close()
	var ips []types.IpInfoResponse
	if expectStatus == http.StatusOK {
		err = json.NewDecoder(resp.Body).Decode(&ips)
		require.NoError(t, err)
	}
	return ips
}

func getIpInfoInvalid(t *testing.T, expectStatus int, field string, value string) {
	resp, err := http.Get(fmt.Sprintf("http://%s:%d%s?%s=%s",
		router.ListenAddr, router.ListenPort, gateway.RouteGroup+gateway.IpInfoPath, field, value))
	require.NoError(t, err)
	require.Equal(t, expectStatus, resp.StatusCode)
	defer resp.Body.Close()
}

func TestGetIpInfo(t *testing.T) {
	for _, ip := range validIP {
		t.Run(fmt.Sprintf("Get ip info valid entry %s", ip), func(tt *testing.T) {
			ips := getIpInfo(tt, http.StatusOK, ip)
			require.GreaterOrEqual(tt, len(ips), 1)
		})
	}
	for _, ip := range invalidIP {
		t.Run(fmt.Sprintf("Get ip info invalid entry %s", ip), func(tt *testing.T) {
			ips := getIpInfo(tt, http.StatusBadRequest, ip)
			require.LessOrEqual(tt, len(ips), 0)
		})
	}

	for _, ent := range []struct {
		field string
		value string
	}{
		{"ip", ""},
		{"ipSHIT", "1.2.3.4"},
		{"", "1.2.3.4"},
		{"", ""},
		{"fuck", "123"},
		{"123", "456"},
		{"123", ""},
		{"#F10-d!!20g?$g", ""},
		{"#F10-d!!20g?$g", "s0d1#R?WF?`ff\\t4"},
	} {
		t.Run(fmt.Sprintf("Get ip info invalid request"), func(tt *testing.T) {
			getIpInfoInvalid(tt, http.StatusBadRequest, ent.field, ent.value)
		})
	}
}
