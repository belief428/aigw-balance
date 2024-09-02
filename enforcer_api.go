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
	//this.engine.Model(&model.Params{}).Where("`key` = ?", "params").Updates(map[string]interface{}{
	//	"`value`": string(_bytes), "updated_at": time.Now(),
	//})
	//return nil
	return utils.TouchJson(this.params.Filepath(), _bytes)
}

func (this *Enforcer) GetParams() map[string]interface{} {
	return utils.StructToMap(this.params)
}

func (this *Enforcer) GetVersion() string {
	return Version
}
