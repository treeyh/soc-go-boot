package resp

type RespResult struct {
	Code      int         `json:"code" xml:"code"`
	Message   string      `json:"message" xml:"message"`
	Data      interface{} `json:"data" xml:"data"`
	Timestamp int64       `json:"timestamp" xml:"timestamp"`
}

type HttpJsonRespResult struct {
	Data interface{}

	HttpStatus int
}

type HttpTextRespResult struct {
	Text string

	HttpStatus int
}

type HttpProtoBufRespResult struct {
	Data interface{}

	HttpStatus int
}

type HttpXmlRespResult struct {
	Data interface{}

	HttpStatus int
}

type HttpFileRespResult struct {
	FilePath string
	FileName string

	HttpStatus int
}

type HttpHtmlRespResult struct {
	Name string
	Data interface{}

	HttpStatus int
}

type HttpRedirectRespResult struct {
	Location string

	HttpStatus int
}

type PageRespResult struct {
	Total int64       `json:"total"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
	Data  interface{} `json:"data"`
}

type ListRespResult struct {
	Total int         `json:"total"`
	Data  interface{} `json:"data"`
}
