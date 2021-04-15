package boot_config

import (
	"github.com/magiconair/properties/assert"
	"github.com/treeyh/soc-go-boot/tests"
	socconfig "github.com/treeyh/soc-go-common/core/config"
	socconsts "github.com/treeyh/soc-go-common/core/consts"
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
	LoadConfigForceEnv(tests.ConfigPath, socconsts.EnvUnitTest)
	assert.Equal(t, env, socconsts.EnvUnitTest)

	socTestConfig := &SocTestConfig{}
	LoadAppConfig(socTestConfig)

	assert.Equal(t, socTestConfig.AppParams.WeChat.Host, "https://api.weixin.qq.com")
}
