package exception

import (
	"github.com/gin-gonic/gin"
	"github.com/treeyh/soc-go-boot/app/model/response"
	"github.com/treeyh/soc-go-common/core/errors"
	"runtime/debug"
	"strings"
)

func SetUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {

				DebugStack := ""
				for _, v := range strings.Split(string(debug.Stack()), "\n") {
					DebugStack += v + ";"
				}
				utilGin := response.GinContext{Ctx: c}
				utilGin.Json(errors.ServerError.Code(), errors.ServerError.Message()+";"+DebugStack, nil)
			}
		}()
		c.Next()
	}
}
