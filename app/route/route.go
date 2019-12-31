package route

import (
	"github.com/gin-gonic/gin"
	"github.com/treeyh/soc-go-boot/app/model/response"
	"github.com/treeyh/soc-go-boot/app/route/middleware/exception"
	"github.com/treeyh/soc-go-boot/app/route/middleware/jaeger"
)

func SetupRouter(engine *gin.Engine) {

	engine.Use(exception.SetUp())
	engine.Use(jaeger.SetUp())

	//404
	engine.NoRoute(func(c *gin.Context) {
		utilGin := response.GinContext{c}
		utilGin.Json(404, "请求方法不存在", nil)
	})

	engine.GET("/sing", func(c *gin.Context) {
		utilGin := response.GinContext{c}
		utilGin.Json(200, "ok", nil)
	})

}
