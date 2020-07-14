package tests

import (
	"github.com/treeyh/soc-go-common/core/logger"
	"github.com/treeyh/soc-go-common/core/utils/file"
)

var (
	ConfigPath = file.GetCurrentPath() + "/../tests/config/"
)

var (
	Log = logger.Logger()
)

func TestStartUp(testFunc func(), initFunc func()) func() {

	return func() {
		//加载测试配置
		if initFunc != nil {
			initFunc()
			//time.Sleep(time.Duration(rand.Int63n(1000)) * time.Millisecond)
		}

		//丢进来的方法立刻执行
		testFunc()
	}
}
