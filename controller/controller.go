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

// OkJson 输出成功Json结果，仅支持0或1个data
func OkJson(c *req.GinContext, data ...interface{}) {
	Json(c, 200, errors.OK.Code(), errors.OK.Message(), data...)
}

// FailJson 输出失败Json结果，仅支持0或1个data
func FailJson(c *req.GinContext, err errors.AppError, data ...interface{}) {
	Json(c, 200, err.Code(), err.Message(), data...)
}

// FailStatusJson 输出失败Json结果，仅支持0或1个data
func FailStatusJson(c *req.GinContext, httpStatus int, err errors.AppError, data ...interface{}) {
	Json(c, httpStatus, err.Code(), err.Message(), data...)
}

// RespJson 输出Json结果，仅支持0或1个data
func Json(c *req.GinContext, httpStatus int, code int, msg string, data ...interface{}) {
	if len(data) > 0 {
		c.Ctx.JSON(httpStatus, resp.RespResult{
			Code:      code,
			Message:   msg,
			Data:      data[0],
			Timestamp: time.Now().Unix(),
		})
	} else {
		c.Ctx.JSON(httpStatus, resp.RespResult{
			Code:      code,
			Message:   msg,
			Timestamp: time.Now().Unix(),
		})
	}
}

func TextHttpRespResult(g *req.GinContext, resp *resp.HttpTextRespResult) {
	if resp.HttpStatus == 0 {
		resp.HttpStatus = 200
	}
	g.Ctx.String(resp.HttpStatus, "%s", resp.Text)
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
