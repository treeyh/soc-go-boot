package logger

import (
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/treeyh/soc-go-boot/boot_config"
	"github.com/treeyh/soc-go-boot/common/boot_consts"
	"sort"
	"strconv"
	"strings"
)

// langInfo 语言结构
type langInfo struct {

	// LangCode  语言编号
	LangCode string `json:"langCode"`

	// Weight  权重
	Weight float64 `json:"weight"`
}

// 语言对象排序支持 begin

type langInfoWrapper struct {
	langs []langInfo
	by    func(p, q *langInfo) bool
}
type SortBy func(p, q *langInfo) bool

func (pw langInfoWrapper) Len() int { // 重写 Len () 方法
	return len(pw.langs)
}
func (pw langInfoWrapper) Swap(i, j int) { // 重写 Swap () 方法
	pw.langs[i], pw.langs[j] = pw.langs[j], pw.langs[i]
}
func (pw langInfoWrapper) Less(i, j int) bool { // 重写 Less () 方法
	return pw.by(&pw.langs[i], &pw.langs[j])
}
func SortLangInfo(langs []langInfo, by SortBy) { // SortPerson 方法
	sort.Sort(langInfoWrapper{langs, by})
}

// 语言对象排序支持 end

//  formatRequestLang 格式化请求头语言  zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6
func formatRequestLang(ctx context.Context, acceptLang string) string {
	if !boot_config.GetSocConfig().I18n.Enable {
		return ""
	}
	if acceptLang == "" {
		return boot_config.GetSocConfig().I18n.DefaultLang
	}

	als := strings.Split(acceptLang, ",")

	langs := make([]langInfo, 0, len(als))
	for _, v := range als {
		ss := strings.Split(v, ";")
		if ss[0] == "" || strings.Contains(ss[0], "=") {
			continue
		}
		if len(ss) <= 1 {
			langs = append(langs, langInfo{
				LangCode: ss[0],
				Weight:   1,
			})
			continue
		}
		ss2 := strings.Split(ss[1], "=")
		if len(ss2) <= 1 {
			langs = append(langs, langInfo{
				LangCode: ss[0],
				Weight:   0.9,
			})
			continue
		}
		q, err := strconv.ParseFloat(ss2[1], 10)
		if err != nil {
			langs = append(langs, langInfo{
				LangCode: ss[0],
				Weight:   0.9,
			})
			continue
		}
		langs = append(langs, langInfo{
			LangCode: ss[0],
			Weight:   q,
		})
	}
	if len(langs) <= 0 {
		return boot_config.GetSocConfig().I18n.DefaultLang
	}

	sort.Sort(langInfoWrapper{langs, func(p, q *langInfo) bool {
		return q.Weight < p.Weight // Weight 递减排序
	}})

	langCode := langs[0].LangCode
	if strings.Contains(langCode, boot_consts.LangZhCn) || strings.Contains(langCode, boot_consts.LangZhChs) {
		return boot_consts.LangZhCn
	} else if strings.Contains(langCode, boot_consts.LangZhTw) || strings.Contains(langCode, boot_consts.LangZhHk) || strings.Contains(langCode, boot_consts.LangZhMo) || strings.Contains(langCode, boot_consts.LangZhCht) {
		return boot_consts.LangZhTw
	} else if strings.Contains(langCode, boot_consts.LangEn) {
		return boot_consts.LangEn
	}
	return boot_config.GetSocConfig().I18n.DefaultLang
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// getClientIp 获取请求ip
func getClientIp(c *gin.Context) string {
	ips := c.Request.Header.Get("X-Forwarded-For")
	ipss := strings.Split(ips, ",")
	if len(ipss) > 0 {
		ip := strings.TrimSpace(ipss[0])
		if ip != "" {
			return ip
		}
	}
	ips = c.Request.Header.Get("HTTP_CLIENT_IP")
	ipss = strings.Split(ips, ",")
	if len(ipss) > 0 {
		ip := strings.TrimSpace(ipss[0])
		if ip != "" {
			return ip
		}
	}
	ips = c.Request.Header.Get("HTTP_X_FORWARDED_FOR")
	ipss = strings.Split(ips, ",")
	if len(ipss) > 0 {
		ip := strings.TrimSpace(ipss[0])
		if ip != "" {
			return ip
		}
	}
	ip := strings.TrimSpace(c.Request.Header.Get("X-Real-IP"))
	if ip != "" {
		return ip
	}
	return c.ClientIP()
}
