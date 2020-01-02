package route

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/treeyh/soc-go-boot/app/config"
	"github.com/treeyh/soc-go-boot/app/controller/user_controller"
	"github.com/treeyh/soc-go-boot/app/model"
	"github.com/treeyh/soc-go-common/core/logger"
	"reflect"
)

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

	userRouter := engine.Group(config.GetSocConfig().App.Server.ContextPath + "/user")
	{
		userRouter.POST("", buildHandler("user_controller.Create", user_controller.Create))
	}
}

// buildHandler 构造 处理handler
func buildHandler(key string, targetFunc interface{}) gin.HandlerFunc {

	// 验证 targetFunc 是否符合规范
	targetType := reflect.TypeOf(targetFunc)
	if reflect.Func != targetType.Kind() {
		logger.Logger().Fatal(" buildHandler " + key + " not func ")
	}
	numIn := targetType.NumIn()
	if numIn < 1 {
		logger.Logger().Fatal(key + " not func ")
	}

	// 构建输入参数列表
	paramTypes := make([]model.ParamsType, 0, numIn)
	for i := 0; i < numIn; i++ {
		elem := targetType.In(i)
		fmt.Println("name:" + elem.Name())
		isPtr := false
		GetObjectTypeIgnorePointer(&isPtr, &elem)
		fmt.Println(isPtr)
		fmt.Println(elem.String())
		fmt.Println(elem.Kind())

		paramTypes = append(paramTypes, model.ParamsType{
			IsPointer: isPtr,
			Type:      elem,
		})
	}

	if paramTypes[0].Type.String() != "gin.Context" {
		logger.Logger().Fatal(" buildHandler " + key + " first params type need gin.Context ")
	}

	//handler := restHandler(targetFunc)
	return func(c *gin.Context) {
		//token := c.GetHeader(consts.APP_HEADER_TOKEN_NAME)
		//handler(c)
	}
}

func GetObjectTypeIgnorePointer(isPtr *bool, elem *reflect.Type) {
	if (*elem).Kind() == reflect.Ptr {
		*elem = (*elem).Elem()
		*isPtr = true
	}
}

//func httpHandler(targetFunc interface{}) gin.HandlerFunc {
//	return func(ctx *gin.Context) {
//
//		var respObj interface{} = nil
//
//		reqInfo, err := NewRequestInfo(ctx)
//		if err != nil {
//			log.Error(err)
//			respObj = vo.VoidErr{Err: vo.NewErr(errs.BuildSystemErrorInfoWithMessage(errs.ServerError, err.Error()))}
//		}
//		startTime := time.Now()
//		respData, err := InjectFunc(targetFunc, *reqInfo)
//		elapsed := time.Since(startTime)
//
//		if err != nil {
//			log.Error(err)
//			respObj = vo.VoidErr{Err: vo.NewErr(errs.BuildSystemErrorInfoWithMessage(errs.ServerError, err.Error()))}
//		} else {
//			if len(respData) > 0 {
//				respObj = respData[0].Interface()
//			}
//		}
//		if respObj != nil {
//			var respBody = ""
//			if s, ok := respObj.(string); ok{
//				respBody = s
//			}else{
//				respBody = json.ToJsonIgnoreError(respObj)
//			}
//
//			respContent(ctx, 200, respBody)
//			log.Infof("请求处理完成，总耗时-> [%dms], url-> [%s], respBody-> [%s]", elapsed/1e6, ctx.Request.URL, respBody)
//		}
//	}
//}

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
//func judgeIsPtr(isPtr *bool, elem *reflect.Type) {
//	if (*elem).Kind() == reflect.Ptr {
//		*elem = (*elem).Elem()
//		*isPtr = true
//	}
//}
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
