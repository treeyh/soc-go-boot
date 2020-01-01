package controller

type IController interface {
	GetVersion() string
}

type BaseController struct {
}

func (BaseController) GetVersion() string {
	return ""
}
