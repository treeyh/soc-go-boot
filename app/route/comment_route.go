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
				Name:       "userId",
				ParamsType: model.ParamsType{},
			},
			{
				Name:       "userName",
				ParamsType: model.ParamsType{},
			},
			{
				Name:       "userReq",
				ParamsType: model.ParamsType{},
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

	routeUrlMethodMapTmp := make(map[string]map[string]map[string][]string)

	routeUrlMethodMapTmp["/user"] = make(map[string]map[string][]string)
	routeUrlMethodMapTmp["/user"]["/create"] = make(map[string][]string)
	routeUrlMethodMapTmp["/user"]["/create"]["*"] = make([]string, 1)
	routeUrlMethodMapTmp["/user"]["/create"]["*"][0] = "UserController.Create"

	routeUrlMethodMap = routeUrlMethodMapTmp

}
