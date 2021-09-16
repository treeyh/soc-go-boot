package boot_error_consts

import "github.com/treeyh/soc-go-common/core/errors"

var (
	// ControllerMethodError 方法不符合规范
	ControllerMethodError = errors.NewResultCode(201000, " Controller 方法不符合规范 %s ")

	// ParseParamError 获取参数失败
	ParseParamError = errors.NewResultCode(201001, "%s 获取参数失败")

	// SkyWalkingNotInit SkyWalking未初始化
	SkyWalkingNotInit = errors.NewResultCode(201051, "SkyWalking未初始化")

	// SignKeyNotExist 签名key不存在
	SignKeyNotExist = errors.NewResultCode(201061, "签名key不存在")

	// SignAuthFail 签名认证失败
	SignAuthFail = errors.NewResultCode(201062, "签名认证失败")

	// RequestTimestampOverLimit 请求时间戳超过阈值
	RequestTimestampOverLimit = errors.NewResultCode(201063, "请求时间戳超过阈值")
)
