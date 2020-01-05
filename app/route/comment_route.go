package route

import (
	"github.com/treeyh/soc-go-boot/app/controller"
	"github.com/treeyh/soc-go-boot/app/model"
)

func init() {

	userController := &controller.UserController{}

	handlerFuncMap = make(map[string]model.HandlerFuncInOut)
	handlerFuncMap["UserController.Create"] = model.HandlerFuncInOut{
		ControllerName: "UserController",
		Name:           "Create",
		RouteMethods:   nil,
		Ins: &[]model.InParamsType{
			{
				Name:       "ctx",
				ParamsType: model.ParamsType{},
			},
			{
				Name:       "updateTime",
				ParamsType: model.ParamsType{},
			},
			{
				Name:       "createTime",
				ParamsType: model.ParamsType{},
			},
			{
				Name:       "userReq",
				ParamsType: model.ParamsType{},
			},
		},
		Outs: &[]model.ParamsType{},
		Func: userController.Create,
	}

	routeUrlMethodMap = make(map[string]map[string]map[string][]string)

	routeUrlMethodMap["/user"] = make(map[string]map[string][]string)
	routeUrlMethodMap["/user"]["/create"] = make(map[string][]string)
	routeUrlMethodMap["/user"]["/create"]["get"] = make([]string, 1)
	routeUrlMethodMap["/user"]["/create"]["get"][0] = "UserController.Create"

}
