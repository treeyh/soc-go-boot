package config

import (
	"fmt"
	"github.com/smartystreets/goconvey/convey"
	"github.com/treeyh/soc-go-boot/tests"
	socconfig "github.com/treeyh/soc-go-common/core/config"
	socconsts "github.com/treeyh/soc-go-common/core/consts"
	"github.com/treeyh/soc-go-common/core/utils/json"
	"testing"
)

type SocTestConfig struct {
	AppParams *AppParamsConfig
}

type AppParamsConfig struct {
	Test   string
	Email  string
	WeChat *socconfig.WeChatConfig
}

func TestLoadConfigForceEnv(t *testing.T) {
	convey.Convey("config TestLoadConfigForceEnv", t, tests.TestStartUp(func() {
		convey.So(env, convey.ShouldEqual, socconsts.EnvUnitTest)

		socTestConfig := &SocTestConfig{}
		LoadCustomConfig(socTestConfig)
		convey.So(socTestConfig.AppParams.WeChat.Host, convey.ShouldEqual, "https://api.weixin.qq.com")
		fmt.Println(json.ToJsonIgnoreError(socTestConfig))

	}, func() {
		LoadConfigForceEnv(tests.ConfigPath, socconsts.EnvUnitTest)
	}))
}
