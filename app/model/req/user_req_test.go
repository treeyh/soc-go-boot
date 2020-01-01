package req

import (
	"fmt"
	"github.com/treeyh/soc-go-boot/app/model/validate"
	"github.com/treeyh/soc-go-common/core/utils/json"
	"testing"
)

func TestCheck(t *testing.T) {
	email := "a@bc.com"
	user := UserReq{
		BaseReq: BaseReq{
			Operator: 0,
		},
		Id:            0,
		Name:          "name",
		CountryCode:   "+86",
		Mobile:        "16666666666",
		Password:      "asssf",
		Sex:           1,
		Email:         email,
		Address:       "",
		Longitude:     2000,
		Latitude:      2000,
		LastLoginIp:   "",
		LastLoginTime: nil,
		Status:        0,
	}
	fmt.Println(json.ToJson(user))

	msg := validate.ValidateObject(user)
	if msg != nil {
		fmt.Println(*msg)
	}
}
