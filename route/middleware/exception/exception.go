package exception

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/treeyh/soc-go-boot/model/resp"
	"github.com/treeyh/soc-go-common/core/errors"
	"github.com/treeyh/soc-go-common/core/logger"
	"runtime/debug"
	"strings"
	"time"
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

				c.JSON(200, resp.RespResult{
					Code:      errors.SystemErr.Code(),
					Message:   errors.SystemErr.Message(),
					Timestamp: time.Now().Unix(),
				})
				c.Abort()
				return
			}
		}()
		c.Next()
	}
}
