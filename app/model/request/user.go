package request

import "github.com/treeyh/soc-go-common/core/types"

type User struct {
	Id       int64       `json:"id"`
	Name     string      `json:"name" validate:"required"`
	Sex      int32       `json:"sex" validate:"required"`
	Birthday *types.Time `json:"birthday" validate:"required"`
	Weight   float64     `json:"weight" validate="min=10, max=1000"`
	Status   int32       `json:"status" validate="oneof=1 2"`
	Remark   *string     `json:"remark"`
}
