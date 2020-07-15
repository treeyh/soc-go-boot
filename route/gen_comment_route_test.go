package route

import (
	"fmt"
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
}
