package req

import (
	"github.com/treeyh/soc-go-common/core/types"
)

type UserReq struct {
	BaseReq

	Id            int64       `json:"id" validate:"omitempty,gt=0"`
	Name          string      `json:"name" validate:"required,max=64"`
	CountryCode   string      `json:"countryCode" validate:"required,max=8,contains=+"`
	Mobile        string      `json:"mobile" validate:"required,max=16"`
	Password      string      `json:"password" validate:"required"`
	Sex           int32       `json:"sex" validate:"omitempty,oneof=0 1 2"`
	Email         string      `json:"email" validate:"omitempty,email,checkEmailHost"`
	Address       string      `json:"address" validate:"omitempty,max=256"`
	Longitude     float64     `json:"longitude" validate:"omitempty,min=0,max=180"`
	Latitude      float64     `json:"latitude" validate:"omitempty,min=0,max=180"`
	LastLoginIp   string      `json:"lastLoginIp" validate:"omitempty,ip"`
	LastLoginTime *types.Time `json:"lastLoginTime"`
	Status        int32       `json:"status" validate:"omitempty,oneof=1 2"`
}

//type UserReq struct {
//	Id                  int64                   `json:"id" validate:"omitempty,gt=0"`
//	Name                string                  `json:"name" validate:"required,max=64"`
//	CountryCode         string                  `json:"countryCode" validate:"required,max=8,contains=+"`
//	Mobile              string                  `json:"mobile" validate:"required,max=16"`
//	Password            string                  `json:"password" validate:"required"`
//	Sex                 int32                   `json:"sex" validate:"omitempty,oneof=0 1 2"`
//	Email               string                  `json:"email" validate:"omitempty,email,checkUserEmail"`
//	Address             string	                `json:"address" validate:"omitempty,max=256"`
//	Longitude           float64                 `json:"longitude" validate:"omitempty,min=0,max=180"`
//	Latitude            float64                 `json:"latitude" validate:"omitempty,min=0,max=180"`
//	LastLoginIp         string                  `json:"lastLoginIp" validate:"omitempty,ip"`
//	LastLoginTime       *time.Time              `json:"lastLoginTime"`
//	Status              int32                   `json:"status" validate:"omitempty,oneof=1 2"`
//	Creator             int64                   `json:"creator"`
//	createTime          time.Time               `json:"createTime"`
//	Updator             int64                   `json:"updator"`
//	UpdateTime          time.Time               `json:"updateTime"`
//	Version             int32                   `json:"version"`
//	DelFlag             int32                   `json:"delFlag"`
//}
