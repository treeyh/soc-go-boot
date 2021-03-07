package route

import (
	"github.com/gin-gonic/gin"
	"github.com/treeyh/soc-go-boot/common/consts/consts_error"
	"github.com/treeyh/soc-go-boot/config"
	"github.com/treeyh/soc-go-boot/controller"
	"github.com/treeyh/soc-go-boot/model"
	socreq "github.com/treeyh/soc-go-boot/model/req"
	"github.com/treeyh/soc-go-boot/model/resp"
	"github.com/treeyh/soc-go-common/core/consts"
	"github.com/treeyh/soc-go-common/core/errors"
	"reflect"
	"strings"
	"time"
)

var (
	groupRouteMap map[string]*gin.RouterGroup

	checkOutParamMsgMap = make(map[model.RouteRespContentType]string)

	responseFuncMap = make(map[model.RouteRespContentType]func(*socreq.GinContext, []reflect.Value, errors.AppError))
)

func init() {

	checkOutParamMsgMap[model.RespContentTypeFile] = "resp.HttpFileRespResult"
	checkOutParamMsgMap[model.RespContentTypeHtml] = "resp.HttpHtmlRespResult"
	checkOutParamMsgMap[model.RespContentTypeRedirect] = "resp.HttpRedirectRespResult"
	checkOutParamMsgMap[model.RespContentTypeXml] = "resp.HttpXmlRespResult"
	checkOutParamMsgMap[model.RespContentTypeJson] = "resp.HttpJsonRespResult"
	checkOutParamMsgMap[model.RespContentTypeProtoBuf] = "resp.HttpProtoBufRespResult"
	checkOutParamMsgMap[model.RespContentTypeText] = "resp.HttpTextRespResult"

	responseFuncMap[model.RespContentTypeFile] = responseFile
	responseFuncMap[model.RespContentTypeHtml] = responseHtml
	responseFuncMap[model.RespContentTypeRedirect] = responseRedirect
	responseFuncMap[model.RespContentTypeXml] = responseXml
	responseFuncMap[model.RespContentTypeJson] = responseJson
	responseFuncMap[model.RespContentTypeProtoBuf] = responseProtoBuf
	responseFuncMap[model.RespContentTypeText] = responseText
}

func RegisterRoute(engine *gin.Engine, controllerStatusPath, controllerPath, goModFilePath, genPath string, routeUrlMethodMap map[string]map[string]map[string][]string, handlerFuncMap map[string]model.HandlerFuncInOut, contrs ...controller.IController) {

	// 构造controllermap
	if consts.GetCurrentEnv() == consts.EnvLocal {
		BuildRouteMap(controllerStatusPath, controllerPath, goModFilePath, genPath, contrs...)
	}

	groupRouteMapTmp := make(map[string]*gin.RouterGroup)
	for preUrl, suffixUrlMethodMap := range routeUrlMethodMap {
		groupRouteMapTmp[preUrl] = engine.Group(config.GetSocConfig().App.Server.ContextPath + preUrl)
		for suffixUrl, methodMap := range suffixUrlMethodMap {
			for method, funcInOutKeys := range methodMap {
				if "*" == method {
					groupRouteMapTmp[preUrl].Any(suffixUrl, buildHandler(method, preUrl, suffixUrl, getHandlerFuncInOutsByKey(funcInOutKeys, handlerFuncMap))...)
					continue
				}
				groupRouteMapTmp[preUrl].Handle(method, suffixUrl, buildHandler(method, preUrl, suffixUrl, getHandlerFuncInOutsByKey(funcInOutKeys, handlerFuncMap))...)
			}
		}
	}
}

