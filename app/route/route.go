package route

import "github.com/gin-gonic/gin"

func SetupRouter(engine *gin.Engine) {

	engine.Use(jaeger.SetUp())

	//404
	engine.NoRoute(func(c *gin.Context) {
		utilGin := util.Gin{Ctx: c}
		utilGin.Response(404, "请求方法不存在", nil)
	})

	engine.GET("/sing", func(c *gin.Context) {
		utilGin := util.Gin{Ctx: c}
		utilGin.Response(1, "sing", nil)
	})
}
