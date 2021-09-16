package boot_consts

// Trace相关Context的key
const (

	// ContextHttpContextKey httpContext的ContextKey
	ContextHttpContextKey = "soc_http_context"
)

// 相关HTTP头
const (
	// HeaderAuthTokenKey 授权token的key
	HeaderAuthTokenKey = "soc-auth-token"
	// HeaderPlatform 平台
	HeaderPlatform = "soc-platform"
	// HeaderApp 应用
	HeaderApp = "soc-app"
	// HeaderChannel 渠道
	HeaderChannel       = "soc-channel"
	HeaderClientVersion = "soc-client-version"

	// HeaderSignKey 请求签名 http header
	HeaderSignKey = "soc-sign"

	// HeaderTimestampKey 请求时间戳 http header
	HeaderTimestampKey = "soc-timestamp"

	// HeaderSignPolicyKey 请求签名策略 http header
	HeaderSignPolicyKey = "soc-sign-policy"
)

const (
	// SignPolicySha256 签名策略 sha256
	SignPolicySha256 = "sha256"

	// SignPolicyMd5 签名策略 md5
	SignPolicyMd5 = "md5"
)
