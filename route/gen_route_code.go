package route

import (
	"github.com/treeyh/soc-go-boot/model"
	"github.com/treeyh/soc-go-common/core/consts"
	"github.com/treeyh/soc-go-common/core/utils/file"
	"github.com/treeyh/soc-go-common/core/utils/strs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	globalRouterTemplate = `package {{.routersDir}}

import (
	"{{.appModuleName}}/internal/controller"
	"github.com/treeyh/soc-go-boot/model"
)

var (

	// routeUrlMethodMap 路由urlmap string:controller的preurl； string:suffixUrl后缀url ； string:method； []string: 匹配到的handlerFuncMap的key
	routeUrlMethodMap map[string]map[string]map[string][]string

	// handlerFuncMap string: routeUrlMethodMap指定的url, 
	handlerFuncMap    map[string]model.HandlerFuncInOut

)

func init() {
	{{.routeCodeInfo}}
}

`

	handlerFuncMapTemplate = `model.HandlerFuncInOut{
		ControllerName: "{{.controllerName}}",
		Name:           "{{.methodName}}",
		RouteMethods:   []model.RouteMethod{
{{.routeMethods}}
		},
		Ins: []model.InParamsType{
{{.inParamsCode}}
		},
		InCount: {{.inParamsCount}},
		Outs: []model.ParamsType{
{{.outParamsCode}}
		},
		OutCount: {{.outParamsCount}},
		Func:     {{.controllerNameMethod}},
	}
`

	handlerFuncInParamTemplate = `			{
				Name:       "{{.inParamName}}",
				AssignType: {{.assignType}},
				ParamsType: model.ParamsType{
					IsPointer:  {{.isPointer}},
					DefaultVal: "{{.defaultVal}}",
					IsNeed:     {{.isNeed}},
					Type:       nil,
					Kind:       0,
				},
			},`

	handlerFuncOutParamTemplate = `			{},`

	handlerFuncRouteMethodsTemplate = `			{
				PreUrl: "{{.preUrl}}",
				Route: "{{.route}}",
				Methods: []string{
{{.methods}}
				},
				ReqContentType: {{.reqContentType}},
				RespContentType: {{.respContentType}},
			},`

	StringSlice = `					"{{.str}}",`

	blankStr = "    "

	genFileName = "comment_route_gen.go"
)

// genRouterCode 构建路由初始化代码
func genRouterCode(genPath, moduleName string, buildRouteMethodMap map[string]model.HandlerFuncInOut, buildRoutePathMap map[string]map[string]map[string][]string) {

	filePath := filepath.Join(genPath, genFileName)

	if file.ExistFile(filePath) {
		os.Remove(filePath)
	}

	pathMapCode := genPathMapCode(buildRoutePathMap)
	funcListCode := genHandlerFuncMapCode(buildRouteMethodMap)

	content := strings.ReplaceAll(globalRouterTemplate, "{{.routersDir}}", "route")
	content = strings.ReplaceAll(content, "{{.appModuleName}}", moduleName)

	content = strings.ReplaceAll(content, "{{.routeCodeInfo}}", blankStr+""+consts.LineSep+pathMapCode+consts.LineSep+consts.LineSep+funcListCode)

	file.WriteFile(filePath, content)

	time.Sleep(3 * time.Second)

	log.Info("Routing code regenerated, please restart ......")

	os.Exit(0)
}

// genPathMapCode 构造url path 的map代码
func genPathMapCode(buildRoutePathMap map[string]map[string]map[string][]string) string {

	routeUrlMethodMapLines := make([]string, 0)

	routeUrlMethodMapLines = append(routeUrlMethodMapLines, blankStr+"routeUrlMethodMapTmp := make(map[string]map[string]map[string][]string)", consts.EmptyStr)
	for preUrl, preRouteMap := range buildRoutePathMap {
		routeUrlMethodMapLines = append(routeUrlMethodMapLines, blankStr+"routeUrlMethodMapTmp[\""+preUrl+"\"] = make(map[string]map[string][]string)")

		for route, methodMap := range preRouteMap {
			routeUrlMethodMapLines = append(routeUrlMethodMapLines, blankStr+"routeUrlMethodMapTmp[\""+preUrl+"\"][\""+route+"\"] = make(map[string][]string)")
			for method, handlerNames := range methodMap {
				routeUrlMethodMapLines = append(routeUrlMethodMapLines, blankStr+"routeUrlMethodMapTmp[\""+preUrl+"\"][\""+route+"\"][\""+method+"\"] = make([]string, "+strconv.Itoa(len(handlerNames))+")")
				for i, handlerName := range handlerNames {
					routeUrlMethodMapLines = append(routeUrlMethodMapLines, blankStr+"routeUrlMethodMapTmp[\""+preUrl+"\"][\""+route+"\"][\""+method+"\"]["+strconv.Itoa(i)+"] = \""+handlerName+"\"")
				}
			}
		}
		routeUrlMethodMapLines = append(routeUrlMethodMapLines, consts.EmptyStr)
	}
	routeUrlMethodMapLines = append(routeUrlMethodMapLines, blankStr+"routeUrlMethodMap = routeUrlMethodMapTmp", consts.EmptyStr) // "+suffix+"
	return strings.Join(routeUrlMethodMapLines, consts.LineSep)
}

// genHandlerFuncMapCode 构造handler方法map的代码
func genHandlerFuncMapCode(buildRouteMethodMap map[string]model.HandlerFuncInOut) string {

	handlerFuncsCode := make([]string, 0)
	handlerFuncsCode = append(handlerFuncsCode, blankStr+"handlerFuncMapTmp := make(map[string]model.HandlerFuncInOut)", consts.EmptyStr)

	controllerNameMap := make(map[string]string)

	for k, handlerFunc := range buildRouteMethodMap {

		controllerInstanceName := strs.LcFirst(handlerFunc.ControllerName)
		controllerNameMap[controllerInstanceName] = handlerFunc.ControllerName

		routeMethodsCode := make([]string, 0)
		for _, routeMethod := range handlerFunc.RouteMethods {
			code := strings.ReplaceAll(handlerFuncRouteMethodsTemplate, "{{.preUrl}}", routeMethod.PreUrl)
			code = strings.ReplaceAll(code, "{{.route}}", routeMethod.Route)
			code = strings.ReplaceAll(code, "{{.methods}}", getStringByStringSlice(routeMethod.Methods))
			code = strings.ReplaceAll(code, "{{.reqContentType}}", getRouteReqContentTypeString(routeMethod.ReqContentType))
			code = strings.ReplaceAll(code, "{{.respContentType}}", getRouteRespContentTypeString(routeMethod.RespContentType))
			routeMethodsCode = append(routeMethodsCode, code)
		}

		inParamsCode := make([]string, 0)

		for _, inParam := range handlerFunc.Ins {
			code := strings.ReplaceAll(handlerFuncInParamTemplate, "{{.inParamName}}", inParam.Name)
			code = strings.ReplaceAll(code, "{{.assignType}}", strconv.Itoa(int(inParam.AssignType)))
			code = strings.ReplaceAll(code, "{{.isPointer}}", strconv.FormatBool(inParam.IsPointer))
			code = strings.ReplaceAll(code, "{{.defaultVal}}", strings.ReplaceAll(inParam.DefaultVal, "\"", "\\\""))
			code = strings.ReplaceAll(code, "{{.isNeed}}", strconv.FormatBool(inParam.IsNeed))
			inParamsCode = append(inParamsCode, code)
		}

		onParamsCode := make([]string, 0)

		outCount := len(handlerFunc.Outs)
		for i := 0; i < outCount; i++ {
			onParamsCode = append(onParamsCode, handlerFuncOutParamTemplate)
		}

		code := strings.ReplaceAll(handlerFuncMapTemplate, "{{.controllerName}}", handlerFunc.ControllerName)
		code = strings.ReplaceAll(code, "{{.methodName}}", handlerFunc.Name)
		code = strings.ReplaceAll(code, "{{.routeMethods}}", strings.Join(routeMethodsCode, consts.LineSep))
		code = strings.ReplaceAll(code, "{{.inParamsCode}}", strings.Join(inParamsCode, consts.LineSep))
		code = strings.ReplaceAll(code, "{{.inParamsCount}}", strconv.Itoa(handlerFunc.InCount))
		code = strings.ReplaceAll(code, "{{.outParamsCode}}", strings.Join(onParamsCode, consts.LineSep))
		code = strings.ReplaceAll(code, "{{.outParamsCount}}", strconv.Itoa(handlerFunc.OutCount))
		code = strings.ReplaceAll(code, "{{.controllerNameMethod}}", controllerInstanceName+"."+handlerFunc.Name)

		handlerFuncsCode = append(handlerFuncsCode, blankStr+"handlerFuncMapTmp[\""+k+"\"] = "+code)
	}

	controllerInstanceCode := make([]string, 0)
	for k, v := range controllerNameMap {
		controllerInstanceCode = append(controllerInstanceCode, blankStr+k+" := &controller."+v+"{}")
	}

	genCode := strings.Join(controllerInstanceCode, consts.LineSep) + consts.LineSep + strings.Join(handlerFuncsCode, consts.LineSep) + consts.LineSep + blankStr + "handlerFuncMap = handlerFuncMapTmp" //" + suffix + "

	return genCode
}

// getRouteReqContentTypeString 构造routeReqContentType字符串
func getRouteReqContentTypeString(contentType model.RouteReqContentType) string {
	switch contentType {
	case model.ReqContentTypeJson:
		return "model.ReqContentTypeJson"
	case model.ReqContentTypeProtoBuf:
		return "model.ReqContentTypeProtoBuf"
	case model.ReqContentTypeFile:
		return "model.ReqContentTypeFile"
	default:
		return "model.ReqContentTypeXml"
	}
}

// getRouteRespContentTypeString 构造routeRespContentType字符串
func getRouteRespContentTypeString(contentType model.RouteRespContentType) string {
	switch contentType {
	case model.RespContentTypeJson:
		return "model.RespContentTypeJson"
	case model.RespContentTypeText:
		return "model.RespContentTypeText"
	case model.RespContentTypeProtoBuf:
		return "model.RespContentTypeProtoBuf"
	case model.RespContentTypeFile:
		return "model.RespContentTypeFile"
	case model.RespContentTypeHtml:
		return "model.RespContentTypeHtml"
	case model.RespContentTypeRedirect:
		return "model.RespContentTypeRedirect"
	case model.RespContentTypeImage:
		return "model.RespContentTypeImage"
	case model.RespContentTypeVideo:
		return "model.RespContentTypeVideo"
	case model.RespContentTypeAudio:
		return "model.RespContentTypeAudio"
	default:
		return "model.RespContentTypeXml"
	}
}

// getStrings 构造字符串数组字符串
func getStringByStringSlice(strs []string) string {
	strss := make([]string, 0)
	for _, v := range strs {
		strss = append(strss, strings.ReplaceAll(StringSlice, "{{.str}}", v))
	}
	return strings.Join(strss, consts.LineSep)
}
