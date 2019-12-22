package main

import (
	"flag"
	"fmt"
	"github.com/treeyh/soc-go-boot/app/config"
	"github.com/treeyh/soc-go-common/core/consts"
	"github.com/treeyh/soc-go-common/core/logger"
	"github.com/treeyh/soc-go-common/core/utils/json"
	"os"
	"path"
	"runtime"
)

const configRelativePath = "/../config"

var (
	env = ""
	log = logger.Logger()
)

func init() {
	env = os.Getenv(consts.EvnRunName)
	if "" == env {
		env = consts.EnvLocal
	}

	flag.StringVar(&env, "env", env, "env")
	flag.Parse()

	configPath := GetCurrentPath() + configRelativePath
	config.LoadEnvConfig(configPath, "application", env)

	for k, v := range *config.GetSocConfig().Logger {
		logger.InitLogger(k, &v, true)
	}

}

func GetCurrentPath() string {
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath := path.Dir(filename)
		return abPath
	}
	return ""
}

func main() {
	//打印程序信息
	log.Info(GetCurrentPath())
	js, _ := json.ToJson(config.GetSocConfig())
	log.Info(js)

	fmt.Println(log)
}
