package route

import (
	"bufio"
	"fmt"
	"github.com/treeyh/soc-go-boot/controller"
	"github.com/treeyh/soc-go-boot/model"
	"github.com/treeyh/soc-go-common/core/logger"
	"github.com/treeyh/soc-go-common/core/utils/file"
	"github.com/treeyh/soc-go-common/core/utils/json"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

const httpMethods = ",GET,POST,DELETE,PATCH,PUT,OPTIONS,HEAD,*,"

var (
	// controllerStatusTmpFileName controller记录文件，判断是否要重新生成路由
	controllerStatusTmpFileName = ".last_controller_status.tmp"
	log                         = logger.Logger()

	routeRegex = regexp.MustCompile(`@Router\s+(\S+)(?:\s+\[(\S+)\])?`)
)

// buildRouteMap 本地环境根据Controller注释构建RouteMap
func buildRouteMap(controllerStatusPath, controllerPath, goModFilePath, genPath string, contrs ...controller.IController) {

	if !checkControllerStatus(controllerStatusPath, controllerPath) {
		return
	}

	buildRouteMap, buildRouteMap2 := buildHandlerFuncMap(controllerPath, contrs...)

	log.Info("buildRouteMap:" + json.ToJsonIgnoreError(buildRouteMap))
	log.Info("buildRouteMap2:" + json.ToJsonIgnoreError(buildRouteMap2))

	moduleName := readGoModModule(goModFilePath)
	genRouterCode(genPath, moduleName, buildRouteMap, buildRouteMap2)
}

// buildHandlerFuncMap 解析controller 构建 HandlerFucMap, 返回两个对象，一个key是controllerName.methodName,value是controller.method，第二个对象key是PreUrl, 子key是RouteUrl, 再子key是httpMethod, value 是 controllerName.methodName 列表
func buildHandlerFuncMap(controllerPath string, contrs ...controller.IController) (map[string]model.HandlerFuncInOut, map[string]map[string]map[string][]string) {

	// 构建需初始化controller列表，没有则全部初始化
	controllerNames := buildRouteControllerMap(contrs...)

	fileSet := token.NewFileSet()
	astPkgs, err := parser.ParseDir(fileSet, controllerPath, func(info os.FileInfo) bool {
		name := info.Name()
		return !info.IsDir() && !strings.HasPrefix(name, ".") && strings.HasSuffix(name, ".go")
	}, parser.ParseComments)

	if err != nil {
		panic("build handler by controller fail. load " + controllerPath + " error: " + err.Error())
	}

	buildRouteMap := make(map[string]model.HandlerFuncInOut)
	buildRouteMap2 := make(map[string]map[string]map[string][]string)
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
						funcPtr := parseHandlerFunc(controllerName, preUrl, specDecl)
						if funcPtr == nil {
							continue
						}
						handlerFunc := *(funcPtr)

						fullName := controllerName + "." + handlerFunc.Name
						if _, ok := buildRouteMap[fullName]; !ok {
							buildRouteMap[fullName] = handlerFunc
						}

						if _, ok := buildRouteMap2[preUrl]; !ok {
							buildRouteMap2[preUrl] = make(map[string]map[string][]string)
						}

						for _, v := range handlerFunc.RouteMethods {
							if _, ok := buildRouteMap2[preUrl][v.Route]; !ok {
								buildRouteMap2[preUrl][v.Route] = make(map[string][]string, 0)
							}
							for _, method := range v.Methods {
								if _, ok := buildRouteMap2[preUrl][v.Route][method]; !ok {
									buildRouteMap2[preUrl][v.Route][method] = make([]string, 0)
								}
								buildRouteMap2[preUrl][v.Route][method] = append(buildRouteMap2[preUrl][v.Route][method], fullName)
							}

						}
					}
				}
			}
		}
	}

	return buildRouteMap, buildRouteMap2
}

