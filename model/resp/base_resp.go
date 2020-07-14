package resp

import (
	"github.com/treeyh/soc-go-common/core/errors"
	"time"
)

type RespResult struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
	Timestamp int64       `json:"timestamp"`
}

type HttpRespResult struct {
	RespResult

	HttpStatus int `json:"httpStatus"`
}

type PageRespResult struct {
	Total int64       `json:"total"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
	Data  interface{} `json:"data"`
}

type ListRespResult struct {
	Total int         `json:"total"`
	Data  interface{} `json:"data"`
}

func BuildRespResult(appError errors.AppError, data ...interface{}) *RespResult {

	if len(data) > 0 {
		return &RespResult{
			Code:      appError.Code(),
			Message:   appError.Message(),
			Data:      data[0],
			Timestamp: time.Now().Unix(),
		}
	}
	return &RespResult{
		Code:      appError.Code(),
		Message:   appError.Message(),
		Timestamp: time.Now().Unix(),
	}

}

func BuildOkRespResult(data interface{}) *RespResult {
	return &RespResult{
		Code:      errors.OK.Code(),
		Message:   errors.OK.Message(),
		Data:      data,
		Timestamp: time.Now().Unix(),
	}
}

func BuildFailRespResult(appError errors.ResultCode, data ...interface{}) *RespResult {
	if len(data) > 0 {
		return &RespResult{
			Code:      appError.Code(),
			Message:   appError.Message(),
			Data:      data[0],
			Timestamp: time.Now().Unix(),
		}
	}
	return &RespResult{
		Code:      appError.Code(),
		Message:   appError.Message(),
		Timestamp: time.Now().Unix(),
	}
}
