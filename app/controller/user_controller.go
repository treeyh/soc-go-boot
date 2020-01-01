package controller

import "github.com/gin-gonic/gin"

type UserController struct {
	BaseController
}

func (UserController) GetUser(ctx gin.Context, userId int64) {

}

func (UserController) CreateUser(ctx gin.Context) {

}
