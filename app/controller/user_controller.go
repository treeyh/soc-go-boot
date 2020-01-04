package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/treeyh/soc-go-boot/app/model/req"
	"github.com/treeyh/soc-go-boot/app/model/resp"
	"time"
)

type UserController struct {
}

// 无参数
// PreUrl Url前缀
func (uc *UserController) PreUrl() string {
	return "/user"
}

// @router /get/:userId [get,post]
func (uc *UserController) Get(ctx *gin.Context, userId int64) *resp.RespResult {
	return nil
}

// @params updateTime
// @router /create [*]
func (uc *UserController) Create(ctx *gin.Context, updateTime, createTime time.Time, userReq *req.UserReq) *resp.RespResult {
	return &resp.RespResult{
		Code:    0,
		Message: "OK",
		Data:    nil,
	}
}

// @router /add [post]
func Create(ctx *gin.Context, updateTime, createTime time.Time, userReq *req.UserReq) *resp.RespResult {
	return &resp.RespResult{
		Code:    0,
		Message: "OK",
		Data:    nil,
	}
}
