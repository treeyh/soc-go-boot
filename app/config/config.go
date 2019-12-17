package config

import (
	"fmt"
	"github.com/spf13/viper"
)

var conf SocConfig = SocConfig{
	Viper:         viper.New(),
	Mysql:         nil,
	Redis:         nil,
	ScheduleTime:  nil,
	Logs:          nil,
	Application:   nil,
	ElasticSearch: nil,
}

type SocConfig struct {
	Viper         *viper.Viper
	Mysql         *MysqlConfig        //数据库配置
	Redis         *RedisConfig        //redis配置
	ScheduleTime  *ScheduleTimeConfig //定时时间配置
	Logs          *map[string]LogConfig
	Application   *ApplicationConfig   //应用配置
	ElasticSearch *ElasticSearchConfig //es配置
}

func loadExtraConfig(dir string, config string, env string, extraConfig interface{}) error {
	err := loadConfig(dir, config, env)
	if err != nil {
		return err
	}
	if err := conf.Viper.Unmarshal(&extraConfig); err != nil {
		return err
	}
	return nil
}

func loadConfig(dir string, config string, env string) error {
	configName := config
	if env != "" {
		configName += "." + env
	}
	if conf.Viper == nil {
		conf.Viper = viper.New()
	}
	conf.Viper.SetConfigName(configName)
	conf.Viper.AddConfigPath(dir)
	conf.Viper.SetConfigType("yaml")
	if err := conf.Viper.MergeInConfig(); err != nil {
		fmt.Println(err)
		return err
	}
	if err := conf.Viper.Unmarshal(&conf); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