func getHandlerFuncInOutsByKey(keys []string, handlerFuncMap map[string]model.HandlerFuncInOut) []model.HandlerFuncInOut {
	funcs := make([]model.HandlerFuncInOut, 0, len(keys))
	for _, v := range keys {
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
			panic(methodName + " The first parameter needs to be *gin.Context, the return value is only one and must be *socresp.HttpJsonRespResult. ")
		}

		// 构建输出参数
		checkAndBuildOutParam(targetType, &handlerFunc)
		handlerFunc.OutCount = len(handlerFunc.Outs)

		urlPaths := strings.Split(suffixUrl, "/")
		// 构建输入参数列表
		maxIndex := len(handlerFunc.Ins) - 1
		for i, _ := range handlerFunc.Ins {
			elem := targetType.In(i)
			isPtr := elem.Kind() == reflect.Ptr
			(handlerFunc.Ins)[i].IsPointer = isPtr
			if isPtr {
				(handlerFunc.Ins)[i].Kind = elem.Elem().Kind()
				(handlerFunc.Ins)[i].Type = elem.Elem()
			} else {
				(handlerFunc.Ins)[i].Kind = elem.Kind()
				(handlerFunc.Ins)[i].Type = elem
			}
			if i == 0 {
				if i == 0 && (!(handlerFunc.Ins)[i].IsPointer || (handlerFunc.Ins)[i].Kind.String() != "struct" || (handlerFunc.Ins)[i].Type.String() != "req.GinContext") {
					panic(methodName + " The first parameter needs to be *gin.Context. ")
				}
				(handlerFunc.Ins)[i].AssignType = model.UnAssign
				continue
			}

			if (handlerFunc.Ins)[i].AssignType == model.UnAssign {
				// 若没有指定获取方式，通过程序判定
				if checkParamExistUrl(urlPaths, (handlerFunc.Ins)[i].Name) {
					(handlerFunc.Ins)[i].AssignType = model.PathAssign
				} else if method == "GET" || i < maxIndex {
					(handlerFunc.Ins)[i].AssignType = model.QueryAssign
				} else if (handlerFunc.Ins)[i].Kind.String() == "struct" && (handlerFunc.Ins)[i].Type.String() != "time.Time" && (handlerFunc.Ins)[i].Type.String() != "types.Time" {
					(handlerFunc.Ins)[i].AssignType = model.BodyAssign
				} else {
					(handlerFunc.Ins)[i].AssignType = model.QueryAssign
				}
			}
			//fmt.Println((*handlerFunc.Ins)[i].AssignType)
			//fmt.Println("====" + (*handlerFunc.Ins)[i].Kind.String())
			//fmt.Println("====" + (*handlerFunc.Ins)[i].Type.String())
		}
		handlerFunc.InCount = maxIndex + 1

		handlers = append(handlers, httpHandler(method, preUrl, suffixUrl, &handlerFunc))
	}
	return handlers
}

// checkOutParam 检查controller方法返回类型是否符合要求
func checkAndBuildOutParam(targetType reflect.Type, handlerFunc *model.HandlerFuncInOut) errors.AppError {

	// 构建输出参数
	for i, _ := range handlerFunc.Outs {
		elem := targetType.Out(i)
		isPtr := elem.Kind() == reflect.Ptr
		(handlerFunc.Outs)[i].IsPointer = isPtr
		if isPtr {
			handlerFunc.Outs[i].Kind = elem.Elem().Kind()
			handlerFunc.Outs[i].Type = elem.Elem()
		} else {
			handlerFunc.Outs[i].Kind = elem.Kind()
			handlerFunc.Outs[i].Type = elem.Elem()
		}
	}

	respContentType := handlerFunc.RouteMethods[0].RespContentType

	msg := string(respContentType) + " 输出类型只允许一个返回参数，且类型为 *" + checkOutParamMsgMap[respContentType]
	if len(handlerFunc.Outs) != 1 {
		return errors.NewAppError(consts_error.ControllerMethodError, msg)
	}

	if !handlerFunc.Outs[0].IsPointer || handlerFunc.Outs[0].Kind.String() != "struct" || handlerFunc.Outs[0].Type.String() != checkOutParamMsgMap[respContentType] {
		return errors.NewAppError(consts_error.ControllerMethodError, msg)
	}
	return nil
}

// checkParamExistUrl 判断参数是否在url中获取
func checkParamExistUrl(urlPaths []string, param string) bool {
	param1 := "*" + param
	param2 := ":" + param
	for _, v := range urlPaths {
		if v == param1 || v == param2 {
			return true
		}
	}
	return false
}

func httpHandler(method, preUrl, suffixUrl string, handlerFunc *model.HandlerFuncInOut) gin.HandlerFunc {

	log.Info("Handler path: " + method + " " + preUrl + suffixUrl + "      -->      " + handlerFunc.ControllerName + "." + handlerFunc.Name)

	respContentType := handlerFunc.RouteMethods[0].RespContentType

	return func(ctx *gin.Context) {

		ginContext := socreq.GinContext{Ctx: ctx}
		respData, err := InjectFunc(&ginContext, handlerFunc)

		if function, ok := responseFuncMap[respContentType]; ok {
			function(&ginContext, respData, err)
		}
	}
}

