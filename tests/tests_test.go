package tests

import (
	"github.com/treeyh/soc-go-common/core/types"
	"testing"
	"time"
)

func TestTimeIn(t *testing.T) {
	now := types.Time(time.Now())
	t.Log(now)
	t.Log(now.ToTime().UnixNano())

	now5 := now.InByOffset(5 * 3600)
	t.Log(now5)
	t.Log(now5.ToTime().UnixNano())
}

func TestTimeIn2(t *testing.T) {
	now := types.UtcTime(time.Now())
	t.Log(now)
	t.Log(now.ToTime().UnixNano())

	now5 := now.InByOffset(5 * 3600)
	t.Log(now5)
	t.Log(now5.ToTime().UnixNano())
}
