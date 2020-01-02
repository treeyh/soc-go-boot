package model

import "reflect"

// ParamsType 参数类型定义
type ParamsType struct {
	IsPointer bool         `json:"isPointer"`
	Type      reflect.Type `json:"type"`
}

// FuncInOut 方法输入输出参数定义
type FuncInOut struct {
	Name string        `json:"name"`
	Ins  *[]ParamsType `json:"ins"`
	Outs *[]ParamsType `json:"outs"`
}
