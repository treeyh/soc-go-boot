package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/treeyh/soc-go-boot/boot_config"
	"github.com/treeyh/soc-go-boot/model"
	"github.com/treeyh/soc-go-boot/model/resp"
	"github.com/treeyh/soc-go-common/core/errors"
	"github.com/treeyh/soc-go-common/core/logger"
	"github.com/treeyh/soc-go-common/library/i18n"
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

func BuildRespResult(appError errors.AppError, data ...interface{}) *resp.RespResult {

	if len(data) > 0 {
		return &resp.RespResult{
			Code:      appError.Code(),
			Message:   appError.Message(),
			Data:      data[0],
			Timestamp: time.Now().Unix(),
		}
	}
	return &resp.RespResult{
		Code:      appError.Code(),
		Message:   appError.Message(),
		Timestamp: time.Now().Unix(),
	}

}

func BuildOkRespResult(data interface{}) *resp.RespResult {
	return &resp.RespResult{
		Code:      errors.OK.Code(),
		Message:   errors.OK.Message(),
		Data:      data,
		Timestamp: time.Now().Unix(),
	}
}

func BuildFailRespResult(appError errors.ResultCode, data ...interface{}) *resp.RespResult {
	if len(data) > 0 {
		return &resp.RespResult{
			Code:      appError.Code(),
			Message:   appError.Message(),
			Data:      data[0],
			Timestamp: time.Now().Unix(),
		}
	}
	return &resp.RespResult{
		Code:      appError.Code(),
		Message:   appError.Message(),
		Timestamp: time.Now().Unix(),
	}
}

func OkHttpRespResult(data ...interface{}) *resp.HttpJsonRespResult {
	var result *resp.HttpJsonRespResult
	if len(data) > 0 {
		result = &resp.HttpJsonRespResult{
			Resp: resp.RespResult{
				Code:      errors.OK.Code(),
				Message:   errors.OK.Message(),
				Timestamp: time.Now().Unix(),
				Data:      data[0],
			},
			HttpStatus: 200,
		}
	} else {
		result = &resp.HttpJsonRespResult{
			Resp: resp.RespResult{
				Code:      errors.OK.Code(),
				Message:   errors.OK.Message(),
				Timestamp: time.Now().Unix(),
			},
			HttpStatus: 200,
		}
	}

	return result
}

func HttpRespResult(respResult *resp.RespResult) *resp.HttpJsonRespResult {
	return &resp.HttpJsonRespResult{
		Resp:       *respResult,
		HttpStatus: 200,
	}
}

// OkJson 输出成功Json结果，仅支持0或1个data
func OkJson(c *gin.Context, data ...interface{}) {
	Json(c, 200, errors.OK.Code(), errors.OK.Message(), data...)
}

// FailJson 输出失败Json结果，仅支持0或1个data
func FailJson(c *gin.Context, err errors.AppError, data ...interface{}) {
	Json(c, 200, err.Code(), err.Message(), data...)
}

// FailStatusJson 输出失败Json结果，仅支持0或1个data
func FailStatusJson(c *gin.Context, httpStatus int, err errors.AppError, data ...interface{}) {
	Json(c, httpStatus, err.Code(), err.Message(), data...)
}

// Json 输出Json结果，仅支持0或1个data
func Json(c *gin.Context, httpStatus int, code int, msg string, data ...interface{}) {

	message := msg
	if boot_config.GetSocConfig().I18n.Enable {
		message = i18n.GetByDefault(model.GetHttpContext(c.Request.Context()).Lang, fmt.Sprintf("ErrorMsg.%d", code), msg)
	}
	if len(data) > 0 {
		c.JSON(httpStatus, resp.RespResult{
			Code:      code,
			Message:   message,
			Data:      data[0],
			Timestamp: time.Now().Unix(),
		})
	} else {
		c.JSON(httpStatus, resp.RespResult{
			Code:      code,
			Message:   message,
			Timestamp: time.Now().Unix(),
		})
	}
}

func TextHttpRespResult(c *gin.Context, resp *resp.HttpTextRespResult) {
	if resp.HttpStatus == 0 {
		resp.HttpStatus = 200
	}
	c.String(resp.HttpStatus, "%s", resp.Text)
}

func HtmlHttpRespResult(c *gin.Context, resp *resp.HttpHtmlRespResult) {
	if resp.HttpStatus == 0 {
		resp.HttpStatus = 200
	}
	c.HTML(resp.HttpStatus, resp.Name, resp.Data)
}

func XmlHttpRespResult(c *gin.Context, resp *resp.HttpXmlRespResult) {
	if resp.HttpStatus == 0 {
		resp.HttpStatus = 200
	}
	c.XML(resp.HttpStatus, resp.Data)
}

func ProtoBufHttpRespResult(c *gin.Context, resp *resp.HttpProtoBufRespResult) {
	if resp.HttpStatus == 0 {
		resp.HttpStatus = 200
	}
	c.ProtoBuf(resp.HttpStatus, resp.Data)
}

func RedirectHttpRespResult(c *gin.Context, resp *resp.HttpRedirectRespResult) {
	if resp.HttpStatus == 0 {
		resp.HttpStatus = 302
	}
	c.Redirect(resp.HttpStatus, resp.Location)
}

func FileHttpRespResult(c *gin.Context, resp *resp.HttpFileRespResult) {
	if resp.HttpStatus == 0 {
		resp.HttpStatus = 200
	}
	if resp.FileName == "" {
		resp.FileName = filepath.Base(resp.FilePath)
	}
	c.FileAttachment(resp.FilePath, resp.FileName)
}
