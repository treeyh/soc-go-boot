package route

import (
	"github.com/gin-gonic/gin"
	"github.com/treeyh/soc-go-boot/app/config"
	"github.com/treeyh/soc-go-boot/app/controller"
	"github.com/treeyh/soc-go-boot/app/model"
	"github.com/treeyh/soc-go-boot/app/model/req"
	"github.com/treeyh/soc-go-boot/app/model/resp"
	"github.com/treeyh/soc-go-common/core/logger"
	"reflect"
	"regexp"
	"strings"
)

const httpMethods = ",GET,POST,DELETE,PATCH,PUT,OPTIONS,HEAD,*,"

var (
	handlerFuncMap    map[string]model.HandlerFuncInOut
	routeUrlMethodMap map[string]map[string]map[string][]string

	groupRouteMap map[string]*gin.RouterGroup

	routeRegex = regexp.MustCompile(`@Router\s+(\S+)(?:\s+\[(\S+)\])?`)
)

func SetupRouter(engine *gin.Engine) {

	//engine.Use(gin.Logger())
	//engine.Use(gin.Recovery())
	//
	//engine.Use(exception.SetUp())
	//engine.Use(jaeger.SetUp())
	//
	////404
	//engine.NoRoute(func(c *gin.Context) {
	//	utilGin := resp.GinContext{Ctx: c}
	//	utilGin.Json(404, "请求方法不存在", nil)
	//})
	//
	//engine.GET("/sing", func(c *gin.Context) {
	//	utilGin := resp.GinContext{Ctx: c}
	//	utilGin.Json(200, "ok", nil)
	//})

	registerRoute(engine, &controller.UserController{}, &controller.ProjectController{})

	//engine.Any("/soc-go-boot-api/user/:create", controller.Create())
}

func registerRoute(engine *gin.Engine, contrs ...controller.IController) {
	buildRouteMap(contrs...)

	groupRouteMapTmp := make(map[string]*gin.RouterGroup)
	for preUrl, suffixUrlMethodMap := range routeUrlMethodMap {
		groupRouteMapTmp[preUrl] = engine.Group(config.GetSocConfig().App.Server.ContextPath + preUrl)
		for suffixUrl, methodMap := range suffixUrlMethodMap {
			for method, funcInOutKeys := range methodMap {
				if "*" == method {
					groupRouteMapTmp[preUrl].Any(suffixUrl, buildHandler(method, preUrl, suffixUrl, getHandlerFuncInOutsByKey(&funcInOutKeys))...)
					continue
				}
				groupRouteMapTmp[preUrl].Handle(method, suffixUrl, buildHandler(method, preUrl, suffixUrl, getHandlerFuncInOutsByKey(&funcInOutKeys))...)
			}
		}
	}
}

func getHandlerFuncInOutsByKey(keys *[]string) []model.HandlerFuncInOut {
	funcs := make([]model.HandlerFuncInOut, 0, len(*keys))
	for _, v := range *keys {
		funcs = append(funcs, handlerFuncMap[v])
	}
	return funcs
}

