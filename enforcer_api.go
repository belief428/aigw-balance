package aibalance

import (
	"encoding/json"
	"errors"
	"github.com/belief428/aigw-balance/persist"
	"github.com/belief428/aigw-balance/utils"
)

func (this *Enforcer) SetParams(params map[string]interface{}) error {
	_bytes, err := json.Marshal(params)

	if err != nil {
		return err
	}
	if err = json.Unmarshal(_bytes, this.params); err != nil {
		return err
	}
	return utils.TouchJson(this.params.Filepath(), _bytes)
}

func (this *Enforcer) GetParams() (map[string]interface{}, error) {
	return utils.StructToMap(this.params), nil
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
		this.data = append(this.data, EnforcerData{
			IGateway:          gateway,
			horizontalArchive: make([]persist.IArchive, 0),
			verticalArchive:   make([]persist.IArchive, 0),
		})
	}
	return nil
}

func (this *Enforcer) RegisterBuild(code string, archive persist.IArchive) error {
	if archive.GetName() == "" {
		return errors.New("Please set archive name or address ")
	}
	// 判断是否存在
	for _, v := range this.data {
		if v.GetCode() == code {
			isExist := false

			for _, iArchive := range v.horizontalArchive {
				if iArchive.GetName() == archive.GetName() {
					iArchive = archive
					isExist = true
					break
				}
			}
			if !isExist {
				v.horizontalArchive = append(v.horizontalArchive, archive)
			}
			break
		}
	}
	return nil
}

func (this *Enforcer) RegisterHouse(code string, archive persist.IArchive) error {
	if archive.GetName() == "" {
		return errors.New("Please set archive name or address ")
	}
	// 判断是否存在
	for _, v := range this.data {
		if v.GetCode() == code {
			isExist := false

			for _, iArchive := range v.verticalArchive {
				if iArchive.GetName() == archive.GetName() {
					iArchive = archive
					isExist = true
					break
				}
			}
			if !isExist {
				v.verticalArchive = append(v.verticalArchive, archive)
			}
			break
		}
	}
	return nil
}
