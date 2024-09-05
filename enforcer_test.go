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

func TestEnforcerWg(t *testing.T) {
	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)

		go func(code int) {
			defer wg.Done()

			time.Sleep(time.Duration(code) * time.Second)

			t.Log("执行了：", code)
		}(i)
	}
	wg.Wait()
	t.Log("终于到我了")
}

func TestTime(t *testing.T) {
	local, _ := time.LoadLocation("Local")
	now := time.Date(2024, 9, 6, 0, 0, 0, 0, local)
	t.Log(now)
	befor := now.AddDate(0, 0, -3)
	t.Log(befor)

	var err error

	engine := orm.NewInstance().GetEngine()

	if err = engine.Where("date < ?", befor).Delete(&model.RegulateBuild{}).Error; err != nil {
		t.Logf("Aigw-balance crontab regulate build error：%v", err)
	}
	if err = engine.Where("date < ?", befor).Delete(&model.RegulateHouse{}).Error; err != nil {
		t.Logf("Aigw-balance crontab regulate house error：%v", err)
	}
	if engine.Mode == "Sqlite" {
		if err = engine.Exec("VACUUM").Error; err != nil {
			t.Logf("Aigw-balance crontab vacuum error：%v", err)
		}
	}
}
