package controller

import (
	"github.com/treeyh/soc-go-boot/model/req"
	"github.com/treeyh/soc-go-boot/model/resp"
	"github.com/treeyh/soc-go-common/core/errors"
	"github.com/treeyh/soc-go-common/core/logger"
	"time"
)

var (
	log = logger.Logger()
)

type IController interface {
	// PreUrl url前缀
	PreUrl() string
}

func Json(g *req.GinContext, httpStatus int, code int, msg string, data interface{}) {
	g.Ctx.JSON(httpStatus, resp.RespResult{
		Code:      code,
		Message:   msg,
		Data:      data,
		Timestamp: time.Now().Unix(),
	})
}

func JsonRespResult(g *req.GinContext, resp *resp.RespResult) {
	g.Ctx.JSON(200, resp)
}

func JsonHttpRespResult(g *req.GinContext, resp *resp.HttpRespResult) {
	g.Ctx.JSON(resp.HttpStatus, resp.RespResult)
}

func StringHttpRespResult(g *req.GinContext, httpStatus int, msg string, values ...interface{}) {
	g.Ctx.String(httpStatus, msg, values)
}

func OkHttpRespResultByData(data ...interface{}) *resp.HttpRespResult {
	result := &resp.HttpRespResult{
		RespResult: resp.RespResult{
			Code:      errors.OK.Code(),
			Message:   errors.OK.Message(),
			Timestamp: time.Now().Unix(),
		},
		HttpStatus: 200,
	}
	if len(data) > 0 {
		result.Data = data[0]
	}
	return result
}

func HttpRespResult(respResult *resp.RespResult) *resp.HttpRespResult {
	return &resp.HttpRespResult{
		RespResult: *respResult,
		HttpStatus: 200,
	}
}

func HttpRespResultHttpStatue(httpStatus int, respResult *resp.RespResult) *resp.HttpRespResult {
	return &resp.HttpRespResult{
		RespResult: *respResult,
		HttpStatus: httpStatus,
	}
}
