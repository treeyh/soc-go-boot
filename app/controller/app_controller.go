package controller

import (
	"github.com/treeyh/soc-go-boot/app/model/req"
	"github.com/treeyh/soc-go-boot/app/model/resp"
	"time"
)

type AppController struct {
}

func (ac *AppController) PreUrl() string {
	return "/app"
}

// @router /get/:appId [get]
func (ac *AppController) Get(ctx *req.GinContext, userId int64) *resp.HttpRespResult {
	return nil
}

// @router /create [post]
func (ac *AppController) Create(ctx *req.GinContext, updateTime, createTime time.Time, userReq *req.UserReq) *resp.HttpRespResult {
	return &resp.HttpRespResult{
		HttpStatus: 200,
		RespResult: resp.RespResult{
			Code:    0,
			Message: "",
			Data:    nil,
		},
	}
}
