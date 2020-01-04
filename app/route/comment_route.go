package route

import (
	"github.com/treeyh/soc-go-boot/app/model"
	"reflect"
	"time"
)

func init() {
	routeMap["/user"] = map[string]model.HandlerFuncInOut{}

	handlerFunc1 := model.HandlerFuncInOut{
		Name: "Create",
		RouteMethods: &[]model.RouteMethod{
			{
				Route: "/create",
				Methods: []string{
					"get",
				},
			},
		},
		Ins: &[]model.InParamsType{
			{
				Name: "createTime",
				ParamsType: model.ParamsType{
					IsPointer: true,
					Type:      reflect.TypeOf(time.Now()),
				},
			},
		},
		Outs: &[]model.ParamsType{
			{
				IsPointer: true,
				Type:      reflect.TypeOf(time.Now()),
			},
		},
	}
	routeMap["/user"]["method"] = handlerFunc1

}
