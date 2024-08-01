package aibalance

import (
	"fmt"
	"github.com/belief428/aigw-balance/model"
	"github.com/belief428/aigw-balance/persist"
	"github.com/belief428/aigw-balance/plugin"
	"gorm.io/gorm"
	"io"
	"sync"
	"time"
)

type EnforcerCache struct {
	engine *gorm.DB

	_vertical   *EnforcerCacheHandle
	_horizontal *EnforcerCacheHandle

	locker *sync.RWMutex
}

type EnforcerCacheHandle struct {
	handle io.WriteCloser
	date   string
}

var _enforcerCache = &EnforcerCache{
	locker: new(sync.RWMutex),
}

// saveVerticalRegulate 垂直调控记录
func (this *EnforcerCache) saveVerticalRegulate(iGateway persist.IGateway, iArchive persist.IArchive, mRegulate *plugin.Regulate) error {
	return nil
}

// saveHorizontalRegulate 水平调控记录
func (this *EnforcerCache) saveHorizontalRegulate(gatewayCode string, mRegulate *plugin.Regulate) error {
	this.locker.Lock()
	defer this.locker.Unlock()

	return this.engine.Create(&model.RegulateBuild{
		Code:        gatewayCode,
		ArchiveCode: mRegulate.Code,
		ArchiveName: mRegulate.Name,
		Params: []model.RegulateParams{
			{
				Key:   "ret_temp",
				Title: "回温",
				Value: fmt.Sprintf("%.3f", mRegulate.RetTemp),
			},
		},
		PrevDeg: mRegulate.PrevDeg,
		NextDeg: mRegulate.NextDeg,
		Status:  mRegulate.Status,
		Remark:  mRegulate.Remark,
		Date:    time.Now(),
	}).Error
}
