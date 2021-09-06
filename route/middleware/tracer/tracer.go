package tracer

import (
	"github.com/SkyAPM/go2sky"
	v3 "github.com/SkyAPM/go2sky-plugins/gin/v3"
	"github.com/SkyAPM/go2sky/reporter"
	"github.com/gin-gonic/gin"
	"github.com/treeyh/soc-go-boot/boot_config"
	"github.com/treeyh/soc-go-boot/common/boot_consts/boot_error_consts"
	"github.com/treeyh/soc-go-common/core/consts"
	"github.com/treeyh/soc-go-common/core/errors"
	"github.com/treeyh/soc-go-common/core/logger"
)

var (
	log = logger.Logger()
)

func SetUp(engine *gin.Engine) gin.HandlerFunc {
	if boot_config.GetSocConfig().Trace.Enable {
		var report go2sky.Reporter
		var err error
		if boot_config.GetSocConfig().Trace.Server != "" {
			report, err = reporter.NewGRPCReporter(boot_config.GetSocConfig().Trace.Server)
		} else {
			report, err = reporter.NewLogReporter()
		}
		if err != nil {
			log.Error("SkyWalking init fail." + err.Error())
			panic(errors.NewAppErrorByExistError(boot_error_consts.SkyWalkingNotInit, err))
		}
		consts.SetReporter(report)

		var tracer *go2sky.Tracer
		//defer rp.Close()
		tracer, err = go2sky.NewTracer(boot_config.GetSocConfig().App.Name, go2sky.WithReporter(report))
		if err != nil {
			log.Error("SkyWalking init tracer fail." + err.Error())
			panic(errors.NewAppErrorByExistError(boot_error_consts.SkyWalkingNotInit, err))
		}
		consts.SetTracer(tracer)
	}

	return v3.Middleware(engine, consts.GetTracer())
}

// CloseReport 关闭report
func CloseReport() {
	if consts.GetReporter() != nil {
		consts.GetReporter().Close()
	}
}
