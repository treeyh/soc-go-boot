package route

import (
	"github.com/treeyh/soc-go-boot/app/controller"
	"github.com/treeyh/soc-go-boot/app/model"
)

func init() {

	userController := &controller.UserController{}

	handlerFuncMapTmp := make(map[string]model.HandlerFuncInOut)
	handlerFuncMapTmp["UserController.Create"] = model.HandlerFuncInOut{
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
	handlerFuncMap = handlerFuncMapTmp

	routeUrlMethodMapTmp := make(map[string]map[string]map[string][]string)

	routeUrlMethodMapTmp["/user"] = make(map[string]map[string][]string)
	routeUrlMethodMapTmp["/user"]["/create"] = make(map[string][]string)
	routeUrlMethodMapTmp["/user"]["/create"]["GET"] = make([]string, 1)
	routeUrlMethodMapTmp["/user"]["/create"]["GET"][0] = "UserController.Create"

	routeUrlMethodMap = routeUrlMethodMapTmp

}
