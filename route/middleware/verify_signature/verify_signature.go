package verify_signature

import (
	"github.com/gin-gonic/gin"
	"github.com/treeyh/soc-go-boot/boot_config"
	"github.com/treeyh/soc-go-boot/common/boot_consts"
	"github.com/treeyh/soc-go-boot/common/boot_consts/boot_error_consts"
	"github.com/treeyh/soc-go-boot/controller"
	"github.com/treeyh/soc-go-boot/model"
	"github.com/treeyh/soc-go-common/core/consts"
	"github.com/treeyh/soc-go-common/core/errors"
	"github.com/treeyh/soc-go-common/core/logger"
	"github.com/treeyh/soc-go-common/core/utils/encrypt"
	"github.com/treeyh/soc-go-common/core/utils/slice"
	"github.com/treeyh/soc-go-common/core/utils/times"
	"sort"
	"strings"
	"time"
)

var (
	log = logger.Logger()

	encryptMap = map[string]func(string) string{
		boot_consts.SignPolicySha256: encrypt.SHA256String,
		boot_consts.SignPolicyMd5:    encrypt.Md5V,
	}
)

// StartVerifySignature 签名校验中间件，该中间件需要在logger中间件后初始化，匹配 logger.isNeedBody的body才会支持签名
func StartVerifySignature(getVerifyConfig func(*gin.Context) *boot_config.VerifyConfig) gin.HandlerFunc {
	return func(c *gin.Context) {

		if !boot_config.GetSocConfig().Signature.Enable ||
			slice.ContainString(c.Request.URL.Path, boot_config.GetSocConfig().Signature.IgnoreUrls) {
			c.Next()
			return
		}

		querys := c.Request.URL.Query()
		env := boot_config.GetEnv()
		if boot_config.GetSocConfig().Signature.IgnoreQuery != consts.EmptyStr && env != consts.EnvProd && env != consts.EnvStag {
			// 判断忽略参数是否存在
			if _, ok := querys[boot_config.GetSocConfig().Signature.IgnoreQuery]; ok {
				c.Next()
				return
			}
		}

		if !checkTimestampOverLimit(c) {
			return
		}

		verifyConfig := getVerifyConfig(c)
		if verifyConfig == nil {
			controller.FailJson(c, errors.NewAppError(boot_error_consts.SignKeyNotExist))
			c.Abort()
			return
		}

		// 签名源字符串 格式为：{timestampStr}&{排序后的keys对(除了时间戳和签名kv) key1=value1&key2=value2&key3=value3}[&{body}]&{签名key}
		sourceStr := ""
		clientVersion := c.Request.Header.Get(boot_consts.HeaderClientVersion)

		if clientVersion == "0.9.0" {
			sourceStr = buildSignSourceStr090(c)
		} else if clientVersion == "1.0.0" {
			sourceStr = buildSignSourceStr100(c)
		} else {
			sourceStr = buildSignSourceStr(c)
		}

		reqSignStr := c.Request.Header.Get(boot_consts.HeaderSignKey)
		signPolicy := c.Request.Header.Get(boot_consts.HeaderSignPolicyKey)

		checkFlag := false
		for _, sk := range verifyConfig.SecretKeys {
			sourceStrTemp := sourceStr + "&" + sk
			sign := ""
			if v, ok := encryptMap[signPolicy]; ok {
				sign = v(sourceStrTemp)
			} else {
				sign = encryptMap[boot_consts.SignPolicySha256](sourceStrTemp)
			}
			if sign == reqSignStr {
				checkFlag = true
				break
			}
			log.InfoCtx(c.Request.Context(), "sign policy: "+signPolicy+"; sign:"+sign+"; reqSign:"+reqSignStr+"; source:"+sourceStrTemp)
		}

		if !checkFlag {
			controller.FailJson(c, errors.NewAppError(boot_error_consts.SignAuthFail))
			c.Abort()
			return
		}

		c.Next()
	}
}

// checkTimestampOverLimit 校验时间戳是否在阈值范围内
func checkTimestampOverLimit(c *gin.Context) bool {
	timestampStr := c.Request.Header.Get(boot_consts.HeaderTimestampKey)
	timestamp, err := times.ParseByFormat(consts.AppSystemTimeFormat8, timestampStr)
	if err != nil {
		log.ErrorCtx2(c.Request.Context(), err, "header timestamp param error. timestamp:"+timestampStr)

		controller.FailJson(c, errors.NewAppError(errors.ParamError, boot_consts.HeaderTimestampKey))
		c.Abort()
		return false
	}

	now := time.Now().UTC()
	if (now.Unix()+boot_config.GetSocConfig().Signature.TimeRange) < timestamp.UTC().Unix() ||
		(now.Unix()-boot_config.GetSocConfig().Signature.TimeRange) > timestamp.UTC().Unix() {

		controller.FailJson(c, errors.NewAppError(boot_error_consts.RequestTimestampOverLimit))
		c.Abort()
		return false
	}
	return true
}

