package exception

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/treeyh/soc-go-boot/controller"
	"github.com/treeyh/soc-go-common/core/errors"
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
				DebugStack := fmt.Sprintf("%v \t", err)
				for _, v := range strings.Split(string(debug.Stack()), "\n") {
					DebugStack += v + ";"
				}

				log.ErrorCtx(c.Request.Context(), DebugStack)

				controller.FailJson(c, errors.NewAppError(errors.SystemErr))
				c.Abort()
				return
			}
		}()
		c.Next()
	}
}
