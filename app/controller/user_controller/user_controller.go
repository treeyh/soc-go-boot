package user_controller

import (
	"github.com/gin-gonic/gin"
	"github.com/treeyh/soc-go-boot/app/model/req"
	"github.com/treeyh/soc-go-boot/app/model/resp"
	"time"
)

func Get(ctx *gin.Context, userId int64) *resp.RespResult {
	return nil
}

func Create(ctx *gin.Context, createTime time.Time, userReq *req.UserReq) *resp.RespResult {
	return &resp.RespResult{
		Code:    0,
		Message: "OK",
		Data:    nil,
	}
}
