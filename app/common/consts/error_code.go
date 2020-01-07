package consts

import "github.com/treeyh/soc-go-common/core/errors"

var (
	// PARSE_PARAM_ERROR 获取参数失败
	PARSE_PARAM_ERROR = errors.NewResultCode(200001, "%s 获取参数失败")
)
