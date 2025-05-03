package logger

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/SkyAPM/go2sky"
	"github.com/treeyh/soc-go-boot/common/boot_consts"
	"github.com/treeyh/soc-go-boot/model"
	"github.com/treeyh/soc-go-common/core/consts"
	"github.com/treeyh/soc-go-common/core/logger"
	"github.com/treeyh/soc-go-common/core/utils/network"
	"github.com/treeyh/soc-go-common/core/utils/slice"
	"github.com/treeyh/soc-go-common/core/utils/strs"
	"github.com/treeyh/soc-go-common/core/utils/times"
	"github.com/treeyh/soc-go-common/core/utils/uuid"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

var _ignoreLogUrls = make([]string, 0)

func newTraceId() string {
	return fmt.Sprintf("%s_%s", network.GetIntranetIp(), uuid.NewUuid())
}

func isBinaryContent(contentType string) bool {
	return strings.Contains(contentType, "image") || strings.Contains(contentType, "video") ||
		strings.Contains(contentType, "audio")
}

func isMultipart(contentType string) bool {
	return strings.Contains(contentType, "multipart/form-data")
}

// isNeedBody 是否需要body体
func isNeedBody(contentType string) bool {
	return contentType == "" ||
		strings.Contains(contentType, "/json") ||
		strings.Contains(contentType, "javascript") ||
		strings.Contains(contentType, "/html") ||
		strings.Contains(contentType, "/plain") ||
		strings.Contains(contentType, "/xhtml")
}

// getTraceIdSpanId 获取traceid和spanid
func getTraceIdSpanId(c *gin.Context) (string, string) {
	traceId := go2sky.TraceID(c.Request.Context())
	spanId := go2sky.SpanID(c.Request.Context())
	// 判断是否已有skywalking traceId
	if "" == traceId || go2sky.EmptyTraceID == traceId {
		traceId = c.Request.Header.Get(boot_consts.HeaderTraceIdKey)
		// 判断是否已有请求 traceId
		if "" == traceId || go2sky.EmptyTraceID == traceId {
			traceId = newTraceId()
		}
	}
	return traceId, strconv.Itoa(int(spanId))
}

func StartTrace(ignoreLogUrls ...string) gin.HandlerFunc {

	_ignoreLogUrls = append(_ignoreLogUrls, ignoreLogUrls...)

	return func(c *gin.Context) {

		traceId, spanId := getTraceIdSpanId(c)

		app := c.Request.Header.Get(boot_consts.HeaderAppCodeKey)
		partnerCode := c.Request.Header.Get(boot_consts.HeaderPartnerCodeKey)
		channelCode := c.Request.Header.Get(boot_consts.HeaderChannel)
		authToken := c.Request.Header.Get(boot_consts.HeaderAuthTokenKey)
		platform := c.Request.Header.Get(boot_consts.HeaderPlatform)
		clientVersion := c.Request.Header.Get(boot_consts.HeaderClientVersion)
		deviceId := c.Request.Header.Get(boot_consts.HeaderDeviceIdKey)
		lang := c.Request.Header.Get(boot_consts.HeaderLangKey)
		if lang == "" {
			acceptLang := c.Request.Header.Get("Accept-Language")
			lang = formatRequestLang(c.Request.Context(), acceptLang)
		}

		contentType := c.ContentType()
		body := ""
		httpContext := &model.HttpContext{
			Request:       c.Request,
			Url:           c.Request.RequestURI,
			Method:        c.Request.Method,
			StartTime:     times.GetNowMillisecond(),
			EndTime:       0,
			Ip:            getClientIp(c),
			App:           app,
			TraceId:       traceId,
			SpanId:        spanId,
			AuthToken:     authToken,
			Platform:      platform,
			ClientVersion: clientVersion,
			PartnerCode:   partnerCode,
			Channel:       channelCode,
			Lang:          lang,
			DeviceId:      deviceId,
		}

		if isNeedBody(contentType) {
			// 判断不是上传文件等大消息体，记录消息体日志
			// c.Request.ParseForm()
			// postForm = c.Request.PostForm.Encode()
			data, err := c.GetRawData()
			if err != nil {
				logger.Logger().Info(err.Error())
			}
			if data != nil {
				body = string(data)
				// 重新写入body
				c.Request.Body = io.NopCloser(bytes.NewBuffer(data))
				httpContext.Body = body
			}
		}

		ctx := context.WithValue(c.Request.Context(), consts.ContextTracerKey, traceId)
		ctx = context.WithValue(ctx, boot_consts.ContextHttpContextKey, httpContext)
		c.Request = c.Request.WithContext(ctx)

		c.Header(boot_consts.HeaderTraceIdKey, traceId)
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		c.Next()

		urlCount := strs.LengthUnicode(httpContext.Url)
		if urlCount <= 6 || httpContext.Url[urlCount-6:] != "health" {
			// 仅记录非心跳日志

			if slice.ContainString(c.Request.URL.Path, _ignoreLogUrls) {
				// 匹配到url忽略日志记录
				return
			}

			httpContext = model.GetHttpContext(c.Request.Context())
			httpContext.Status = c.Writer.Status()
			httpContext.EndTime = times.GetNowMillisecond()
			runtime := httpContext.EndTime - httpContext.StartTime
			runtimes := strconv.FormatInt(runtime, 10)
			httpStatus := strconv.Itoa(httpContext.Status)
			// msg := fmt.Sprintf("request|traceId=%s|clientVersion=%s|authToken=%s|platform=%s|start=%s|ip=%s|contentType=%s|method=%s|url=%s|body=%s|------response|end=%s|time=%s|status=%s|body=%s|",
			//	httpContext.TraceId, clientVersion, authToken, platform, times.GetDateTimeStrByMillisecond(httpContext.StartTime), httpContext.Ip, contentType,
			//	httpContext.Method, httpContext.Url, strings.ReplaceAll(body, "\n", "\\n"), times.GetDateTimeStrByMillisecond(httpContext.EndTime),
			//	runtimes, httpStatus, blw.body.String())

			logger.Logger().InfoCtx(c.Request.Context(),
				strings.ReplaceAll(blw.body.String(), "\n", "\\n"),
				zap.String("clientVersion", clientVersion),
				zap.String("authToken", authToken),
				zap.String("duration", runtimes),
				zap.String("app", app),
				zap.String("platform", platform),
				zap.String("channel", channelCode),
				zap.String("partner", partnerCode),
				zap.String("requestBody", strings.ReplaceAll(body, "\n", "\\n")),
				zap.String("start", times.GetDateTimeMillisecondStrByMillisecond(httpContext.StartTime)),
				zap.String("end", times.GetDateTimeMillisecondStrByMillisecond(httpContext.EndTime)),
				zap.String("ip", httpContext.Ip),
				zap.String("contentType", contentType),
				zap.String("spanId", httpContext.SpanId),
				zap.String("method", httpContext.Method),
				zap.String("url", httpContext.Url),
				zap.String("lang", httpContext.Lang),
				zap.String("httpStatus", httpStatus),
				zap.String("deviceId", httpContext.DeviceId),
				zap.String("socLog", "rr"))

		}
	}
}
