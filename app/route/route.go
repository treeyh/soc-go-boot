package route

import (
	"bufio"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/treeyh/soc-go-boot/app/config"
	"github.com/treeyh/soc-go-boot/app/controller"
	"github.com/treeyh/soc-go-boot/app/model"
	"github.com/treeyh/soc-go-boot/app/model/resp"
	"github.com/treeyh/soc-go-common/core/consts"
	"github.com/treeyh/soc-go-common/core/errors"
	"github.com/treeyh/soc-go-common/core/logger"
	"github.com/treeyh/soc-go-common/core/utils/file"
	"github.com/treeyh/soc-go-common/core/utils/json"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"path/filepath"
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

// buildRouteMap 本地环境根据Controller注释构建RouteMap
func buildRouteMap(contrs ...controller.IController) {
	if consts.GetCurrentEnv() != consts.EnvLocal {
		return
	}

	//goModPath := filepath.Join(file.GetCurrentPath(), "..", "..", "go.mod")
	//moduleName := readGoModModule(goModPath)

	buildRouteMap, buildRouteMap2 := buildHandlerFuncMap(contrs...)

	fmt.Println(json.ToJsonIgnoreError(buildRouteMap))
	fmt.Println(json.ToJsonIgnoreError(buildRouteMap2))

	//genRouterCode(pkgRealpath)
	//savetoFile(pkgRealpath)
}

// buildHandlerFuncMap 解析controller 构建 HandlerFucMap, 返回两个对象，一个key是preurl,子key是controllerName.methodName，第二个对象key是PreUrl, 子key是RouteUrl
func buildHandlerFuncMap(contrs ...controller.IController) (*map[string]map[string]model.HandlerFuncInOut, *map[string]map[string][]model.HandlerFuncInOut) {

	// 构建需初始化controller列表，没有则全部初始化
	controllerNames := *buildRouteControllerMap(contrs...)

	// 获取controller路径
	controllerPath := filepath.Join(file.GetCurrentPath(), "..", "controller")

	fileSet := token.NewFileSet()
	astPkgs, err := parser.ParseDir(fileSet, controllerPath, func(info os.FileInfo) bool {
		name := info.Name()
		return !info.IsDir() && !strings.HasPrefix(name, ".") && strings.HasSuffix(name, ".go")
	}, parser.ParseComments)

	if err != nil {
		panic("build handler by controller fail. load " + controllerPath + " error: " + err.Error())
	}

	buildRouteMap := make(map[string]map[string]model.HandlerFuncInOut)
	buildRouteMap2 := make(map[string]map[string][]model.HandlerFuncInOut)
	for _, pkg := range astPkgs {
		for _, fl := range pkg.Files {
			for _, d := range fl.Decls {
				switch specDecl := d.(type) {
				case *ast.FuncDecl:
					if specDecl.Recv != nil && len(specDecl.Type.Params.List) > 0 {
						exp, ok := specDecl.Recv.List[0].Type.(*ast.StarExpr)
						if !ok {
							continue
						}
						controllerName := fmt.Sprint(exp.X)

						// 不能空默认加载所有路由，因为需要根据Controller获取PreUrl
						//if len(contrs) > 0 {
						//
						//}
						preUrl := ""
						if preUrl, ok = controllerNames[controllerName]; !ok {
							// 不需要初始化
							continue
						}
						handlerFunc := *(parseHandlerFunc(controllerName, preUrl, specDecl))

						if _, ok := buildRouteMap[preUrl]; !ok {
							buildRouteMap[preUrl] = make(map[string]model.HandlerFuncInOut)
							buildRouteMap2[preUrl] = make(map[string][]model.HandlerFuncInOut)
						}
						buildRouteMap[preUrl][controllerName+"."+handlerFunc.Name] = handlerFunc
						for _, v := range *handlerFunc.RouteMethods {
							if _, ok := buildRouteMap2[preUrl][v.Route]; !ok {
								buildRouteMap2[preUrl][v.Route] = make([]model.HandlerFuncInOut, 0)
							}

							buildRouteMap2[preUrl][v.Route] = append(buildRouteMap2[preUrl][v.Route], handlerFunc)
						}
					}
				}
			}
		}
	}

	return &buildRouteMap, &buildRouteMap2
}

// parseHandlerFunc 构建HandlerFunc
func parseHandlerFunc(controllerName, preUrl string, specDecl *ast.FuncDecl) *model.HandlerFuncInOut {

	handlerRoute := model.HandlerFuncRoute{}
	handlerRoute.PreUrl = preUrl
	routeMethods := make([]model.RouteMethod, 0)
	for _, v := range specDecl.Doc.List {
		t := strings.TrimSpace(strings.TrimLeft(v.Text, "//"))
		if !strings.HasPrefix(t, "@router") {
			continue
		}
		matches := routeRegex.FindStringSubmatch(t)
		routeMethod := model.RouteMethod{}
		if len(matches) != 3 {
			panic(" @route format does not to the rules. " + v.Text)
		}

		routeMethod.Route = matches[1]
		routeMethod.PreUrl = preUrl
		methods := strings.ToUpper(matches[2])
		if matches[2] == "" {
			routeMethod.Methods = []string{"GET"}
		} else {
			routeMethod.Methods = strings.Split(methods, ",")
			for _, httpMethod := range routeMethod.Methods {
				if !strings.Contains(httpMethods, ","+httpMethod+",") {
					panic(" @route http method format does not to the rules. " + httpMethod)
				}
			}
		}
		routeMethods = append(routeMethods, routeMethod)
	}

	if len(routeMethods) == 0 {
		//无需设置路由
		return nil
	}

	handlerFunc := model.HandlerFuncInOut{}
	handlerFunc.Name = specDecl.Name.Name
	handlerFunc.ControllerName = controllerName
	handlerFunc.RouteMethods = &routeMethods
	ins := make([]model.InParamsType, 0)
	for _, param := range specDecl.Type.Params.List {
		for _, pn := range param.Names {
			fmt.Println("param:" + pn.Name)
			ins = append(ins, model.InParamsType{Name: pn.Name})
		}
	}
	handlerFunc.Ins = &ins
	outs := make([]model.ParamsType, 0)
	for _, param := range specDecl.Type.Results.List {
		for _, pn := range param.Names {
			fmt.Println("return:" + pn.Name)
		}
		outs = append(outs, model.ParamsType{})
	}
	handlerFunc.Outs = &outs
	return &handlerFunc
}

// buildRouteControllerMap 获取需要路由的Constroller map
func buildRouteControllerMap(contrs ...controller.IController) *map[string]string {
	controllerNames := make(map[string]string)
	for _, controller := range contrs {
		reflectVal := reflect.ValueOf(controller)
		contr := reflect.Indirect(reflectVal).Type()
		if contr.Kind() != reflect.Struct || !strings.HasSuffix(contr.Name(), "Controller") {
			panic("build handler by controller fail. " + contr.String() + " not struct or struct name not has 'Controller' suffix.")
		}
		controllerNames[contr.Name()] = controller.PreUrl()
	}
	return &controllerNames
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
			panic(methodName + " The first parameter needs to be *gin.Context, the return value is only one and must be *resp.RespResult. ")
		}
		if targetType.NumOut() != 1 {
			panic(methodName + " The first parameter needs to be *gin.Context, the return value is only one and must be *resp.RespResult. ")
		}
		urlPaths := strings.Split(suffixUrl, "/")
		// 构建输入参数列表
		for i, inParam := range *handlerFunc.Ins {
			elem := targetType.In(i)
			isPtr := elem.Kind() == reflect.Ptr
			if isPtr {
				inParam.Kind = elem.Elem().Kind()
				inParam.Type = elem.Elem()
			} else {
				inParam.Kind = elem.Kind()
				inParam.Type = elem
			}
			inParam.IsPointer = isPtr
			if checkParamExistUrl(&urlPaths, inParam.Name) {
				inParam.AssignType = model.UrlAssign
			}
		}

		// 构建输出参数
		for i, outParam := range *handlerFunc.Outs {
			elem := targetType.Out(i)
			isPtr := elem.Kind() == reflect.Ptr
			fmt.Println("aaaaa:" + elem.Kind().String())
			if isPtr {
				outParam.Kind = elem.Elem().Kind()
				outParam.Type = elem.Elem()
			} else {
				outParam.Kind = elem.Kind()
				outParam.Type = elem.Elem()
			}

			fmt.Println(outParam.Kind.String())
			fmt.Println(outParam.Type.String())
		}
		fmt.Println(json.ToJson(handlerFunc))
		handlers = append(handlers, httpHandler(&handlerFunc))
	}

	//
	//if paramTypes[0].Type.String() != "gin.Context" {
	//	logger.Logger().Fatal(" buildHandler " + key + " first params type need gin.Context ")
	//}

	//handler := restHandler(targetFunc)
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

