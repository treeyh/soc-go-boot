package controller

import (
	"github.com/gin-gonic/gin"
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
func (ac *AppController) Get(ctx *gin.Context, userId int64) *resp.RespResult {
	return nil
}

// @router /create [post]
func (ac *AppController) Create(ctx *gin.Context, updateTime, createTime time.Time, userReq *req.UserReq) *resp.RespResult {
	return &resp.RespResult{
		Code:    0,
		Message: "OK",
		Data:    nil,
	}
}
