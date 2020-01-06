package req

import "github.com/gin-gonic/gin"

type BaseReq struct {
	Operator int64 `json:"operator" validate:"omitempty,gt=0"`
}

type GinContext struct {
	Ctx *gin.Context
}
