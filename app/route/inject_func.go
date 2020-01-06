package route

import (
	"github.com/treeyh/soc-go-boot/app/model"
	"github.com/treeyh/soc-go-boot/app/model/req"
	"github.com/treeyh/soc-go-common/core/errors"
	"reflect"
)

func injectFunc(ctx *req.GinContext, handlerFunc *model.HandlerFuncInOut) (*[]reflect.Value, errors.AppError) {

	//inputValues := make([]reflect.Value, handlerFunc.InCount)
	//if handlerFunc.InCount > 0 {
	//	for i := 0; i < handlerFunc.InCount; i++ {
	//
	//
	//		if assemblyInputValues(&inputValues, elem, isPtr, reqInfo, i) {
	//			continue
	//		}
	//		if elem.Kind() == reflect.Struct {
	//			value, err := ParseValue(elem, isPtr, reqInfo)
	//			if err != nil {
	//				log.Error(fmt.Sprintf("%#v", err))
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
