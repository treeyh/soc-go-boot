package model

import "reflect"

type HttpParamsAssignType int

const (
	UrlAssign   HttpParamsAssignType = 1
	QueryAssign HttpParamsAssignType = iota
	BodyAssign
	HeaderAssign
)

// InParamsType 输入参数类型定义
type InParamsType struct {
	Name       string               `json:"name"`
	AssignType HttpParamsAssignType `json:"assignType"`
	ParamsType
}

// OutParamsType 输出参数类型定义
type ParamsType struct {
	IsPointer bool         `json:"isPointer"`
	Type      reflect.Type `json:"type"`
	Kind      reflect.Kind `json:"kind"`
}

type RouteMethod struct {
	PreUrl  string   `json:"preUrl"`
	Route   string   `json:"route"`
	Methods []string `json:"methods"`
}

// HandlerFuncInOut handler方法输入输出参数定义
type HandlerFuncInOut struct {
	ControllerName string          `json:"controllerName"`
	Name           string          `json:"name"`
	RouteMethods   *[]RouteMethod  `json:"routeMethods"`
	Ins            *[]InParamsType `json:"ins"`
	Outs           *[]ParamsType   `json:"outs"`
	Func           interface{}     `json:"-"`
}

// HandlerFuncRoute 路由策略
type HandlerFuncRoute struct {
	PreUrl            string              `json:"preUrl"`
	HandlerFuncInOuts *[]HandlerFuncInOut `json:"handlerFuncInOuts"`
}