// readGoModModule 读取go.mod的模块
func readGoModModule(goModPath string) string {
	file, err := os.OpenFile(goModPath, os.O_RDONLY, 0666)
	if err != nil {
		panic("read go mod file fail! " + err.Error())
	}
	defer func() {
		if file != nil {
			file.Close()
		}
	}()

	buf := bufio.NewReader(file)
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("File read over!")
				break
			} else {
				panic("read go mod file fail! " + err.Error())
			}
		}

		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "module ") {
			continue
		}
		line = strings.Replace(line, "module ", "", 1)
		return strings.TrimSpace(line)
	}
	return ""
}

func httpHandler(handlerFunc *model.HandlerFuncInOut) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		ginContext := resp.GinContext{Ctx: ctx}
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
		ginContext.JsonRespResult(&respObj)
	}
}

func injectFunc(ctx *resp.GinContext, handlerFunc *model.HandlerFuncInOut) (*[]reflect.Value, errors.AppError) {

	//inputValues := make([]reflect.Value, numIn)

	return nil, nil
}

//func InjectFunc(targetFunc interface{}, reqInfo RequestInfo) ([]reflect.Value, error) {
//	targetType := reflect.TypeOf(targetFunc)
//	if reflect.Func != targetType.Kind() {
//		return nil, errors.New("target is not func")
//	}
//	numIn := targetType.NumIn()
//	inputValues := make([]reflect.Value, numIn)
//	if numIn > 0 {
//		for i := 0; i < numIn; i++ {
//			elem := targetType.In(i)
//			isPtr := false
//			//if elem.Kind() == reflect.Ptr {
//			//	elem = elem.Elem()
//			//	isPtr = true
//			//}
//			judgeIsPtr(&isPtr, &elem)
//			//if elem.String() == "context.Context" {
//			//	if isPtr {
//			//		inputValues[i] = reflect.ValueOf(&reqInfo.Ctx)
//			//	} else {
//			//		inputValues[i] = reflect.ValueOf(reqInfo.Ctx)
//			//	}
//			//	continue
//			//}
//			if assemblyInputValues(&inputValues, elem, isPtr, reqInfo, i) {
//				continue
//			}
//			if elem.Kind() == reflect.Struct {
//				value, err := ParseValue(elem, isPtr, reqInfo)
//				if err != nil {
//					log.Error(fmt.Sprintf("%#v", err))
//					return nil, err
//				}
//				inputValues[i] = *value
//			}
//		}
//	}
//	return reflect.ValueOf(targetFunc).Call(inputValues), nil
//}
//
//
//func assemblyInputValues(inputValues *[]reflect.Value, elem reflect.Type, isPtr bool, reqInfo RequestInfo, i int) bool {
//	if elem.String() == "context.Context" {
//		if isPtr {
//			(*inputValues)[i] = reflect.ValueOf(&reqInfo.Ctx)
//		} else {
//			(*inputValues)[i] = reflect.ValueOf(reqInfo.Ctx)
//		}
//		return true
//	}
//	return false
//}
//
//func ParseValue(elem reflect.Type, isPtr bool, reqInfo RequestInfo) (*reflect.Value, error) {
//	reqObj := reflect.New(elem).Elem()
//	for i := 0; i < elem.NumField(); i++ {
//		field := elem.Field(i)
//		fieldType := field.Type
//		isPtr := false
//		if fieldType.Kind() == reflect.Ptr {
//			fieldType = fieldType.Elem()
//			isPtr = true
//		}
//		var target *reflect.Value = nil
//		var err error = nil
//
//		if converter, ok := BasicTypeConverter[fieldType.String()]; ok {
//			target, err = ParseQuery(field.Name, isPtr, converter, reqInfo)
//		} else if isBodyFlag(fieldType.Kind()) {
//			target, err = ParseBody(fieldType, isPtr, reqInfo)
//		}
//		if err != nil {
//			log.Error(fmt.Sprintf("%#v", err))
//			return nil, err
//		}
//
//		if target != nil {
//			reqObj.FieldByName(field.Name).Set(*target)
//		}
//	}
//	if isPtr {
//		reqObj = reqObj.Addr()
//	}
//	return &reqObj, nil
//}
//
//func ParseBody(elem reflect.Type, isPtr bool, reqInfo RequestInfo) (*reflect.Value, error) {
//	body := reqInfo.Body
//	if body == "" && isPtr {
//		return nil, nil
//	}
//	newStrut := reflect.New(elem)
//	targetInterface := newStrut.Interface()
//	err := json.FromJson(body, &targetInterface)
//	if err != nil {
//		log.Error(err)
//		return nil, err
//	}
//	if !isPtr {
//		newStrut = newStrut.Elem()
//	}
//	return &newStrut, nil
//}
//
//func ParseQuery(fieldName string, isPtr bool, converter func(v string) (interface{}, error), reqInfo RequestInfo) (*reflect.Value, error) {
//	paramValue := reqInfo.Parameters[str.LcFirst(fieldName)]
//	//指针直接返回空
//	if paramValue == "" && isPtr {
//		return nil, nil
//	}
//	v, err := converter(paramValue)
//	if err != nil {
//		log.Error(fmt.Sprintf("%#v", err))
//		return nil, err
//	}
//	va := reflect.ValueOf(v)
//	if !isPtr {
//		va = va.Elem()
//	}
//	return &va, nil
//}
