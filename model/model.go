package model

import (
	"context"
	"fmt"
	"github.com/treeyh/soc-go-boot/common/boot_consts"
	"github.com/treeyh/soc-go-common/core/consts"
	"github.com/treeyh/soc-go-common/core/utils/network"
	"github.com/treeyh/soc-go-common/core/utils/uuid"
	"net/http"
	"reflect"
)

const (
	UnAssign   HttpParamsAssignType = 0
	PathAssign HttpParamsAssignType = iota
	QueryAssign
	PostFormAssign
	BodyAssign
	HeaderAssign
)

var (
	httpContent HttpContext
)

type HttpParamsAssignType int
type RouteReqContentType string
type RouteRespContentType string

const (
	ReqContentTypeJson     RouteReqContentType = "json"
	ReqContentTypeXml      RouteReqContentType = "xml"
	ReqContentTypeProtoBuf RouteReqContentType = "protobuf"
	ReqContentTypeFile     RouteReqContentType = "file"

	RespContentTypeJson     RouteRespContentType = "json"
	RespContentTypeText     RouteRespContentType = "text"
	RespContentTypeXml      RouteRespContentType = "xml"
	RespContentTypeProtoBuf RouteRespContentType = "protobuf"
	RespContentTypeFile     RouteRespContentType = "file"
	RespContentTypeHtml     RouteRespContentType = "html"
	RespContentTypeRedirect RouteRespContentType = "redirect"
	RespContentTypeImage    RouteRespContentType = "image"
	RespContentTypeVideo    RouteRespContentType = "video"
	RespContentTypeAudio    RouteRespContentType = "audio"
)

// InParamsType 输入参数类型定义
type InParamsType struct {
	Name       string               `json:"name"`
	AssignType HttpParamsAssignType `json:"assignType"`
	ParamsType
}

// ParamsType 输出参数类型定义
type ParamsType struct {
	IsPointer  bool         `json:"isPointer"`
	DefaultVal string       `json:"defaultVal"`
	IsNeed     bool         `json:"isNeed"`
	Type       reflect.Type `json:"type"`
	Kind       reflect.Kind `json:"kind"`
}

type RouteMethod struct {
	PreUrl          string               `json:"preUrl"`
	Route           string               `json:"route"`
	Methods         []string             `json:"methods"`
	ReqContentType  RouteReqContentType  `json:"reqContentType"`
	RespContentType RouteRespContentType `json:"respContentType"`
}

// HandlerFuncInOut handler方法输入输出参数定义
type HandlerFuncInOut struct {
	ControllerName string         `json:"controllerName"`
	Name           string         `json:"name"`
	RouteMethods   []RouteMethod  `json:"routeMethods"`
	Ins            []InParamsType `json:"ins"`
	InCount        int            `json:"inCount"`
	Outs           []ParamsType   `json:"outs"`
	OutCount       int            `json:"outCount"`
	Func           interface{}    `json:"-"`
}

// HandlerFuncRoute 路由策略
type HandlerFuncRoute struct {
	PreUrl            string             `json:"preUrl"`
	HandlerFuncInOuts []HandlerFuncInOut `json:"handlerFuncInOuts"`
}

func (h HttpParamsAssignType) String() string {
	switch h {
	case UnAssign:
		return "UnAssign"
	case PathAssign:
		return "PathAssign"
	case QueryAssign:
		return "QueryAssign"
	case PostFormAssign:
		return "PostFormAssign"
	case BodyAssign:
		return "BodyAssign"
	case HeaderAssign:
		return "HeaderAssign"
	default:
		return "UNKNOWN"
	}
}

type HttpContext struct {
	Request       *http.Request
	Url           string
	Method        string
	StartTime     int64
	EndTime       int64
	TraceId       string
	SpanId        string
	Ip            string
	Status        int
	PartnerCode   string
	App           string
	Body          string
	AuthToken     string
	ClientVersion string
	Platform      string
	Channel       string
	Lang          string
}

// GetNewContext 获取一个新的ctx
func GetNewContext() context.Context {
	ctx := context.Background()
	traceId := fmt.Sprintf("%s_%s", network.GetIntranetIp(), uuid.NewUuid())
	ctx = context.WithValue(ctx, consts.ContextTracerKey, traceId)
	return ctx
}

func GetHttpContext(ctx context.Context) *HttpContext {
	val := ctx.Value(boot_consts.ContextHttpContextKey)
	if val == nil {
		return &httpContent
	}
	return val.(*HttpContext)
}
