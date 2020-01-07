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
	groupRouteMap     map[string]*gin.RouterGroup

	routeRegex = regexp.MustCompile(`@router\s+(\S+)(?:\s+\[(\S+)\])?`)

	demoString = "aaabbb"

	goTemplate = `package route


func init(){
	demoString = "{{.globalinfo}}"
}`
)

//func DemoPrint() {
//	path := file.GetCurrentPath()
//
//	fmt.Println(filepath.Join(path, "abc.go"))
//	f, err := os.Create(filepath.Join(path, "abc.go"))
//	if err != nil {
//		panic(err)
//	}
//	defer f.Close()
//
//	content := strings.ReplaceAll(goTemplate, "{{.globalinfo}}", "cccccc")
//	f.WriteString(content)
//
//	fmt.Println(demoString)
//}

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

	registerRoute(engine, &controller.UserController{})

	//engine.Any("/soc-go-boot-api/user/:create", controller.Create())
}

func registerRoute(engine *gin.Engine, contrs ...controller.IController) {
	//buildRouteMap(contrs...)

	groupRouteMapTmp := make(map[string]*gin.RouterGroup)

	for preUrl, suffixUrlMethodMap := range routeUrlMethodMap {
		groupRouteMapTmp[preUrl] = engine.Group(config.GetSocConfig().App.Server.ContextPath + preUrl)
		for suffixUrl, methodMap := range suffixUrlMethodMap {
			for method, funcInOutKeys := range methodMap {
				if "*" == method {
					groupRouteMapTmp[preUrl].Any(suffixUrl, buildHandler(method, suffixUrl, *getHandlerFuncInOutsByKey(&funcInOutKeys))...)
					continue
				}
				groupRouteMapTmp[preUrl].Handle(method, suffixUrl, buildHandler(method, suffixUrl, *getHandlerFuncInOutsByKey(&funcInOutKeys))...)
			}
		}
	}
}

func getHandlerFuncInOutsByKey(keys *[]string) *[]model.HandlerFuncInOut {
	funcs := make([]model.HandlerFuncInOut, 0, len(*keys))
	for _, v := range *keys {
		funcs = append(funcs, handlerFuncMap[v])
	}
	return &funcs
}

// buildHandler 构造 处理handler
func buildHandler(method, suffixUrl string, handlerFuncs []model.HandlerFuncInOut) []gin.HandlerFunc {

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
		for i, inParam := range *handlerFunc.Ins {
			elem := targetType.In(i)
			isPtr := elem.Kind() == reflect.Ptr
			inParam.IsPointer = isPtr
			if isPtr {
				inParam.Kind = elem.Elem().Kind()
				inParam.Type = elem.Elem()
			} else {
				inParam.Kind = elem.Kind()
				inParam.Type = elem
			}
			//fmt.Println(inParam.Name)
			//fmt.Println(inParam.Kind.String())
			//fmt.Println(inParam.Type.String())
			//fmt.Println("====")
			if i == 0 {
				if i == 0 && (!inParam.IsPointer || inParam.Kind.String() != "struct" || inParam.Type.String() != "req.GinContext") {
					panic(methodName + " The first parameter needs to be *gin.Context. ")
				}
				inParam.AssignType = model.UnAssign
			} else if checkParamExistUrl(&urlPaths, inParam.Name) {
				inParam.AssignType = model.PathAssign
			} else if method == "GET" || i < maxIndex {
				inParam.AssignType = model.QueryAssign
			} else if inParam.Kind.String() == "struct" && inParam.Type.String() != "time.Time" {
				inParam.AssignType = model.BodyAssign
			} else {
				inParam.AssignType = model.QueryAssign
			}
		}
		handlerFunc.InCount = maxIndex + 1

		// 构建输出参数
		for i, outParam := range *handlerFunc.Outs {
			elem := targetType.Out(i)
			isPtr := elem.Kind() == reflect.Ptr
			outParam.IsPointer = isPtr
			if isPtr {
				outParam.Kind = elem.Elem().Kind()
				outParam.Type = elem.Elem()
			} else {
				outParam.Kind = elem.Kind()
				outParam.Type = elem.Elem()
			}
			if i == 0 && (!isPtr || outParam.Kind.String() != "struct" || outParam.Type.String() != "resp.HttpRespResult") {
				panic(methodName + " The return value is only one and must be *resp.HttpRespResult. ")
			}
		}
		handlerFunc.InCount = len(*handlerFunc.Outs)
		handlers = append(handlers, httpHandler(&handlerFunc))
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

func httpHandler(handlerFunc *model.HandlerFuncInOut) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		ginContext := req.GinContext{Ctx: ctx}
		respData, err := injectFunc(&ginContext, handlerFunc)

		var respObj resp.RespResult
		if err != nil {
			logger.Logger().Error(err)
			respObj = resp.RespResult{
				Code:    err.Code(),
				Message: err.Message(),
			}
		} else {
			if len(*respData) > 0 {
				respObj = (*respData)[0].Interface().(resp.RespResult)
			}
		}
		resp.JsonRespResult(&ginContext, &respObj)
	}
}