// buildSignSourceStr 构造签名源字符串
func buildSignSourceStr(c *gin.Context) string {
	// 收集 query和header参数
	params := make(map[string]string)
	keys := make([]string, 0)
	querys := c.Request.URL.Query()
	for k, _ := range querys {
		params[k] = querys[k][0]
		keys = append(keys, k)
	}
	for k, _ := range c.Request.Header {
		lowerKey := strings.ToLower(k)
		if !slice.ContainString(lowerKey, boot_config.GetSocConfig().Signature.Headers) {
			continue
		}
		params[lowerKey] = c.Request.Header.Get(k)
		keys = append(keys, lowerKey)
	}
	// 排序
	sort.Strings(keys)

	// 签名源字符串 格式为：{timestampStr}&{排序后的keys对(除了时间戳和签名kv) key1=value1&key2=value2&key3=value3}[&{body}]&{签名key}
	sourceStr := c.Request.Header.Get(boot_consts.HeaderTimestampKey)

	for _, v := range keys {
		if v == boot_consts.HeaderSignKey || v == boot_consts.HeaderTimestampKey {
			continue
		}
		sourceStr += "&" + v + "=" + params[v]
	}
	body := model.GetHttpContext(c.Request.Context()).Body
	if body != "" {
		sourceStr += "&" + body
	}
	return sourceStr

}

// buildSignSourceStr100 构造签名源字符串 1.0.0版本, 返回值适配content-type没有的情况, 为第一版兜底, TODO 后续可去除该逻辑
func buildSignSourceStr100(c *gin.Context) string {
	// 收集 query和header参数
	params := make(map[string]string)
	keys := make([]string, 0)
	querys := c.Request.URL.Query()
	for k, _ := range querys {
		params[k] = querys[k][0]
		keys = append(keys, k)
	}
	for k, _ := range c.Request.Header {
		lowerKey := strings.ToLower(k)
		if !slice.ContainString(lowerKey, boot_config.GetSocConfig().Signature.Headers) {
			if "accept-language" == lowerKey {
				params[lowerKey] = c.Request.Header.Get(k)
				keys = append(keys, lowerKey)
			}
			continue
		}
		params[lowerKey] = c.Request.Header.Get(k)
		keys = append(keys, lowerKey)
	}
	// 排序
	sort.Strings(keys)

	// 签名源字符串 格式为：{timestampStr}&{排序后的keys对(除了时间戳和签名kv) key1=value1&key2=value2&key3=value3}[&{body}]&{签名key}
	sourceStr := c.Request.Header.Get(boot_consts.HeaderTimestampKey)

	for _, v := range keys {
		if v == boot_consts.HeaderSignKey || v == boot_consts.HeaderTimestampKey {
			continue
		}
		sourceStr += "&" + v + "=" + params[v]
	}
	body := model.GetHttpContext(c.Request.Context()).Body
	if body != "" {
		sourceStr += "&" + body
	}
	return sourceStr
}

// buildSignSourceStr090 构造签名源字符串 0.9.0版本, 返回值适配content-type没有的情况, 为第一版兜底, TODO 后续可去除该逻辑
func buildSignSourceStr090(c *gin.Context) string {
	// 收集 query和header参数
	params := make(map[string]string)
	keys := make([]string, 0)
	querys := c.Request.URL.Query()
	for k, _ := range querys {
		params[k] = querys[k][0]
		keys = append(keys, k)
	}
	for k, _ := range c.Request.Header {
		lowerKey := strings.ToLower(k)
		if !slice.ContainString(lowerKey, boot_config.GetSocConfig().Signature.Headers) {
			if "content-type" == lowerKey || "accept-language" == lowerKey {
				params[lowerKey] = c.Request.Header.Get(k)
				keys = append(keys, lowerKey)
			}
			continue
		}
		params[lowerKey] = c.Request.Header.Get(k)
		keys = append(keys, lowerKey)
	}
	// 排序
	sort.Strings(keys)

	// 签名源字符串 格式为：{timestampStr}&{排序后的keys对(除了时间戳和签名kv) key1=value1&key2=value2&key3=value3}[&{body}]&{签名key}
	sourceStr := c.Request.Header.Get(boot_consts.HeaderTimestampKey)

	for _, v := range keys {
		if v == boot_consts.HeaderSignKey || v == boot_consts.HeaderTimestampKey {
			continue
		}
		sourceStr += "&" + v + "=" + params[v]
	}
	body := model.GetHttpContext(c.Request.Context()).Body
	if body != "" {
		sourceStr += "&" + body
	}
	return sourceStr
}
