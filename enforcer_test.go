package aibalance

import (
	"encoding/json"
	"github.com/belief428/aigw-balance/persist"
	"github.com/belief428/aigw-balance/plugin"
	"github.com/belief428/aigw-balance/utils"
	"sync"
	"testing"
)

func TestNewEnforcer(t *testing.T) {
	enforcer := NewEnforcer(WithPort(1111))

	enforcer.params = &plugin.Params{
		VerticalTime:   1,
		HorizontalTime: 2,
	}
	if err := enforcer.Enforcer(); err != nil {
		t.Log("Enforcer：", err)
		return
	}
	select {}
}

type A struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (this *A) MarshalJSON() ([]byte, error) {
	data := map[string]interface{}{
		"name":   this.Name,
		"age":    18,
		"weight": 68,
	}
	return json.Marshal(data)
}

func TestJsonMarshal(t *testing.T) {
	dist := &A{
		Name: "Hanson",
		Age:  30,
	}
	_bytes, _ := json.Marshal(dist)
	t.Log(string(_bytes))

	utils.LoadConfig("data/params.json", dist)

	_bytes, _ = json.Marshal(dist)
	t.Log(string(_bytes))

	utils.LoadConfig("data/params1.json", dist)

	_bytes, _ = json.Marshal(dist)
	t.Log(string(_bytes))
}

func TestEnforcerCache(t *testing.T) {
	var src = &EnforcerCache{
		locker: new(sync.RWMutex),
	}
	gateway := NewGateway()
	gateway.code = "123"
	err := src.saveHorizontalRegulate(gateway.code, plugin.NewRegulate())

	if err != nil {
		t.Log(err)
	}
	function := func(params *persist.WatcherArchiveParams) []persist.IArchive {
		return nil
	}
	NewWatcher().SetArchiveFunc(function)
}
