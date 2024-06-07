package aibalance

import (
	"github.com/belief428/aigw-balance/model"
	"testing"
)

func TestNewEnforcer(t *testing.T) {
	enforcer := NewEnforcer(WithPort(1111))

	enforcer.params = &model.Params{
		VerticalTime:   1,
		HorizontalTime: 2,
	}
	enforcer.RegisterHouse(NewArchive())

	if err := enforcer.Enforcer(); err != nil {
		t.Log("Enforcer：", err)
		return
	}
	select {}
}
