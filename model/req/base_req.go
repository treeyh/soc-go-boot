package req

import "github.com/gin-gonic/gin"

type BaseReq struct {
}

type GinContext struct {
	Ctx *gin.Context
}

type PageReq struct {
	Page int `json:"page" validate:"required,min=1"`
	Size int `json:"size" validate:"required,min=1,max=200"`
}
