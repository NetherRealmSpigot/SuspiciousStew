package router

import (
	"bytes"
	"fmt"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"io"
	"io/ioutil"
	"net/http"
	"stew/logging"
	"time"
)

const RequestIdHeader = "X-Request-Id"

func randoms() string {
	id, err := gonanoid.Generate("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", 24)
	if err != nil {
		panic(err)
	}
	return id
}

var addRequestIdHeader = requestid.New(requestid.WithGenerator(randoms), requestid.WithCustomHeaderStrKey(RequestIdHeader))

func routerLogger(c *gin.Context) {
	startTime := time.Now()

	bodyBuffer := &bytes.Buffer{}
	reader := io.TeeReader(c.Request.Body, bodyBuffer)
	bodyBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		logging.AppLogger.Panic(err)
		panic(err)
	}
	c.Request.Body = io.NopCloser(bodyBuffer)

	c.Next()

	latencyTime := time.Since(startTime)

	logging.AppLogger.Info(fmt.Sprintf(
		"\033[34m%s \033[37m[\033[33m%s\033[37m] \033[36m%s \033[37m- \033[35m%s \033[37m%s - \033[34m%d\033[37m - \033[31m%s\033[0m",
		latencyTime,
		c.Request.Header.Get(RequestIdHeader),
		c.ClientIP(),
		c.Request.Method,
		c.Request.URL.String(),
		c.Writer.Status(),
		string(bodyBytes),
	))
}
