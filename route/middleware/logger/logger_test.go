package logger

import (
	"context"
	"fmt"
	"testing"
)

func TestLang(t *testing.T) {
	lang := "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6"

	fmt.Println(formatRequestLang(context.Background(), lang))
}
