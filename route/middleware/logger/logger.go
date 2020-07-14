package logger

import (
	"bytes"
	"context"
	"fmt"
	"github.com/treeyh/soc-go-boot/common/consts"
	coreconsts "github.com/treeyh/soc-go-common/core/consts"
	"github.com/treeyh/soc-go-common/core/logger"
	"github.com/treeyh/soc-go-common/core/model"
	"github.com/treeyh/soc-go-common/core/utils/network"
	"github.com/treeyh/soc-go-common/core/utils/strs"
	"github.com/treeyh/soc-go-common/core/utils/times"
	"github.com/treeyh/soc-go-common/core/utils/uuid"
	"go.uber.org/zap"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func newTraceId() string {
	return fmt.Sprintf("%s_%s", network.GetIntranetIp(), uuid.NewUuid())
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func isBinaryContent(contentType string) bool {
	return strings.Contains(contentType, "image") || strings.Contains(contentType, "video") ||
		strings.Contains(contentType, "audio")
}

func isMultipart(contentType string) bool {
	return strings.Contains(contentType, "multipart/form-data")
}

func StartTrace() gin.HandlerFunc {
	return func(c *gin.Context) {

		traceId := c.Request.Header.Get(coreconsts.TraceIdKey)
		if "" == traceId {
			traceId = newTraceId()
		}

		app := c.Request.Header.Get(consts.HeaderApp)
		authToken := c.Request.Header.Get(consts.HeaderAuthTokenKey)
		platform := c.Request.Header.Get(consts.HeaderPlatform)
		clientVersion := c.Request.Header.Get(consts.HeaderClientVersion)

		contentType := c.ContentType()
		body := ""
		httpContext := &model.HttpContext{
			Request:       c.Request,
			Url:           c.Request.RequestURI,
			Method:        c.Request.Method,
			StartTime:     times.GetNowMillisecond(),
			EndTime:       0,
			Ip:            c.ClientIP(),
			App:           app,
			TraceId:       traceId,
			AuthToken:     authToken,
			Platform:      platform,
			ClientVersion: clientVersion,
		}

		if !isBinaryContent(contentType) && !isMultipart(contentType) {
			// 判断不是上传文件等大消息体，记录消息体日志
			//c.Request.ParseForm()
			//postForm = c.Request.PostForm.Encode()
			data, err := c.GetRawData()
			if err != nil {
				logger.Logger().Info(err.Error())
			}
			if data != nil {
				body = string(data)
				// 重新写入body
				c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(data))
			}
		}

		ctx := context.WithValue(c.Request.Context(), coreconsts.TraceIdKey, traceId)
		ctx = context.WithValue(ctx, consts.TracerHttpContextKey, httpContext)
		c.Request = c.Request.WithContext(ctx)

		c.Header(coreconsts.TraceIdKey, traceId)
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		c.Next()

		urlCount := strs.LengthUnicode(httpContext.Url)
		if urlCount <= 6 || httpContext.Url[urlCount-6:] != "health" {
			// 仅记录非心跳日志

			httpContext = c.Request.Context().Value(consts.TracerHttpContextKey).(*model.HttpContext)
			httpContext.Status = c.Writer.Status()
			httpContext.EndTime = times.GetNowMillisecond()
			runtime := httpContext.EndTime - httpContext.StartTime
			runtimes := strconv.FormatInt(runtime, 10)
			httpStatus := strconv.Itoa(httpContext.Status)
			//msg := fmt.Sprintf("request|traceId=%s|clientVersion=%s|authToken=%s|platform=%s|start=%s|ip=%s|contentType=%s|method=%s|url=%s|body=%s|------response|end=%s|time=%s|status=%s|body=%s|",
			//	httpContext.TraceId, clientVersion, authToken, platform, times.GetDateTimeStrByMillisecond(httpContext.StartTime), httpContext.Ip, contentType,
			//	httpContext.Method, httpContext.Url, strings.ReplaceAll(body, "\n", "\\n"), times.GetDateTimeStrByMillisecond(httpContext.EndTime),
			//	runtimes, httpStatus, blw.body.String())

			logger.LoggerByName("rr").InfoCtx(c.Request.Context(), blw.body.String(), zap.String("clientVersion", clientVersion), zap.String("authToken", authToken),
				zap.String("duration", runtimes), zap.String("app", app), zap.String("platform", platform), zap.String("requestBody", strings.ReplaceAll(body, "\n", "\\n")),
				zap.String("start", times.GetDateTimeStrByMillisecond(httpContext.StartTime)), zap.String("end", times.GetDateTimeStrByMillisecond(httpContext.EndTime)),
				zap.String("ip", httpContext.Ip), zap.String("contentType", contentType),
				zap.String("method", httpContext.Method), zap.String("url", httpContext.Url), zap.String("httpStatus", httpStatus))

		}
	}
}
