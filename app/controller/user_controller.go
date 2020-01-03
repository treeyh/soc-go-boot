package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/treeyh/soc-go-boot/app/model/req"
	"github.com/treeyh/soc-go-boot/app/model/resp"
	"time"
)

type UserController struct {
}

//func (uc *UserController) Version() string {
//	return "v1"
//}

func (uc *UserController) Get(ctx *gin.Context, userId int64) *resp.RespResult {
	return nil
}

func (uc *UserController) Create(ctx *gin.Context, createTime time.Time, userReq *req.UserReq) *resp.RespResult {
	return &resp.RespResult{
		Code:    0,
		Message: "OK",
		Data:    nil,
	}
}
