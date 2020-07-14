package limiter

import (
	"github.com/gin-gonic/gin"
	"github.com/treeyh/soc-go-boot/controller"
	"github.com/treeyh/soc-go-boot/model/req"
	"github.com/treeyh/soc-go-common/core/errors"
	"github.com/treeyh/soc-go-common/core/logger"
	"golang.org/x/time/rate"
	"time"
)

func SetUp(maxBurstSize int) gin.HandlerFunc {

	limiter := rate.NewLimiter(rate.Every(time.Second*1), maxBurstSize)
	return func(c *gin.Context) {
		if limiter.Allow() {
			c.Next()
			return
		}
		logger.Logger().Error("Too many requests . ")
		utilGin := req.GinContext{Ctx: c}
		controller.Json(&utilGin, 200, errors.LimitExceed.Code(), errors.LimitExceed.Message(), nil)
		c.Abort()
		return
	}
}
