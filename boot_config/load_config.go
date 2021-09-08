package boot_config

import (
	"flag"
	"github.com/treeyh/soc-go-common/core/consts"
	"github.com/treeyh/soc-go-common/core/logger"
	"github.com/treeyh/soc-go-common/core/utils/json"
	"github.com/treeyh/soc-go-common/library/aliyun"
	"github.com/treeyh/soc-go-common/library/database"
	"github.com/treeyh/soc-go-common/library/redis"
	"github.com/treeyh/soc-go-common/library/tracing"
	"github.com/treeyh/soc-go-common/library/wechat"
)

var (
	env = ""
	log = logger.Logger()
)

// GetEnv 获取环境
func GetEnv() string {
	return env
}

// LoadConfig 加载配置文件， configRelativePath 配置文件路径
func LoadConfig(configRelativePath string) string {
	// 获取环境变量
	env = consts.GetCurrentEnv()

	flag.StringVar(&env, "env", env, "env")
	flag.Parse()

	loadConfigInfo(configRelativePath)

	return env
}

// LoadConfigForceEnv 加载配置文件， configRelativePath 配置文件路径， forceEnv 强制环境变量
func LoadConfigForceEnv(configRelativePath string, forceEnv string) string {
	// 获取环境变量
	env = forceEnv

	loadConfigInfo(configRelativePath)

	return env
}

// LoadAppConfig 加载应用配置
func LoadAppConfig(customConfig interface{}) {
	conf.Viper.Unmarshal(customConfig)
}

// loadConfigInfo 加载配置文件， configRelativePath 配置文件路径
func loadConfigInfo(configRelativePath string) string {
	//configPath := file.GetCurrentPath() + configRelativePath
	LoadEnvConfig(configRelativePath, "application", env)
	effectConfig()

	log.Info("app env: " + env)
	jstr, err := json.ToJson(GetSocConfig())
	if err != nil {
		panic(" init config json fail....")
	}
	log.Info("app config: " + jstr)

	return env
}

// effectConfig 生效系统配置
func effectConfig() {
	// 重新初始化日志配置
	for k, v := range *GetSocConfig().Logger {
		logger.InitLogger(k, &v, true)
	}

	log.Info(json.ToJsonIgnoreError(*GetSocConfig()))

	// 初始化sky walking配置,需要优先初始化
	if GetSocConfig().Trace != nil {
		tracing.InitTracing(*GetSocConfig().Trace, GetSocConfig().App.Name)
	}

	// 初始化缓存配置
	if GetSocConfig().Redis != nil {
		redis.InitRedisPool(*GetSocConfig().Redis)
	}

	// 初始化数据库配置
	if GetSocConfig().DataSource != nil {
		database.InitDataSource(*GetSocConfig().DataSource)
	}

	// 初始化微信配置
	if GetSocConfig().WeChat != nil {
		wechat.InitWeChatConfig(*GetSocConfig().WeChat)
	}

	// 初始化阿里云配置
	if GetSocConfig().ALiYun != nil {
		aliyun.InitALiYunConfig(*GetSocConfig().ALiYun)
	}

}
