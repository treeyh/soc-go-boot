package model

import "reflect"

// InParamsType 输入参数类型定义
type InParamsType struct {
	Name string `json:"name"`
	ParamsType
}

// OutParamsType 输出参数类型定义
type ParamsType struct {
	IsPointer bool         `json:"isPointer"`
	Type      reflect.Type `json:"type"`
}

type RouteMethod struct {
	PreUrl  string   `json:"preUrl"`
	Route   string   `json:"route"`
	Methods []string `json:"methods"`
}

// HandlerFuncInOut handler方法输入输出参数定义
type HandlerFuncInOut struct {
	Name         string          `json:"name"`
	RouteMethods *[]RouteMethod  `json:"routeMethods"`
	Ins          *[]InParamsType `json:"ins"`
	Outs         *[]ParamsType   `json:"outs"`
}

// HandlerFuncRoute 路由策略
type HandlerFuncRoute struct {
	PreUrl            string              `json:"preUrl"`
	HandlerFuncInOuts *[]HandlerFuncInOut `json:"handlerFuncInOuts"`
}
