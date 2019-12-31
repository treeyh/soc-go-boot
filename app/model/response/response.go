package response

import "github.com/gin-gonic/gin"

type GinContext struct {
	Ctx *gin.Context
}

type RespResult struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (g *GinContext) Json(code int, msg string, data interface{}) {
	g.Ctx.JSON(200, RespResult{
		Code:    code,
		Message: msg,
		Data:    data,
	})
}
