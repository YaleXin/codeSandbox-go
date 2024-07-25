package middleware

import (
	"bytes"
	"codeSandbox/responses"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"math"
	"net/http"
	"os"
	"time"
)

// 自定义 ResponseWriter
type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r *responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

type LogData struct {
	Ip              string      `json:"ip"`
	Method          string      `json:"method"`
	Path            string      `json:"path"`
	Code            int         `json:"code"`
	ApplicationCode int         `json:"applicationCode"`
	Referer         string      `json:"referer"`
	UserAgent       string      `json:"userAgent"`
	ResData         interface{} `json:"resData"`
	Latency         int64       `json:"latency"` // ms
}

// Logger is the logrus logger handler
// https://github.com/toorop/gin-logrus
func Logger(logger log.FieldLogger, notLogged ...string) gin.HandlerFunc {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknow"
	}

	var skip map[string]struct{}

	if length := len(notLogged); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, p := range notLogged {
			skip[p] = struct{}{}
		}
	}

	return func(c *gin.Context) {
		// other handler can change c.Path so:
		path := c.Request.URL.Path

		// 创建一个缓存来存储响应数据
		writer := &responseBodyWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = writer

		start := time.Now()
		c.Next()
		stop := time.Since(start)

		// 获取响应体
		var resData responses.Response
		responseBodyBytes := writer.body.String()
		err := json.Unmarshal([]byte(responseBodyBytes), &resData)
		if err != nil {

		}
		latency := int64(math.Ceil(float64(stop.Nanoseconds()) / 1000000.0))
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		clientUserAgent := c.Request.UserAgent()
		referer := c.Request.Referer()
		dataLength := c.Writer.Size()
		if dataLength < 0 {
			dataLength = 0
		}

		if _, ok := skip[path]; ok {
			return
		}
		logData := LogData{
			Ip:              clientIP,
			Method:          c.Request.Method,
			Path:            path,
			Code:            statusCode,
			ApplicationCode: resData.Code,
			Referer:         referer,
			UserAgent:       clientUserAgent,
			Latency:         latency,
		}
		// 异常时，把响应信息拿到
		if resData.Code != 200 {
			logData.ResData = resData.Data
		}
		entry := logger.WithFields(log.Fields{
			"hostname":   hostname,
			"statusCode": statusCode,
			//"latency":    latency, // time to process
			"clientIP": clientIP,
			//"method":     c.Request.Method,
			//"path":       path,
			"referer":    referer,
			"dataLength": dataLength,
			//"userAgent":  clientUserAgent,
		})

		if len(c.Errors) > 0 {
			entry.Error(c.Errors.ByType(gin.ErrorTypePrivate).String())
		} else {
			marshal, _ := json.Marshal(logData)
			if statusCode >= http.StatusInternalServerError {
				entry.Error(string(marshal))
			} else if statusCode >= http.StatusBadRequest || resData.Code != 200 {
				entry.Warn(string(marshal))
			} else if resData.Code != 200 {
			} else {
				entry.Info(string(marshal))
			}
		}
	}
}

func LoggerDebug(httpMethod, absolutePath, handlerName string, nuHandlers int) {
	entry := log.WithFields(log.Fields{
		//"httpMethod":   httpMethod,
		//"absolutePath": absolutePath,
		"handlerName": handlerName,
		"nuHandlers":  nuHandlers,
	})
	msg := fmt.Sprintf("%-6s %-25s", httpMethod, absolutePath)
	entry.Debug(msg)
}
