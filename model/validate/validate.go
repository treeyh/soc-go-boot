package validate

import (
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/treeyh/soc-go-common/core/logger"
	"strings"
)

var (
	vali              *validator.Validate
	trans             ut.Translator
	customValidateMap = map[string]CustomCheckInfo{
		"checkEmailHost": {
			fun:     checkEmailHost,
			message: "Email需要以@email.com结尾",
		},
	}
)

type CustomCheckInfo struct {
	fun     func(fl validator.FieldLevel) bool
	message string
}

func init() {
	//中文翻译器
	zh_ch := zh.New()
	uni := ut.New(zh_ch)
	trans, _ = uni.GetTranslator("zh")
	//验证器
	vali = validator.New()

	//添加自定义参数校验
	for k, v := range customValidateMap {
		err := vali.RegisterValidation(k, v.fun)
		if err != nil {
			logger.Logger().Error(err)
		}
	}
	//验证器注册翻译器
	zh_translations.RegisterDefaultTranslations(vali, trans)
}

func checkEmailHost(fl validator.FieldLevel) bool {
	email := fl.Field().String()
	if strings.Contains(email, "@email.com") {
		return true
	}
	return false
}

// ValidateObject 校验对象参数
func ValidateObject(obj interface{}) string {

	errs := make([]string, 0)
	err := vali.Struct(obj)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			//翻译错误信息
			if v, ok := customValidateMap[err.Tag()]; ok {
				errs = append(errs, v.message)
			} else {
				errs = append(errs, err.Translate(trans))
			}
		}
		msg := strings.Join(errs, ",")
		return msg
	}

	return ""
}
