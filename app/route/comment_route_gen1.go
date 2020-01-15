package route

import (
	"github.com/treeyh/soc-go-boot/app/controller"
	"github.com/treeyh/soc-go-boot/app/model"
)

func init() {
	routeCodeMd5 = "125ce18e142df552252149ab47b5bf20"

}

func initRouteInfo() {
	routeCodeMd5 = "125ce18e142df552252149ab47b5bf20"

	routeUrlMethodMapTmp := make(map[string]map[string]map[string][]string)

	routeUrlMethodMapTmp["/user"] = make(map[string]map[string][]string)
	routeUrlMethodMapTmp["/user"]["/get/:userId"] = make(map[string][]string)
	routeUrlMethodMapTmp["/user"]["/get/:userId"]["GET"] = make([]string, 1)
	routeUrlMethodMapTmp["/user"]["/get/:userId"]["GET"][0] = "UserController.Get"
	routeUrlMethodMapTmp["/user"]["/get/:userId"]["POST"] = make([]string, 1)
	routeUrlMethodMapTmp["/user"]["/get/:userId"]["POST"][0] = "UserController.Get"
	routeUrlMethodMapTmp["/user"]["/create"] = make(map[string][]string)
	routeUrlMethodMapTmp["/user"]["/create"]["*"] = make([]string, 1)
	routeUrlMethodMapTmp["/user"]["/create"]["*"][0] = "UserController.Create"

	routeUrlMethodMap1 = routeUrlMethodMapTmp

	userController := &controller.UserController{}
	handlerFuncMapTmp := make(map[string]model.HandlerFuncInOut)

	handlerFuncMapTmp["UserController.Get"] = model.HandlerFuncInOut{
		ControllerName: "UserController",
		Name:           "Get",
		RouteMethods:   nil,
		Ins: &[]model.InParamsType{
			{
				Name:       "ctx",
				AssignType: 0,
				ParamsType: model.ParamsType{
					IsPointer:  false,
					DefaultVal: "",
					IsNeed:     false,
					Type:       nil,
					Kind:       0,
				},
			},
			{
				Name:       "userId",
				AssignType: 2,
				ParamsType: model.ParamsType{
					IsPointer:  false,
					DefaultVal: "1",
					IsNeed:     true,
					Type:       nil,
					Kind:       0,
				},
			},
		},
		InCount: 0,
		Outs: &[]model.ParamsType{
			{},
		},
		OutCount: 0,
		Func:     userController.Get,
	}

	handlerFuncMapTmp["UserController.Create"] = model.HandlerFuncInOut{
		ControllerName: "UserController",
		Name:           "Create",
		RouteMethods:   nil,
		Ins: &[]model.InParamsType{
			{
				Name:       "ctx",
				AssignType: 0,
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
					DefaultVal: "1",
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
					DefaultVal: "2012-12-12 12:12:11",
					IsNeed:     false,
					Type:       nil,
					Kind:       0,
				},
			},
			{
				Name:       "userReq",
				AssignType: 0,
				ParamsType: model.ParamsType{
					IsPointer:  false,
					DefaultVal: "",
					IsNeed:     false,
					Type:       nil,
					Kind:       0,
				},
			},
		},
		InCount: 0,
		Outs: &[]model.ParamsType{
			{},
		},
		OutCount: 0,
		Func:     userController.Create,
	}

	handlerFuncMap1 = handlerFuncMapTmp
}
