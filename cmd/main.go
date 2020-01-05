package main

import (
	"context"
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/treeyh/soc-go-boot/app/common/buildinfo"
	"github.com/treeyh/soc-go-boot/app/config"
	"github.com/treeyh/soc-go-boot/app/route"
	"github.com/treeyh/soc-go-common/core/consts"
	"github.com/treeyh/soc-go-common/core/logger"
	"github.com/treeyh/soc-go-common/core/utils/json"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

const configRelativePath = "./config"

var (
	env = ""
	log = logger.Logger()
)

func init() {
	// 获取环境变量
	env = consts.GetCurrentEnv()

	flag.StringVar(&env, "env", env, "env")
	flag.Parse()

	//configPath := file.GetCurrentPath() + configRelativePath
	config.LoadEnvConfig(configRelativePath, "application", env)

	// 重新初始化日志配置
	for k, v := range *config.GetSocConfig().Logger {
		logger.InitLogger(k, &v, true)
	}

	log.Info("app env: " + env)
	jstr, err := json.ToJson(config.GetSocConfig())
	if err != nil {
		panic(" init config json fail....")
	}
	log.Info("app config: " + jstr)

}

// gracefulShutdown 优雅关机
func gracefulShutdown(srv *http.Server) {
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown: " + err.Error())
	}

	log.Info("Server exiting")
}

func main() {
	//打印程序信息

	logger.Logger().Info(buildinfo.StringifySingleLine())

	port := strconv.Itoa(config.GetSocConfig().App.Server.Port)

	engine := gin.New()

	route.SetupRouter(engine)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: engine,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s", err)
		}
	}()

	gracefulShutdown(srv)

	//os.Remove("/Users/tree/work/99_tree/03_github/soc-go-boot/app/route/abc.go")
	//route.DemoPrint()
}
