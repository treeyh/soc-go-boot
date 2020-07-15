package consts_error

import "github.com/treeyh/soc-go-common/core/errors"

var (
	// ParseParamError 获取参数失败
	ParseParamError = errors.NewResultCode(201001, "%s 获取参数失败")
)
