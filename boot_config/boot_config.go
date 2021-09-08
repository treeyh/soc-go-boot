package boot_config

import (
	"github.com/spf13/viper"
	"github.com/treeyh/soc-go-common/core/config"
)

var conf = SocBootConfig{
	Viper:   viper.New(),
	SocBoot: nil,
}

type SocBootConfig struct {
	Viper   *viper.Viper
	SocBoot *SocConfig
}

type SocConfig struct {
	App        *AppConfig
	DataSource *map[string]config.DBConfig    //数据库配置
	Redis      *map[string]config.RedisConfig //redis配置
	Logger     *map[string]config.LogConfig
	Trace      *config.TraceConfig
	WeChat     *map[string]config.WeChatConfig
	ALiYun     *map[string]config.ALiYunConfig
	Params     *map[string]string
}

// AppConfig 应用配置
type AppConfig struct {
	Name    string
	Server  *ServerConfig
	AppCode string
	AppKey  string
}

// ServerConfig 服务配置
type ServerConfig struct {
	Port        int
	ContextPath string
}

func GetSocConfig() *SocConfig {
	return conf.SocBoot
}

func LoadEnvConfig(dir string, config string, env string) {
	loadConfig(dir, config, "")

	if env != "" {
		loadConfig(dir, config, env)
	}
}

// loadConfig 加载配置
func loadConfig(dir string, config string, env string) {
	configName := config
	if env != "" {
		configName += "." + env
	}
	if conf.Viper == nil {
		conf.Viper = viper.New()
	}

	conf.Viper.SetConfigName(configName)
	conf.Viper.AddConfigPath(dir)
	conf.Viper.SetConfigType("yml")
	if err := conf.Viper.MergeInConfig(); err != nil {
		panic("Load config file fail. " + err.Error())
	}
	if err := conf.Viper.Unmarshal(&conf); err != nil {
		panic("Load config file fail. " + err.Error())
	}
}
