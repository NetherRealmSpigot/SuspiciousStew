package utils

import (
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"stew/embeds"
	"stew/logging"
)

func PGUUIDToString(uuid pgtype.UUID) string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid.Bytes[0:4], uuid.Bytes[4:6], uuid.Bytes[6:8], uuid.Bytes[8:10], uuid.Bytes[10:16])
}

func Info() {
	logging.AppLogger.Info(fmt.Sprintf("%s (ver %s)", embeds.Code, embeds.ExecutableVersion))
}
