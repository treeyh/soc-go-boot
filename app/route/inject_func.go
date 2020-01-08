package route

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	socconsts "github.com/treeyh/soc-go-boot/app/common/consts"
	"github.com/treeyh/soc-go-boot/app/model"
	"github.com/treeyh/soc-go-boot/app/model/req"
	"github.com/treeyh/soc-go-common/core/consts"
	"github.com/treeyh/soc-go-common/core/errors"
	"github.com/treeyh/soc-go-common/core/logger"
	"github.com/treeyh/soc-go-common/core/types"
	"reflect"
	"strconv"
	"time"
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

func injectFunc(ctx *req.GinContext, handlerFunc model.HandlerFuncInOut) ([]reflect.Value, errors.AppError) {
	inputValues := make([]reflect.Value, handlerFunc.InCount)
	if handlerFunc.InCount > 0 {
		for i, inParam := range *handlerFunc.Ins {
			if i == 0 {
				inputValues[i] = reflect.ValueOf(ctx)
				continue
			}
			//fmt.Println("i===" + strconv.Itoa(i))
			//fmt.Println(inParam.Kind.String())
			//fmt.Println(inParam.Type.String())
			//fmt.Println(inputValues[0].Kind().String())
			//fmt.Println(inputValues[0].String())
			if inParam.AssignType == model.UnAssign {
				continue
			}

			if inParam.AssignType == model.BodyAssign {
				val, err := parseBody(ctx.Ctx, &inParam)
				if err != nil {
					logger.Logger().Error(err)
					return nil, err
				}
				//fmt.Println(json.ToJson(val))
				inputValues[i] = *val
				continue
			}

			val := getParamString(ctx.Ctx, &inParam)
			if v, ok := ParamConverter[inParam.Kind]; ok {
				value, err := parseBaseType(val, v, &inParam)
				if err != nil {
					logger.Logger().Error(err)
					return nil, err
				}
				inputValues[i] = *value
				continue
			}

			if inParam.Type.String() == "time.Time" || inParam.Type.String() == "types.Time" {
				value, err := parseTimeType(val, &inParam)
				if err != nil {
					logger.Logger().Error(err)
					return nil, err
				}
				inputValues[i] = *value
				continue
			}
		}
	}

	return reflect.ValueOf(handlerFunc.Func).Call(inputValues), nil
}

// parseBaseType 转换基本类型
func parseBaseType(val string, funcc func(v string) (interface{}, error), inParam *model.InParamsType) (*reflect.Value, errors.AppError) {
	value, err := funcc(val)
	if err != nil {
		logger.Logger().Error(err)
		return nil, errors.NewAppErrorByExistError(socconsts.PARSE_PARAM_ERROR, err, inParam.Name)
	}
	va := reflect.ValueOf(value)
	if !inParam.IsPointer {
		va = va.Elem()
	}
	return &va, nil
}

// parseTimeType 转换时间类型
func parseTimeType(val string, inParam *model.InParamsType) (*reflect.Value, errors.AppError) {
	t, err := time.Parse(consts.AppTimeFormat, val)
	if err != nil {
		logger.Logger().Error(err)
		return nil, errors.NewAppErrorByExistError(socconsts.PARSE_PARAM_ERROR, err, inParam.Name)
	}
	var va reflect.Value
	if inParam.Type.String() == "types.Time" {
		tt := types.Time(t)
		va = reflect.ValueOf(tt)
	} else {
		va = reflect.ValueOf(t)
	}
	if inParam.IsPointer {
		va = va.Addr()
	}
	return &va, nil
}

// getParamString 根据获取类型，获取参数string
func getParamString(ctx *gin.Context, inParam *model.InParamsType) string {

	switch inParam.AssignType {
	case model.PathAssign:
		return ctx.Param(inParam.Name)
	case model.QueryAssign:
		return ctx.DefaultQuery(inParam.Name, inParam.DefaultVal)
	case model.HeaderAssign:
		return ctx.GetHeader(inParam.Name)
	case model.PostFormAssign:
		return ctx.DefaultPostForm(inParam.Name, inParam.DefaultVal)
	}

	return ""
}

func parseBody(ctx *gin.Context, inParam *model.InParamsType) (*reflect.Value, errors.AppError) {
	newStrut := reflect.New(inParam.Type)
	targetInterface := newStrut.Interface()
	err := ctx.ShouldBindBodyWith(&targetInterface, binding.JSON)
	if err != nil {
		logger.Logger().Error(err)
		return nil, errors.NewAppErrorByExistError(socconsts.PARSE_PARAM_ERROR, err, inParam.Name)
	}
	if !inParam.IsPointer {
		newStrut = newStrut.Elem()
	}
	return &newStrut, nil
}

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
