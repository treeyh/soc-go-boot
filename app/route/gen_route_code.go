package route

import (
	"fmt"
	"github.com/treeyh/soc-go-boot/app/common/consts"
	"github.com/treeyh/soc-go-boot/app/model"
	"github.com/treeyh/soc-go-common/core/logger"
	"github.com/treeyh/soc-go-common/core/utils/encrypt"
	"github.com/treeyh/soc-go-common/core/utils/file"
	"github.com/treeyh/soc-go-common/core/utils/json"
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
	"{{.appModuleName}}/app/controller"
	"github.com/treeyh/soc-go-boot/app/model"
)

func init() {
{{.md5CodeInfo}}
}

func initRouteInfo() {
{{.md5CodeInfo}}

{{.routeCodeInfo}}
}
`

	handlerFuncMapTemplate = `model.HandlerFuncInOut{
		ControllerName: "{{.controllerName}}",
		Name:           "{{.methodName}}",
		RouteMethods:   nil,
		Ins: &[]model.InParamsType{
{{.inParamsCode}}
		},
		InCount: {{.inParamsCount}},
		Outs: &[]model.ParamsType{
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

	blankStr = "    "

	genFileName1 = "comment_route_gen1.go"

	genFileName2 = "comment_route_gen2.go"
)

// genRouterCode 构建路由初始化代码
func genRouterCode(moduleName string, buildRouteMethodMap *map[string]model.HandlerFuncInOut, buildRoutePathMap *map[string]map[string]map[string][]string) {

	routeCodeMd5New := encrypt.Md5V(json.ToJsonIgnoreError(buildRouteMethodMap) + json.ToJsonIgnoreError(buildRoutePathMap))

	logger.Logger().Info(" Gen route code md5: " + routeCodeMd5New + "; old md5: " + routeCodeMd5)
	if routeCodeMd5 == routeCodeMd5New {
		initRouteInfo()
		logger.Logger().Info(" Gen route code  md5 same ignore generatio. " + routeCodeMd5)
		return
	}

	path := file.GetCurrentPath()

	filePath1 := filepath.Join(path, genFileName1)
	filePath2 := filepath.Join(path, genFileName2)
	suffix := "1"

	if file.ExistFile(filePath1) {
		os.Remove(filePath1)
		handlerFuncMap1 = nil
		routeUrlMethodMap1 = nil
		suffix = "2"
	}
	if file.ExistFile(filePath2) {
		os.Remove(filePath2)
		handlerFuncMap2 = nil
		routeUrlMethodMap2 = nil
		suffix = "1"
	}

	pathMapCode := genPathMapCode(buildRoutePathMap, suffix)
	funcListCode := genHandlerFuncMapCode(buildRouteMethodMap, suffix)

	md5Code := blankStr + "routeCodeMd5 = \"" + routeCodeMd5New + "\" " + consts.LineSep

	content := strings.ReplaceAll(globalRouterTemplate, "{{.routersDir}}", "route")
	content = strings.ReplaceAll(content, "{{.appModuleName}}", moduleName)
	content = strings.ReplaceAll(content, "{{.md5CodeInfo}}", md5Code)
	content = strings.ReplaceAll(content, "{{.routeCodeInfo}}", blankStr+""+consts.LineSep+pathMapCode+consts.LineSep+consts.LineSep+funcListCode)

	if suffix == "1" {
		file.WriteFile(filePath1, content)
	} else {
		file.WriteFile(filePath2, content)
	}

	time.Sleep(3 * time.Second)

	initRouteInfo()

	fmt.Println("Routing code regenerated, please restart ......" + routeCodeMd5)

	//os.Exit(0)
}

// genPathMapCode 构造url path 的map代码
func genPathMapCode(buildRoutePathMap *map[string]map[string]map[string][]string, suffix string) string {

	routeUrlMethodMapLines := make([]string, 0)

	routeUrlMethodMapLines = append(routeUrlMethodMapLines, blankStr+"routeUrlMethodMapTmp := make(map[string]map[string]map[string][]string)", consts.EmptyStr)
	for preUrl, preRouteMap := range *buildRoutePathMap {
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
	routeUrlMethodMapLines = append(routeUrlMethodMapLines, blankStr+"routeUrlMethodMap"+suffix+" = routeUrlMethodMapTmp", consts.EmptyStr)
	return strings.Join(routeUrlMethodMapLines, consts.LineSep)
}

// genHandlerFuncMapCode 构造handler方法map的代码
func genHandlerFuncMapCode(buildRouteMethodMap *map[string]model.HandlerFuncInOut, suffix string) string {

	handlerFuncsCode := make([]string, 0)
	handlerFuncsCode = append(handlerFuncsCode, blankStr+"handlerFuncMapTmp := make(map[string]model.HandlerFuncInOut)", consts.EmptyStr)

	controllerNameMap := make(map[string]string)

	for k, handlerFunc := range *buildRouteMethodMap {

		controllerInstanceName := strs.LcFirst(handlerFunc.ControllerName)
		controllerNameMap[controllerInstanceName] = handlerFunc.ControllerName

		inParamsCode := make([]string, 0)

		for _, inParam := range *handlerFunc.Ins {
			code := strings.ReplaceAll(handlerFuncInParamTemplate, "{{.inParamName}}", inParam.Name)
			code = strings.ReplaceAll(code, "{{.assignType}}", strconv.Itoa(int(inParam.AssignType)))
			code = strings.ReplaceAll(code, "{{.isPointer}}", strconv.FormatBool(inParam.IsPointer))
			code = strings.ReplaceAll(code, "{{.defaultVal}}", strings.ReplaceAll(inParam.DefaultVal, "\"", "\\\""))
			code = strings.ReplaceAll(code, "{{.isNeed}}", strconv.FormatBool(inParam.IsNeed))
			inParamsCode = append(inParamsCode, code)
		}

		onParamsCode := make([]string, 0)

		outCount := len(*handlerFunc.Outs)
		for i := 0; i < outCount; i++ {
			onParamsCode = append(onParamsCode, handlerFuncOutParamTemplate)
		}

		code := strings.ReplaceAll(handlerFuncMapTemplate, "{{.controllerName}}", handlerFunc.ControllerName)
		code = strings.ReplaceAll(code, "{{.methodName}}", handlerFunc.Name)
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

	genCode := strings.Join(controllerInstanceCode, consts.LineSep) + consts.LineSep + strings.Join(handlerFuncsCode, consts.LineSep) + consts.LineSep + blankStr + "handlerFuncMap" + suffix + " = handlerFuncMapTmp"

	return genCode
}
