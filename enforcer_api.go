package aibalance

import (
	"encoding/json"
	"errors"
	"github.com/belief428/aigw-balance/model"
	"github.com/belief428/aigw-balance/persist"
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

func (this *Enforcer) Register(gateway persist.IGateway) error {
	if gateway.GetCode() == "" {
		return errors.New("Please set gateway code ")
	}
	isExist := false
	// 判断是否存在
	for _, v := range this.data {
		if v.GetCode() == gateway.GetCode() {
			v.IGateway = gateway
			isExist = true
			break
		}
	}
	if !isExist {
		this.data = append(this.data, EnforcerData[persist.IArchive]{
			IGateway: gateway,
			build:    make([]persist.IArchive, 0),
			house:    make([]persist.IArchive, 0),
		})
	}
	return nil
}

func (this *Enforcer) RegisterBuild(code string, archive persist.IArchive) error {
	if archive.GetName() == "" {
		return errors.New("Please set archive name or address ")
	}
	// 判断是否存在
	for k, v := range this.data {
		if v.GetCode() == code {
			this.data[k].build = this.fill(v.build, archive)
			break
		}
	}
	return nil
}

func (this *Enforcer) RegisterHouse(code string, archive persist.IArchive) error {
	if archive.GetName() == "" {
		return errors.New("Please set archive name or address ")
	}
	for k, v := range this.data {
		if v.GetCode() == code {
			this.data[k].house = this.fill(v.house, archive)
			break
		}
	}
	return nil
}

func (this *Enforcer) ReportRegulateLog(code string, archive persist.IArchive) {
	model.NewRegulate()
}

func (this *Enforcer) fill(src []persist.IArchive, dist persist.IArchive) []persist.IArchive {
	isExist := false

	for key, iArchive := range src {
		if iArchive.GetName() == dist.GetName() {
			src[key] = dist
			isExist = true
			break
		}
	}
	if !isExist {
		src = append(src, dist)
	}
	return src
}