// buildHandler 构造 处理handler
func buildHandler(method, preUrl, suffixUrl string, handlerFuncs []model.HandlerFuncInOut) []gin.HandlerFunc {

	handlers := make([]gin.HandlerFunc, 0, len(handlerFuncs))

	for _, handlerFunc := range handlerFuncs {
		// 验证 handlerFunc 是否符合规范
		targetType := reflect.TypeOf(handlerFunc.Func)
		methodName := handlerFunc.ControllerName + "." + handlerFunc.Name
		if reflect.Func != targetType.Kind() {
			panic(" buildHandler " + methodName + " not func. ")
		}
		if targetType.NumIn() < 1 {
			panic(methodName + " The first parameter needs to be *gin.Context, the return value is only one and must be *resp.HttpRespResult. ")
		}
		if targetType.NumOut() != 1 {
			panic(methodName + " The first parameter needs to be *gin.Context, the return value is only one and must be *resp.HttpRespResult. ")
		}
		urlPaths := strings.Split(suffixUrl, "/")
		// 构建输入参数列表
		maxIndex := len(*handlerFunc.Ins) - 1
		for i, _ := range *handlerFunc.Ins {
			elem := targetType.In(i)
			isPtr := elem.Kind() == reflect.Ptr
			(*handlerFunc.Ins)[i].IsPointer = isPtr
			if isPtr {
				(*handlerFunc.Ins)[i].Kind = elem.Elem().Kind()
				(*handlerFunc.Ins)[i].Type = elem.Elem()
			} else {
				(*handlerFunc.Ins)[i].Kind = elem.Kind()
				(*handlerFunc.Ins)[i].Type = elem
			}
			if i == 0 {
				if i == 0 && (!(*handlerFunc.Ins)[i].IsPointer || (*handlerFunc.Ins)[i].Kind.String() != "struct" || (*handlerFunc.Ins)[i].Type.String() != "req.GinContext") {
					panic(methodName + " The first parameter needs to be *gin.Context. ")
				}
				(*handlerFunc.Ins)[i].AssignType = model.UnAssign
				continue
			}

			if (*handlerFunc.Ins)[i].AssignType == model.UnAssign {
				// 若没有指定获取方式，通过程序判定
				if checkParamExistUrl(&urlPaths, (*handlerFunc.Ins)[i].Name) {
					(*handlerFunc.Ins)[i].AssignType = model.PathAssign
				} else if method == "GET" || i < maxIndex {
					(*handlerFunc.Ins)[i].AssignType = model.QueryAssign
				} else if (*handlerFunc.Ins)[i].Kind.String() == "struct" && (*handlerFunc.Ins)[i].Type.String() != "time.Time" && (*handlerFunc.Ins)[i].Type.String() != "types.Time" {
					(*handlerFunc.Ins)[i].AssignType = model.BodyAssign
				} else {
					(*handlerFunc.Ins)[i].AssignType = model.QueryAssign
				}
			}
			//fmt.Println((*handlerFunc.Ins)[i].AssignType)
			//fmt.Println("====" + (*handlerFunc.Ins)[i].Kind.String())
			//fmt.Println("====" + (*handlerFunc.Ins)[i].Type.String())
		}
		handlerFunc.InCount = maxIndex + 1

		// 构建输出参数
		for i, _ := range *handlerFunc.Outs {
			elem := targetType.Out(i)
			isPtr := elem.Kind() == reflect.Ptr
			(*handlerFunc.Outs)[i].IsPointer = isPtr
			if isPtr {
				(*handlerFunc.Outs)[i].Kind = elem.Elem().Kind()
				(*handlerFunc.Outs)[i].Type = elem.Elem()
			} else {
				(*handlerFunc.Outs)[i].Kind = elem.Kind()
				(*handlerFunc.Outs)[i].Type = elem.Elem()
			}
			if i == 0 && (!isPtr || (*handlerFunc.Outs)[i].Kind.String() != "struct" || (*handlerFunc.Outs)[i].Type.String() != "resp.HttpRespResult") {
				panic(methodName + " The return value is only one and must be *resp.HttpRespResult. ")
			}
		}
		handlerFunc.OutCount = len(*handlerFunc.Outs)

		handlers = append(handlers, httpHandler(method, preUrl, suffixUrl, handlerFunc))
	}
	return handlers
}

// checkParamExistUrl 判断参数是否在url中获取
func checkParamExistUrl(urlPaths *[]string, param string) bool {
	param1 := "*" + param
	param2 := ":" + param
	for _, v := range *urlPaths {
		if v == param1 || v == param2 {
			return true
		}
	}
	return false
}

func httpHandler(method, preUrl, suffixUrl string, handlerFunc model.HandlerFuncInOut) gin.HandlerFunc {

	logger.Logger().Info("Handler path: " + method + " " + preUrl + suffixUrl + "      -->      " + handlerFunc.ControllerName + "." + handlerFunc.Name)

	return func(ctx *gin.Context) {

		ginContext := req.GinContext{Ctx: ctx}
		respData, err := injectFunc(&ginContext, handlerFunc)

		var respObj *resp.HttpRespResult
		if err != nil {
			logger.Logger().Error(err)
			respObj = &resp.HttpRespResult{
				HttpStatus: 500,
				RespResult: resp.RespResult{
					Code:    err.Code(),
					Message: err.Message(),
				},
			}
		} else {
			if len(respData) > 0 {
				respObj = (respData)[0].Interface().(*resp.HttpRespResult)
			}
		}
		resp.JsonHttpRespResult(&ginContext, respObj)
	}
}
