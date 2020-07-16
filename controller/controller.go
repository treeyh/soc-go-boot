package controller

import (
	"github.com/treeyh/soc-go-boot/model/req"
	"github.com/treeyh/soc-go-boot/model/resp"
	"github.com/treeyh/soc-go-common/core/errors"
	"github.com/treeyh/soc-go-common/core/logger"
	"path/filepath"
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

func JsonHttpRespResult(g *req.GinContext, resp *resp.HttpJsonRespResult) {
	if resp.HttpStatus == 0 {
		resp.HttpStatus = 200
	}
	g.Ctx.JSON(resp.HttpStatus, resp.RespResult)
}

func TextHttpRespResult(g *req.GinContext, resp *resp.HttpTextRespResult) {
	if resp.HttpStatus == 0 {
		resp.HttpStatus = 200
	}
	g.Ctx.String(resp.HttpStatus, resp.Format, resp.Values)
}

func HtmlHttpRespResult(g *req.GinContext, resp *resp.HttpHtmlRespResult) {
	if resp.HttpStatus == 0 {
		resp.HttpStatus = 200
	}
	g.Ctx.HTML(resp.HttpStatus, resp.Name, resp.Data)
}

func XmlHttpRespResult(g *req.GinContext, resp *resp.HttpXmlRespResult) {
	if resp.HttpStatus == 0 {
		resp.HttpStatus = 200
	}
	g.Ctx.XML(resp.HttpStatus, resp.Data)
}

func ProtoBufHttpRespResult(g *req.GinContext, resp *resp.HttpProtoBufRespResult) {
	if resp.HttpStatus == 0 {
		resp.HttpStatus = 200
	}
	g.Ctx.ProtoBuf(resp.HttpStatus, resp.Data)
}

func RedirectHttpRespResult(g *req.GinContext, resp *resp.HttpRedirectRespResult) {
	if resp.HttpStatus == 0 {
		resp.HttpStatus = 302
	}
	g.Ctx.Redirect(resp.HttpStatus, resp.Location)
}

func FileHttpRespResult(g *req.GinContext, resp *resp.HttpFileRespResult) {
	if resp.HttpStatus == 0 {
		resp.HttpStatus = 200
	}
	if resp.FileName == "" {
		resp.FileName = filepath.Base(resp.FilePath)
	}
	g.Ctx.FileAttachment(resp.FilePath, resp.FileName)
}

func OkHttpRespResultByData(data ...interface{}) *resp.HttpJsonRespResult {
	result := &resp.HttpJsonRespResult{
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

func HttpRespResult(respResult *resp.RespResult) *resp.HttpJsonRespResult {
	return &resp.HttpJsonRespResult{
		RespResult: *respResult,
		HttpStatus: 200,
	}
}
