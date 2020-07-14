package exception

import (
	"github.com/gin-gonic/gin"
	"github.com/treeyh/soc-go-common/core/logger"
	"runtime/debug"
	"strings"
)

var (
	log = logger.Logger()
)

func SetUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				DebugStack := ""
				for _, v := range strings.Split(string(debug.Stack()), "\n") {
					DebugStack += v + ";"
				}

				log.ErrorCtx(c.Request.Context(), DebugStack)
				//utilGin := req.GinContext{Ctx: c}
				//controller.Json(&utilGin, 503, errors.ServerError.Code(), errors.ServerError.Message()+";"+DebugStack, nil)
			}
		}()
		c.Next()
	}
}
