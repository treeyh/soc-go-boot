package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/treeyh/soc-go-boot/app/model/req"
	"github.com/treeyh/soc-go-boot/app/model/resp"
	"github.com/treeyh/soc-go-boot/app/model/validate"
	"github.com/treeyh/soc-go-common/core/utils/json"
	"time"
)

// @PreUrl /user 暂不支持
type UserController struct {
}

// 无参数
// PreUrl Url前缀
func (uc *UserController) PreUrl() string {
	return "/user"
}

// @Param	userId		    query	 true	1	"The email for login"
// @Router /get/:userId [get,post]
func (uc *UserController) Get(ctx *req.GinContext, userId int64) *resp.HttpRespResult {
	return nil
}

// Param@   参数名（对应方法中参数名）      取值来源（formData、query、path、body、header(参数名"-"用"_"符号代替)）   是否必须(true,false)  "默认值（可以没有）"   "注释"
// @Param	updateTime		query	 true	"2012-12-12 12:12:11"	"The email for login"
// @Param	createTime		query	 true	"2012-12-12 12:12:11"	"The email for login"
// @Param	userId		    query	 true	1	"The email for login"
// @Param	userName		query	 33222	"2012-12-12 12:12:11"	"The email for login"
// @Router /create [*]
func (uc *UserController) Create(ctx *req.GinContext, updateTime, createTime time.Time, userId int64, userName string, userReq *req.UserReq) *resp.HttpRespResult {
	fmt.Println(updateTime)
	fmt.Println(ctx.Ctx.Param("updateTime"))
	fmt.Println(ctx.Ctx.GetPostForm("updateTime"))
	fmt.Println(json.ToJson(userReq))

	return resp.OkHttpRespResult(&resp.RespResult{
		Code:    0,
		Message: updateTime.String(),
		Data:    userReq,
	})
}

// @router /add [post]
func Create() gin.HandlerFunc {
	return func(context *gin.Context) {

		//context.ShouldBindBodyWith()

		fmt.Println(context.Param("create"))
		fmt.Println(context.Param("updateTime"))
		fmt.Println(context.Query("updateTime"))
		fmt.Println(context.GetPostForm("updateTime"))
		user := req.UserReq{}
		err := context.ShouldBindBodyWith(&user, binding.JSON)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(json.ToJson(user))

		msg := validate.ValidateObject(user)
		if msg != nil {
			fmt.Println(*msg)
		}

		context.JSON(200, resp.RespResult{
			Code:    0,
			Message: "OK",
			Data:    nil,
		})
	}
}
