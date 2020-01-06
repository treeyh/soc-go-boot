package resp

import (
	"github.com/treeyh/soc-go-boot/app/model/req"
)

type RespResult struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type HttpRespResult struct {
	RespResult

	HttpStatus int `json:"httpStatus"`
}

func Json(g *req.GinContext, code int, msg string, data interface{}) {
	g.Ctx.JSON(200, RespResult{
		Code:    code,
		Message: msg,
		Data:    data,
	})
}

func JsonRespResult(g *req.GinContext, resp *RespResult) {
	g.Ctx.JSON(200, resp)
}

func OkHttpRespResult(resp *RespResult) *HttpRespResult {
	return &HttpRespResult{
		RespResult: *resp,
		HttpStatus: 200,
	}
}

func FailHttpRespResult(httpStatus int, resp *RespResult) *HttpRespResult {
	return &HttpRespResult{
		RespResult: *resp,
		HttpStatus: httpStatus,
	}
}
