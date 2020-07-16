package route

import (
	"fmt"
	"github.com/treeyh/soc-go-boot/model"
	"regexp"
	"testing"
)

var (
	routeRegex2 = regexp.MustCompile(`@Router\s+(\S+)(?:\s+\[(\S+)\])?(?:\s+(\S+))?(?:\s+(\S+))?`)
)

func TestRouteRegexFindStringSubmatch(t *testing.T) {

	str := " @Router /create    [post,get,delete]  xml  string"

	matches := routeRegex2.FindStringSubmatch(str)
	fmt.Println(matches)
	for k, v := range matches {
		fmt.Println(k, v)
	}

	str = " @Router /create    [post] "

	matches = routeRegex2.FindStringSubmatch(str)

	fmt.Println(matches)
	for k, v := range matches {
		fmt.Println(k, v)
	}

	var test model.RouteReqContentType
	fmt.Println("===")
	fmt.Println(test)

	test = model.ReqContentTypeFile
	fmt.Println("===")
	fmt.Println(test)

	test = ""
	fmt.Println("===")
	fmt.Println(test)

}
