package config

import (
	"fmt"
	"github.com/spf13/viper"
	socconfig "github.com/treeyh/soc-go-common/core/config"
	"github.com/treeyh/soc-go-common/core/errors"
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
	DataSource *map[string]socconfig.DBConfig    //数据库配置
	Redis      *map[string]socconfig.RedisConfig //redis配置
	Logger     *map[string]socconfig.LogConfig
	Trace      *TraceConfig
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

// Trace trace配置
type TraceConfig struct {
	// Enable 是否开启
	Enable bool
	// Server 服务地址
	Server string
}

func GetSocConfig() *SocConfig {
	return conf.SocBoot
}

func LoadEnvConfig(dir string, config string, env string) errors.AppError {
	err := loadConfig(dir, config, "")
	if err != nil {
		return err
	}
	if env != "" {
		err = loadConfig(dir, config, env)
		if err != nil {
			return err
		}
	}
	return nil
}

func loadConfig(dir string, config string, env string) errors.AppError {
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
		fmt.Println(err)
		return errors.NewAppErrorByExistError(errors.LoadConfigFileFail, err)
	}
	if err := conf.Viper.Unmarshal(&conf); err != nil {
		fmt.Println(err)
		return errors.NewAppErrorByExistError(errors.LoadConfigFileFail, err)
	}
	return nil
}
