package aibalance

import (
	"encoding/json"
	"fmt"
	"github.com/belief428/aigw-balance/lib/orm"
	"github.com/belief428/aigw-balance/model"
	"github.com/belief428/aigw-balance/persist"
	"github.com/belief428/aigw-balance/plugin"
	"github.com/belief428/aigw-balance/utils"
	"sync"
	"testing"
	"time"
)

func TestNewEnforcer(t *testing.T) {
	enforcer := NewEnforcer(WithPort(11118))

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

func TestEnforceOrm(t *testing.T) {
	func(a int) {
		t.Log(a)
	}(1)
	return
	_orm := orm.NewInstance()

	err := _orm.GetEngine().Create(&model.RegulateBuild{
		GatewayCode: "88372100",
		ArchiveCode: "11223344",
		ArchiveName: "测试",
		Params: []model.RegulateParam{
			{
				Key:   "ret_temp",
				Title: "回温",
				Value: fmt.Sprintf("%.3f", 11.22),
			},
		},
		PrevDeg: 11,
		NextDeg: 22,
		Status:  1,
		Remark:  "测试",
		Date:    time.Now(),
	}).Error
	t.Log(err)
}
