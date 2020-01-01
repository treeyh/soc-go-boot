package request

type BaseReq struct {
	Operator int64 `json:"operator" validate:"omitempty,gt=0"`
}
