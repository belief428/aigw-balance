package aibalance

import (
	"encoding/json"
	"github.com/belief428/aigw-balance/utils"
)

func (this *Enforcer) SetParams(params map[string]interface{}) error {
	_params := utils.StructToMap(this.params)

	for key, val := range params {
		_params[key] = val
	}
	_bytes, err := json.Marshal(_params)

	if err != nil {
		return err
	}
	if err = json.Unmarshal(_bytes, this.params); err != nil {
		return err
	}
	return utils.TouchJson(this.params.Filepath(), _bytes)
}

func (this *Enforcer) GetParams() map[string]interface{} {
	return utils.StructToMap(this.params)
}
