package logging

import (
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

var AppLogger = logrus.New()

func LoadLogger() {
	AppLogger.Formatter = &logrus.TextFormatter{
		FullTimestamp:          true,
		TimestampFormat:        time.RFC1123Z,
		ForceColors:            true,
		DisableLevelTruncation: true,
		PadLevelText:           true,
		QuoteEmptyFields:       true,
	}
	AppLogger.Out = os.Stdout
	AppLogger.Level = logrus.InfoLevel
}
