package route

import (
	"github.com/treeyh/soc-go-boot/app/controller"
	"github.com/treeyh/soc-go-boot/app/model"
)

func init() {
	routeUrlMethodMapTmp := make(map[string]map[string]map[string][]string)

	routeUrlMethodMapTmp["/user"] = make(map[string]map[string][]string)
	routeUrlMethodMapTmp["/user"]["/create"] = make(map[string][]string)
	routeUrlMethodMapTmp["/user"]["/create"]["*"] = make([]string, 1)
	routeUrlMethodMapTmp["/user"]["/create"]["*"][0] = "UserController.Create"

	routeUrlMethodMap = routeUrlMethodMapTmp

	userController := &controller.UserController{}

	handlerFuncMapTmp := make(map[string]model.HandlerFuncInOut)
	handlerFuncMapTmp["UserController.Create"] = model.HandlerFuncInOut{
		ControllerName: "UserController",
		Name:           "Create",
		RouteMethods:   nil,
		Ins: &[]model.InParamsType{
			{
				Name: "ctx",
				ParamsType: model.ParamsType{
					IsPointer:  false,
					DefaultVal: "",
					IsNeed:     false,
					Type:       nil,
					Kind:       0,
				},
			},
			{
				Name:       "updateTime",
				AssignType: 2,
				ParamsType: model.ParamsType{
					IsPointer:  false,
					DefaultVal: "2012-12-12 12:12:11",
					IsNeed:     true,
					Type:       nil,
					Kind:       0,
				},
			},
			{
				Name:       "createTime",
				AssignType: 2,
				ParamsType: model.ParamsType{
					IsPointer:  false,
					DefaultVal: "2012-12-12 12:12:11",
					IsNeed:     true,
					Type:       nil,
					Kind:       0,
				},
			},
			{
				Name:       "userId",
				AssignType: 2,
				ParamsType: model.ParamsType{
					IsPointer:  false,
					DefaultVal: "12",
					IsNeed:     true,
					Type:       nil,
					Kind:       0,
				},
			},
			{
				Name:       "userName",
				AssignType: 2,
				ParamsType: model.ParamsType{
					IsPointer:  false,
					DefaultVal: "12",
					IsNeed:     true,
					Type:       nil,
					Kind:       0,
				},
			},
			{
				Name:       "userReq",
				AssignType: 4,
				ParamsType: model.ParamsType{
					IsPointer:  true,
					DefaultVal: "",
					IsNeed:     true,
					Type:       nil,
					Kind:       0,
				},
			},
		},
		InCount: 6,
		Outs: &[]model.ParamsType{
			{},
		},
		OutCount: 1,
		Func:     userController.Create,
	}
	handlerFuncMap = handlerFuncMapTmp

}
