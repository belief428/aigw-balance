package aibalance

import (
	"fmt"
	"github.com/belief428/aigw-balance/persist"
	"time"
)

type Watcher struct {
	getArchiveFunc        func(params *persist.WatcherArchiveParams) []persist.IArchive
	regulateCallbackFunc  func(params *persist.WatcherRegulateParams) persist.IWatchRegulate
	setParamsCallbackFunc func(params map[string]interface{})
}

type WatcherRegulate struct {
	status int
	remark string
}

func (this *Watcher) GetArchiveFunc() func(params *persist.WatcherArchiveParams) []persist.IArchive {
	return this.getArchiveFunc
}

func (this *Watcher) SetArchiveFunc(function func(params *persist.WatcherArchiveParams) []persist.IArchive) {
	this.getArchiveFunc = function
}

func (this *Watcher) SetRegulateCallbackFunc(function func(params *persist.WatcherRegulateParams) persist.IWatchRegulate) {
	this.regulateCallbackFunc = function
}

func (this *Watcher) GetRegulateCallbackFunc() func(params *persist.WatcherRegulateParams) persist.IWatchRegulate {
	return this.regulateCallbackFunc
}

func (this *Watcher) SetParamsCallbackFunc(function func(params map[string]interface{})) {
	this.setParamsCallbackFunc = function
}

func (this *Watcher) GetParamsCallbackFunc() func(params map[string]interface{}) {
	return this.setParamsCallbackFunc
}

func NewWatcher() *Watcher {
	return &Watcher{
		getArchiveFunc: func(params *persist.WatcherArchiveParams) []persist.IArchive {
			fmt.Println("getArchiveFunc time：", time.Now(), " code：", params.Code, " kind：", params.Kind)
			return nil
		},
		regulateCallbackFunc: func(params *persist.WatcherRegulateParams) persist.IWatchRegulate {
			fmt.Println("regulateCallbackFunc time：", time.Now(), " code：", params.Code, " archiveCode：", params.ArchiveCode,
				" kind：", params.Kind, " value：", params.Value)
			return NewWatcherRegulate()
		},
		setParamsCallbackFunc: func(params map[string]interface{}) {
			fmt.Println("setParamsCallbackFunc time：", time.Now(), " params：", params)
		},
	}
}

func (this *WatcherRegulate) SetStatus(status int) {
	this.status = status
}

func (this *WatcherRegulate) GetStatus() int {
	return this.status
}

func (this *WatcherRegulate) SetRemark(remark string) {
	this.remark = remark
}

func (this *WatcherRegulate) GetRemark() string {
	return this.remark
}

func NewWatcherRegulate() *WatcherRegulate {
	return &WatcherRegulate{status: 1}
}
