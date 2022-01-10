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

	// HeaderAppCodeKey 用于http header
	HeaderAppCodeKey = "soc-app"

	// HeaderPartnerCodeKey 合作方id的http header
	HeaderPartnerCodeKey = "soc-partner-code"

	// HeaderChannel 渠道
	HeaderChannel = "soc-channel"

	HeaderClientVersion = "soc-client-version"

	// HeaderTraceIdKey 用于http header
	HeaderTraceIdKey = "soc-trace-id"

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

const (
	// LangZhCn LangZhChs 中文简体  中间件会统一转为 LangZhCn
	LangZhCn  string = "zh-CN"
	LangZhChs string = "zh-CHS"

	// LangZhTw LangZhHk LangZhMo LangZhCht 中文繁体 中间件会统一转为 LangZhTw
	LangZhTw  string = "zh-TW"
	LangZhHk  string = "zh-HK"
	LangZhMo  string = "zh-MO"
	LangZhCht string = "zh-CHT"

	// LangEn 英文
	LangEn string = "en"
)
