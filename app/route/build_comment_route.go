package route

import (
	"bufio"
	"fmt"
	"github.com/treeyh/soc-go-boot/app/controller"
	"github.com/treeyh/soc-go-boot/app/model"
	"github.com/treeyh/soc-go-common/core/consts"
	"github.com/treeyh/soc-go-common/core/utils/file"
	"github.com/treeyh/soc-go-common/core/utils/json"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

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
			ins = append(ins, model.InParamsType{Name: pn.Name})
		}
	}
	handlerFunc.Ins = &ins
	outs := make([]model.ParamsType, 0)
	for _, param := range specDecl.Type.Results.List {
		if len(param.Names) > 0 {
			for _, pn := range param.Names {
				fmt.Sprint("return:" + pn.Name)
				outs = append(outs, model.ParamsType{})
			}
		} else {
			outs = append(outs, model.ParamsType{})
		}
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
