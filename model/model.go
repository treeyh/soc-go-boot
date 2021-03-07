package model

import "reflect"

type HttpParamsAssignType int

const (
	UnAssign   HttpParamsAssignType = 0
	PathAssign HttpParamsAssignType = iota
	QueryAssign
	PostFormAssign
	BodyAssign
	HeaderAssign
)

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

// OutParamsType 输出参数类型定义
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
