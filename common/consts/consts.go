package consts

// Trace相关Context的key
const (

	// TracerHttpContextKey httpContext的ContextKey
	TracerHttpContextKey = "SOC-HttpContext"
)

// 相关HTTP头
const (
	// 授权token的key
	HeaderAuthTokenKey  = "SOC-Auth-Token"
	HeaderPlatform      = "SOC-Platform"
	HeaderApp           = "SOC-App"
	HeaderClientVersion = "SOC-Client-Version"
)
