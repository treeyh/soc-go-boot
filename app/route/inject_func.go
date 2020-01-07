package route

import (
	"github.com/gin-gonic/gin"
	"github.com/treeyh/soc-go-boot/app/model"
	"github.com/treeyh/soc-go-boot/app/model/req"
	"github.com/treeyh/soc-go-common/core/errors"
	"reflect"
	"strconv"
)

var (
	ParamConverter = map[reflect.Kind]func(v string) (interface{}, error){}
)

func init() {
	intFunc := func(bitSize int) func(v string) (interface{}, error) {
		return func(v string) (interface{}, error) {
			obj, err := strconv.ParseInt(v, 10, bitSize)
			if err != nil {
				return 0, err
			}
			switch bitSize {
			case 0:
				obj := int(obj)
				return &obj, nil
			case 8:
				obj := int8(obj)
				return &obj, nil
			case 16:
				obj := int16(obj)
				return &obj, nil
			case 32:
				obj := int32(obj)
				return &obj, nil
			}
			return &obj, nil
		}
	}
	floatFunc := func(bitSize int) func(v string) (interface{}, error) {
		return func(v string) (interface{}, error) {
			obj, err := strconv.ParseFloat(v, bitSize)
			if err != nil {
				return 0, err
			}
			switch bitSize {
			case 32:
				obj := float32(obj)
				return &obj, nil
			}
			return &obj, nil
		}
	}
	uintFunc := func(bitSize int) func(v string) (interface{}, error) {
		return func(v string) (interface{}, error) {
			obj, err := strconv.ParseUint(v, 10, bitSize)
			if err != nil {
				return 0, err
			}
			switch bitSize {
			case 0:
				obj := uint(obj)
				return &obj, nil
			case 8:
				obj := uint8(obj)
				return &obj, nil
			case 16:
				obj := uint16(obj)
				return &obj, nil
			case 32:
				obj := uint32(obj)
				return &obj, nil
			}
			return &obj, nil
		}
	}
	ParamConverter[reflect.Bool] = func(v string) (interface{}, error) {
		b, e := strconv.ParseBool(v)
		return &b, e
	}

	ParamConverter[reflect.Int] = intFunc(0)
	ParamConverter[reflect.Int8] = intFunc(8)
	ParamConverter[reflect.Int16] = intFunc(16)
	ParamConverter[reflect.Int32] = intFunc(32)
	ParamConverter[reflect.Int64] = intFunc(64)
	ParamConverter[reflect.Uint] = uintFunc(0)
	ParamConverter[reflect.Uint8] = uintFunc(8)
	ParamConverter[reflect.Uint16] = uintFunc(16)
	ParamConverter[reflect.Uint32] = uintFunc(32)
	ParamConverter[reflect.Uint64] = uintFunc(64)
	ParamConverter[reflect.Float32] = floatFunc(32)
	ParamConverter[reflect.Float64] = floatFunc(64)
	ParamConverter[reflect.String] = func(v string) (interface{}, error) {
		return &v, nil
	}
	//ParamConverter["date.Time"] = func(v string) (interface{}, error) {
	//	t, e := time.Parse(consts.AppTimeFormat, v)
	//	return &t, e
	//}
	//ParamConverter["types.Time"] = func(v string) (interface{}, error) {
	//	t, err := time.Parse(consts.AppTimeFormat, v)
	//	tt := types.Time(t)
	//	return &tt, err
	//}
}

func injectFunc(ctx *req.GinContext, handlerFunc *model.HandlerFuncInOut) (*[]reflect.Value, errors.AppError) {
	//inputValues := make([]reflect.Value , 0, handlerFunc.InCount)
	//if handlerFunc.InCount > 0 {
	//	for i, inParam := range *handlerFunc.Ins {
	//		if i == 0 {
	//			inputValues[i] = reflect.ValueOf(ctx)
	//			continue
	//		}
	//
	//		if v, ok := ParamConverter[inParam.Kind]; ok {
	//			value, err := v()
	//			if err != nil {
	//				logger.Logger().Error(fmt.Sprintf("%#v", err))
	//				return nil, err
	//			}
	//			inputValues[i] = *value
	//		}
	//	}
	//
	//}
	//
	//return &reflect.ValueOf(handlerFunc.Func).Call(inputValues), nil

	return nil, nil
}

// getParamString 根据获取类型，获取参数string
func getParamString(key, defaultVal string, assignType *model.HttpParamsAssignType, ctx *gin.Context) string {

	//switch *assignType {
	//case model.PathAssign:
	//	return ctx.Param(key)
	//case model.QueryAssign:
	//	return ctx.DefaultQuery(key, defaultVal)
	//case model.HeaderAssign:
	//	return ctx.GetHeader(key)
	//case model.PostFormAssign:
	//	return ctx.DefaultPostForm(key, defaultVal)
	//case model.BodyAssign:
	//	return ctx.ShouldBindJSON()
	//}

	return ""

}

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
