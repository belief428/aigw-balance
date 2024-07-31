package aibalance

import (
	"fmt"
	"github.com/belief428/aigw-balance/persist"
	"time"
)

type Watcher struct {
	regulateCallbackFunc  func(code, archiveCode string, kind int, value uint8) persist.IWatchRegulate
	setParamsCallbackFunc func(params map[string]interface{})
}

type WatcherRegulate struct {
	status int
	remark string
}

func (this *Watcher) SetRegulateCallback(function func(code, archiveCode string, kind int, value uint8) persist.IWatchRegulate) {
	this.regulateCallbackFunc = function
}

func (this *Watcher) GetRegulateCallback() func(code, archiveCode string, kind int, value uint8) persist.IWatchRegulate {
	return this.regulateCallbackFunc
}

func (this *Watcher) SetParamsCallback(function func(params map[string]interface{})) {
	this.setParamsCallbackFunc = function
}

func (this *Watcher) GetParamsCallback() func(params map[string]interface{}) {
	return this.setParamsCallbackFunc
}

func NewWatcher() *Watcher {
	return &Watcher{
		regulateCallbackFunc: func(code, archiveCode string, kind int, value uint8) persist.IWatchRegulate {
			fmt.Println("regulateCallbackFunc time：", time.Now(), " code：", code, " archiveCode：", archiveCode, " kind：", kind, " value：", value)
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
