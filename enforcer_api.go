package aibalance

import (
	"encoding/json"
	"github.com/belief428/aigw-balance/utils"
)

func (this *Enforcer) SetParams(params map[string]interface{}) error {
	_bytes, _ := json.Marshal(params)
	return utils.TouchJson("data/params.json", _bytes)
}

func (this *Enforcer) GetParams() (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}

func (this *Enforcer) RegisterBuild() {

}

func (this *Enforcer) RegisterHouse() {

}
