package aibalance

import (
	"encoding/csv"
	"fmt"
	"github.com/belief428/aigw-balance/model"
	"github.com/belief428/aigw-balance/persist"
	"github.com/belief428/aigw-balance/utils"
	"io"
	"os"
	"sync"
	"time"
)

type EnforcerCache struct {
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

var _path = "data/regulate/horizontal"

// saveVerticalRegulate 垂直调控记录
func (this *EnforcerCache) saveVerticalRegulate(iGateway persist.IGateway, iArchive persist.IArchive, mRegulate *model.Regulate) error {
	return nil
}

// saveHorizontalRegulate 水平调控记录
func (this *EnforcerCache) saveHorizontalRegulate(gatewayCode string, mRegulate *model.Regulate) error {
	//this.locker.RLock()
	//defer this.locker.RUnlock()

	now := time.Now().Format("20060102")

	isNewWrite := false

	filepath := fmt.Sprintf("%s/%s.csv", _path, now)

	if isExist, err := utils.FileExists(filepath); err != nil {
		return err
	} else if !isExist {
		if this._horizontal != nil && this._horizontal.handle != nil {
			this._horizontal.handle.Close()
		}
		var _file *os.File

		if _file, err = utils.Create(fmt.Sprintf("%s/%s.csv", _path, now)); err != nil {
			return err
		}
		this._horizontal = &EnforcerCacheHandle{handle: _file, date: now}

		isNewWrite = true
	} else {
		if this._horizontal == nil || this._horizontal.handle == nil {
			var _file *os.File

			if _file, err = os.OpenFile(fmt.Sprintf("%s/%s.csv", _path, now), os.O_WRONLY|os.O_APPEND, 0666); err != nil {
				return err
			}
			this._horizontal = &EnforcerCacheHandle{handle: _file, date: now}
		}
	}
	// 创建csv writer
	writer := csv.NewWriter(this._horizontal.handle)

	var err error

	if isNewWrite {
		if err = writer.Write([]string{"网关", "地址", "模式", "回温", "调控前开度", "调控后开度", "状态", "备注信息", "调控时间"}); err != nil {
			return err
		}
	}
	if err = writer.Write([]string{gatewayCode, mRegulate.Name, mRegulate.Code,
		fmt.Sprintf("%.3f", mRegulate.RetTemp),
		fmt.Sprintf("%d", mRegulate.PrevDeg), fmt.Sprintf("%d", mRegulate.NextDeg),
		fmt.Sprintf("%d", mRegulate.Status),
		mRegulate.Remark,
		mRegulate.CreatedAt.Format("2006-01-02 15:04:05")}); err != nil {
		return err
	}
	writer.Flush()

	return nil
}
