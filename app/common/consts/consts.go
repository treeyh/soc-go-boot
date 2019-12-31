package consts

// Trace相关Context的key
const (
	// TraceIdKey 用于http header
	TraceIdKey = "SOC-TRACE-ID"

	TracerContextKey = "SOC-Trace"

	// TracerHttpContextKey httpContext的ContextKey
	TracerHttpContextKey = "SOC-HttpContext"

	TraceParentSpanContextKey = "SOC-ParentSpanContext"
)
