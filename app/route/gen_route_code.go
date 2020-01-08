package route

import (
	"fmt"
	"github.com/treeyh/soc-go-boot/app/model"
	"strconv"
	"strings"
)

const (
	globalRouterTemplate = `package {{.routersDir}}

import (
	"{{.appModuleName}}/app/controller"
	"github.com/treeyh/soc-go-boot/app/model"
)

func init() {
{{.globalInfo}}
}
`

	BlankStr = "	"
)

// genRouterCode 构建路由初始化代码
func genRouterCode(moduleName string, buildRouteMethodMap *map[string]model.HandlerFuncInOut, buildRoutePathMap *map[string]map[string]map[string][]string) {
	routeUrlMethodMapLines := make([]string, 0)

	routeUrlMethodMapLines = append(routeUrlMethodMapLines, BlankStr+"routeUrlMethodMapTmp := make(map[string]map[string]map[string][]string)", "")

	for preUrl, preRouteMap := range *buildRoutePathMap {
		routeUrlMethodMapLines = append(routeUrlMethodMapLines, BlankStr+"routeUrlMethodMapTmp[\""+preUrl+"\"] = make(map[string]map[string][]string)")

		for route, methodMap := range preRouteMap {
			routeUrlMethodMapLines = append(routeUrlMethodMapLines, BlankStr+"routeUrlMethodMapTmp[\""+preUrl+"\"][\""+route+"\"] = make(map[string][]string)")
			for method, handlerNames := range methodMap {
				routeUrlMethodMapLines = append(routeUrlMethodMapLines, BlankStr+"routeUrlMethodMapTmp[\""+preUrl+"\"][\""+route+"\"][\""+method+"\"] = make([]string, 0)")
				for i, handlerName := range handlerNames {
					routeUrlMethodMapLines = append(routeUrlMethodMapLines, BlankStr+"routeUrlMethodMapTmp[\""+preUrl+"\"][\""+route+"\"][\""+method+"\"]["+strconv.Itoa(i)+"] = \""+handlerName+"\"")
				}
				routeUrlMethodMapLines = append(routeUrlMethodMapLines, "")
			}
			routeUrlMethodMapLines = append(routeUrlMethodMapLines, "")
		}
		routeUrlMethodMapLines = append(routeUrlMethodMapLines, "")
	}

	fmt.Println(strings.Join(routeUrlMethodMapLines, "\n"))

}