// responseJson 输出json
func responseJson(ctx *socreq.GinContext, results []reflect.Value, err errors.AppError) {
	var respObj *resp.HttpJsonRespResult

	if err != nil {
		log.ErrorCtx(ctx.Ctx.Request.Context(), err)
		respObj = &resp.HttpJsonRespResult{
			HttpStatus: 500,
			Data: resp.RespResult{
				Code:      err.Code(),
				Message:   err.Message(),
				Timestamp: time.Now().Unix(),
			},
		}
		controller.JsonHttpRespResult(ctx, respObj)
		return
	}
	if len(results) > 0 {
		respObj = (results)[0].Interface().(*resp.HttpJsonRespResult)
	}

	controller.JsonHttpRespResult(ctx, respObj)
}

// responseHtml
func responseHtml(ctx *socreq.GinContext, results []reflect.Value, err errors.AppError) {
	var respObj *resp.HttpHtmlRespResult

	if err != nil {
		log.ErrorCtx(ctx.Ctx.Request.Context(), err)
		respObj = &resp.HttpHtmlRespResult{
			HttpStatus: 500,
		}
		controller.HtmlHttpRespResult(ctx, respObj)
		return
	}
	if len(results) > 0 {
		respObj = (results)[0].Interface().(*resp.HttpHtmlRespResult)
	}

	controller.HtmlHttpRespResult(ctx, respObj)
}

// responseText
func responseText(ctx *socreq.GinContext, results []reflect.Value, err errors.AppError) {
	var respObj *resp.HttpTextRespResult

	if err != nil {
		log.ErrorCtx(ctx.Ctx.Request.Context(), err)
		respObj = &resp.HttpTextRespResult{
			HttpStatus: 500,
		}
		controller.TextHttpRespResult(ctx, respObj)
		return
	}
	if len(results) > 0 {
		respObj = (results)[0].Interface().(*resp.HttpTextRespResult)
	}

	controller.TextHttpRespResult(ctx, respObj)
}

// responseXml
func responseXml(ctx *socreq.GinContext, results []reflect.Value, err errors.AppError) {
	var respObj *resp.HttpXmlRespResult

	if err != nil {
		log.ErrorCtx(ctx.Ctx.Request.Context(), err)
		respObj = &resp.HttpXmlRespResult{
			HttpStatus: 500,
		}
		controller.XmlHttpRespResult(ctx, respObj)
		return
	}
	if len(results) > 0 {
		respObj = (results)[0].Interface().(*resp.HttpXmlRespResult)
	}

	controller.XmlHttpRespResult(ctx, respObj)
}

// responseHtml
func responseProtoBuf(ctx *socreq.GinContext, results []reflect.Value, err errors.AppError) {
	var respObj *resp.HttpProtoBufRespResult

	if err != nil {
		log.ErrorCtx(ctx.Ctx.Request.Context(), err)
		respObj = &resp.HttpProtoBufRespResult{
			HttpStatus: 500,
		}
		controller.ProtoBufHttpRespResult(ctx, respObj)
		return
	}
	if len(results) > 0 {
		respObj = (results)[0].Interface().(*resp.HttpProtoBufRespResult)
	}

	controller.ProtoBufHttpRespResult(ctx, respObj)
}

// responseHtml
func responseRedirect(ctx *socreq.GinContext, results []reflect.Value, err errors.AppError) {
	var respObj *resp.HttpRedirectRespResult

	if err != nil {
		log.ErrorCtx(ctx.Ctx.Request.Context(), err)
		respObj = &resp.HttpRedirectRespResult{
			HttpStatus: 500,
		}
		controller.RedirectHttpRespResult(ctx, respObj)
		return
	}
	if len(results) > 0 {
		respObj = (results)[0].Interface().(*resp.HttpRedirectRespResult)
	}

	controller.RedirectHttpRespResult(ctx, respObj)
}

// responseHtml
func responseFile(ctx *socreq.GinContext, results []reflect.Value, err errors.AppError) {
	var respObj *resp.HttpFileRespResult

	if err != nil {
		log.ErrorCtx(ctx.Ctx.Request.Context(), err)
		respObj = &resp.HttpFileRespResult{
			HttpStatus: 500,
		}
		controller.FileHttpRespResult(ctx, respObj)
		return
	}
	if len(results) > 0 {
		respObj = (results)[0].Interface().(*resp.HttpFileRespResult)
	}

	controller.FileHttpRespResult(ctx, respObj)
}
