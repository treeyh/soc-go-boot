package tracer

import (
	v3 "github.com/SkyAPM/go2sky-plugins/gin/v3"
	"github.com/gin-gonic/gin"
	"github.com/treeyh/soc-go-common/library/tracing"
)

func SetUp(engine *gin.Engine) gin.HandlerFunc {

	return v3.Middleware(engine, tracing.GetTracer())
}
