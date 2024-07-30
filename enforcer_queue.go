package aibalance

import (
	"fmt"
	"github.com/belief428/aigw-balance/model"
	"github.com/belief428/aigw-balance/persist"
	"time"
)

type EnforcerQueueData[G persist.IGateway, T persist.IArchive] struct {
	gateway G
	archive T
	mode    string
	kind    int // 调控类型，1：垂直，2：水平
	value   uint8

	watcher persist.IWatcher
	logger  persist.Logger
}

var triggerCount = 2

func (this *EnforcerQueueData[G, T]) Priority() int {
	return 0
}

func (this *EnforcerQueueData[G, T]) Delay() int {
	return 0
}

func (this *EnforcerQueueData[G, T]) Call(args ...interface{}) {
	if this.watcher == nil {
		return
	}
	// 最多只能执行两次
	for i := 1; i <= triggerCount; i++ {
		resp := this.watcher.GetRegulateCallback()(this.gateway.GetCode(), this.archive.GetCode(), this.kind, this.value)

		if resp.GetStatus() == 1 || i == triggerCount {
			func(gateway persist.IGateway, archive persist.IArchive) {
				_regulate := model.NewRegulate()
				_regulate.Name = archive.GetName()
				_regulate.Code = archive.GetCode()
				_regulate.Mode = this.mode
				_regulate.RetTemp = archive.GetRetTemp()
				_regulate.PrevDeg = archive.GetDeg()
				_regulate.NextDeg = this.value
				_regulate.Status = resp.GetStatus()
				_regulate.Remark = resp.GetRemark()
				_regulate.CreatedAt = time.Now()

				var err error

				if this.kind == EnforcerKindForVertical { // 垂直计算

				} else if this.kind == EnforcerKindForHorizontal { // 水平计算
					err = _enforcerCache.saveHorizontalRegulate(gateway, _regulate)
				}
				if err != nil {
					fmt.Println(err)
					this.logger.Errorf("Aigw-balance cache save：%d error：%v", this.kind, err)
				}
			}(this.gateway, this.archive)
			break
		}
	}
}

// consume 开启队列
func (this *Enforcer) consume() {
	if this.queue == nil {
		return
	}
	defer func() {
		if err := recover(); err != nil {
			this.errorf("Aigw-balance consume recover error：%v", err)
		}
		go this.consume()
	}()
	for {
		_queue := this.queue.LPop()

		if _queue == nil {
			time.Sleep(time.Second)
			continue
		}
		_queue.Call()
	}
}
