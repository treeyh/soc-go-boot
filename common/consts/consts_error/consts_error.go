package consts_error

import "github.com/treeyh/soc-go-common/core/errors"

var (
	// ControllerMethodError 方法不符合规范
	ControllerMethodError = errors.NewResultCode(201000, " Controller 方法不符合规范 %s ")

	// ParseParamError 获取参数失败
	ParseParamError = errors.NewResultCode(201001, "%s 获取参数失败")
)