// parseHandlerFunc 构建HandlerFunc
func parseHandlerFunc(controllerName, preUrl string, specDecl *ast.FuncDecl) *model.HandlerFuncInOut {

	handlerRoute := model.HandlerFuncRoute{}
	handlerRoute.PreUrl = preUrl
	routeMethods := make([]model.RouteMethod, 0)
	paramMaps := make(map[string]model.InParamsType, 0)
	if specDecl.Doc == nil {
		return nil
	}

	for _, v := range specDecl.Doc.List {
		t := strings.TrimSpace(strings.TrimLeft(v.Text, "//"))
		if !strings.HasPrefix(t, "@Router") && !strings.HasPrefix(t, "@Param") {
			continue
		}

		if strings.HasPrefix(t, "@Router") {
			matches := routeRegex.FindStringSubmatch(t)
			routeMethod := model.RouteMethod{}
			if len(matches) != 3 {
				panic(" @Route format does not to the rules. " + v.Text)
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
						panic(" @Route http method format does not to the rules. " + httpMethod)
					}
				}
			}
			routeMethods = append(routeMethods, routeMethod)
		}
		if strings.HasPrefix(t, "@Param") {
			pv := getParams(strings.TrimSpace(strings.TrimLeft(t, "@Param")))
			if len(pv) < 4 {
				logger.Logger().Error("Invalid @Param format. Needs at least 4 parameters")
			}
			param := model.InParamsType{}
			param.Name = pv[0]
			param.AssignType = getParamAssignType(pv[1])
			param.DefaultVal = pv[2]

			switch len(pv) {
			case 4:
				param.IsNeed, _ = strconv.ParseBool(pv[2])
				param.DefaultVal = ""
			case 5:
				param.IsNeed, _ = strconv.ParseBool(pv[2])
				param.DefaultVal = pv[3]
			}
			paramMaps[param.Name] = param
		}
	}

	if len(routeMethods) == 0 {
		//无需设置路由
		return nil
	}

	handlerFunc := model.HandlerFuncInOut{}
	handlerFunc.Name = specDecl.Name.Name
	handlerFunc.ControllerName = controllerName
	handlerFunc.RouteMethods = routeMethods
	ins := make([]model.InParamsType, 0)
	for _, param := range specDecl.Type.Params.List {
		for _, pn := range param.Names {
			if v, ok := paramMaps[pn.Name]; ok {
				if v.AssignType == model.HeaderAssign {
					v.Name = strings.ReplaceAll(v.Name, "_", "-")
				}
				ins = append(ins, v)
				continue
			}
			ins = append(ins, model.InParamsType{Name: pn.Name})
		}
	}
	handlerFunc.Ins = ins
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
	handlerFunc.Outs = outs
	return &handlerFunc
}

// getParamAssignType 获取参数注解
func getParamAssignType(val string) model.HttpParamsAssignType {
	switch strings.ToLower(val) {
	case "formdata":
		return model.PostFormAssign
	case "query":
		return model.QueryAssign
	case "path":
		return model.PathAssign
	case "body":
		return model.BodyAssign
	case "header":
		return model.HeaderAssign
	}
	return model.UnAssign
}

// getParams 解析@Params备注
func getParams(str string) []string {
	var s []rune
	var j int
	var start bool
	var r []string
	var quoted int8
	for _, c := range str {
		if unicode.IsSpace(c) && quoted == 0 {
			if !start {
				continue
			} else {
				start = false
				j++
				r = append(r, string(s))
				s = make([]rune, 0)
				continue
			}
		}

		start = true
		if c == '"' {
			quoted ^= 1
			continue
		}
		s = append(s, c)
	}
	if len(s) > 0 {
		r = append(r, string(s))
	}
	return r
}

// buildRouteControllerMap 获取需要路由的Constroller map
func buildRouteControllerMap(contrs ...controller.IController) map[string]string {
	controllerNames := make(map[string]string)
	for _, controller := range contrs {
		reflectVal := reflect.ValueOf(controller)
		contr := reflect.Indirect(reflectVal).Type()
		if contr.Kind() != reflect.Struct || !strings.HasSuffix(contr.Name(), "Controller") {
			panic("build handler by controller fail. " + contr.String() + " not struct or struct name not has 'Controller' suffix.")
		}
		controllerNames[contr.Name()] = controller.PreUrl()
	}
	return controllerNames
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

// checkControllerStatus 获取controller文件状态
func checkControllerStatus(controllerStatusPath, controllerPath string) bool {
	t := false

	controllerStatusFilePath := filepath.Join(controllerStatusPath, controllerStatusTmpFileName)
	if !file.ExistFile(controllerStatusFilePath) {
		t = true
	}

	tmpContent, err := ioutil.ReadFile(controllerStatusFilePath)
	if err != nil {
		t = true
	}
	tmpJson := string(tmpContent)
	tmpMap := make(map[string]int64)
	json.FromJson(tmpJson, &tmpMap)

	files, err := file.GetDirSon(controllerPath)
	if err != nil {
		return true
	}
	for _, v := range files {
		if v.IsDir() {
			continue
		}

		if fi, ok := tmpMap[v.Name()]; ok {
			if fi != v.ModTime().Unix() {
				if strings.HasSuffix(v.Name(), "_controller.go") {
					t = true
					tmpMap[v.Name()] = v.ModTime().Unix()
				}
			}
			continue
		}

		if strings.HasSuffix(v.Name(), "controller.go") {
			t = true
			tmpMap[v.Name()] = v.ModTime().Unix()
		}
	}

	if t {
		file.WriteFile(controllerStatusFilePath, json.ToJsonIgnoreError(tmpMap))
	}
	return t
}
