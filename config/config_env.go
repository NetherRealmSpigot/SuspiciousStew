package config

import (
	"fmt"
	"os"
	"stew/embeds"
	"strconv"
	"strings"
)

func key(key string) string {
	return fmt.Sprintf(strings.ToUpper(fmt.Sprintf("%sAPI_%s", embeds.Code, key)))
}

func find(key string) (string, bool) {
	return os.LookupEnv(key)
}

func read(key string, fallback string) string {
	v, found := find(key)
	if !found {
		return fallback
	}
	return v
}

func readStr(key string, fallback string) string {
	return read(key, fallback)
}

func parseInt(key string, bitsize int) (int64, error) {
	return strconv.ParseInt(key, 10, bitsize)
}

func parseUInt(key string, bitsize int) (uint64, error) {
	return strconv.ParseUint(key, 10, bitsize)
}

func readUInt16(key string, fallback uint16) uint16 {
	v, err := parseUInt(key, 16)
	if err != nil {
		return fallback
	}
	return uint16(v)
}

func readInt32(key string, fallback int32) int32 {
	v, err := parseInt(key, 32)
	if err != nil {
		return fallback
	}
	return int32(v)
}

func readBool(key string, fallback bool) bool {
	v, err := strconv.ParseBool(read(key, strconv.FormatBool(fallback)))
	if err != nil {
		return fallback
	}
	return v
}
