package jaeger

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/treeyh/soc-go-boot/app/common/consts"
	"github.com/treeyh/soc-go-boot/app/common/utils/jaeger_trace"
	"github.com/treeyh/soc-go-boot/app/config"
)

func SetUp() gin.HandlerFunc {

	return func(c *gin.Context) {
		if config.GetSocConfig().Trace.Enable {

			var parentSpan opentracing.Span

			tracer, closer := jaeger_trace.NewJaegerTracer(config.GetSocConfig().App.Name, config.GetSocConfig().Trace.Server)
			defer closer.Close()

			spCtx, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
			if err != nil {
				parentSpan = tracer.StartSpan(c.Request.URL.Path)
				defer parentSpan.Finish()
			} else {
				parentSpan = opentracing.StartSpan(
					c.Request.URL.Path,
					opentracing.ChildOf(spCtx),
					opentracing.Tag{Key: string(ext.Component), Value: "HTTP"},
					ext.SpanKindRPCServer,
				)
				defer parentSpan.Finish()
			}
			c.Set(consts.TracerContextKey, tracer)
			c.Set(consts.TraceParentSpanContextKey, parentSpan.Context())
		}
		c.Next()
	}
}
