package config

import (
	"github.com/smartystreets/goconvey/convey"
	"github.com/treeyh/soc-go-boot/tests"
	socconsts "github.com/treeyh/soc-go-common/core/consts"
	"testing"
)

func TestLoadConfigForceEnv(t *testing.T) {
	convey.Convey("config TestLoadConfigForceEnv", t, tests.TestStartUp(func() {
		convey.So(env, convey.ShouldEqual, socconsts.EnvUnitTest)
	}, func() {
		LoadConfigForceEnv(tests.ConfigPath, socconsts.EnvUnitTest)
	}))
}
