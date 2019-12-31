package http

import (
	"crypto/tls"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/treeyh/soc-go-boot/app/common/consts"
	"github.com/treeyh/soc-go-boot/app/config"
	"github.com/treeyh/soc-go-common/core/errors"
	"github.com/treeyh/soc-go-common/core/logger"
	"io/ioutil"
	"net/http"
	neturl "net/url"
	"strings"
	"time"
)

func Get(c *gin.Context, url string, querys map[string]string) (string, int, errors.AppError) {
	return do(c, "GET", url, querys, nil, nil)
}

func Post(c *gin.Context, url string, querys map[string]string, body *string) (string, int, errors.AppError) {
	return do(c, "POST", url, querys, nil, body)
}

func Put(c *gin.Context, url string, querys map[string]string, body *string) (string, int, errors.AppError) {
	return do(c, "PUT", url, querys, nil, body)
}

func Delete(c *gin.Context, url string, querys map[string]string, body *string) (string, int, errors.AppError) {
	return do(c, "DELETE", url, querys, nil, body)
}

func do(c *gin.Context, method string, url string, querys map[string]string, headers map[string]string, body *string) (string, int, errors.AppError) {

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Timeout:   time.Second * 3, //默认3秒超时时间
		Transport: tr,
	}

	// 拼接url
	reqUrl := url
	if querys != nil {
		reqUrl += "?" + ConvertToQueryParams(querys)
	}

	// 构建request
	var reqBody *strings.Reader
	if body != nil {
		reqBody = strings.NewReader(*body)
	}
	req, err := http.NewRequest(method, reqUrl, reqBody)
	if err != nil {
		return "", 0, errors.NewAppErrorByExistError(errors.HttpCreateRequestFail, err)
	}

	//设置header
	if headers != nil {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}

	if config.GetSocConfig().Trace.Enable {

		tracer, _ := c.Get(consts.TracerContextKey)
		parentSpanContext, _ := c.Get(consts.TraceParentSpanContextKey)

		span := opentracing.StartSpan(
			"call Http "+method,
			opentracing.ChildOf(parentSpanContext.(opentracing.SpanContext)),
			opentracing.Tag{Key: string(ext.Component), Value: "HTTP"},
			ext.SpanKindRPCClient,
		)

		span.Finish()

		injectErr := tracer.(opentracing.Tracer).Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
		if injectErr != nil {
			logger.Logger().Fatalf("%s: Couldn't inject headers", err)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", 0, errors.NewAppErrorByExistError(errors.HttpRequestFail, err)
	}

	b, err := ioutil.ReadAll(resp.Body)
	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()

	if err != nil {
		logger.Logger().Error(err)
		return "", resp.StatusCode, errors.NewAppErrorByExistError(errors.HttpRequestFail, err)
	}
	return string(b), resp.StatusCode, nil
}

func ConvertToQueryParams(queryParams map[string]string) string {
	rq := neturl.Values{}
	for k, v := range queryParams {
		rq.Add(k, v)
	}
	return rq.Encode()
}
