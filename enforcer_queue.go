package aibalance

import (
	"github.com/belief428/aigw-balance/persist"
	"github.com/belief428/aigw-balance/plugin"
	"time"
)

type EnforcerQueueData[T persist.IArchive] struct {
	gatewayCode string
	archive     T
	kind        int // 调控类型，1：垂直，2：水平
	value       uint8

	watcher persist.IWatcher
	logger  persist.Logger
}

var triggerCount = 2

func (this *EnforcerQueueData[T]) Priority() int {
	return 0
}

func (this *EnforcerQueueData[T]) Delay() int {
	return 0
}

func (this *EnforcerQueueData[T]) Call(args ...interface{}) {
	if this.watcher == nil {
		return
	}
	// 最多只能执行两次
	for i := 1; i <= triggerCount; i++ {
		resp := this.watcher.GetRegulateCallbackFunc()(&persist.WatcherRegulateParams{
			Code:        this.gatewayCode,
			ArchiveCode: this.archive.GetCode(),
			Kind:        this.kind,
			Value:       this.value,
		})
		if resp.GetStatus() == 1 || i == triggerCount {
			func(gatewayCode string, archive persist.IArchive) {
				_regulate := plugin.NewRegulate()
				_regulate.Name = archive.GetName()
				_regulate.Code = archive.GetCode()
				_regulate.RetTemp = archive.GetRetTemp()
				_regulate.PrevDeg = archive.GetDeg()
				_regulate.NextDeg = this.value
				_regulate.Status = resp.GetStatus()
				_regulate.Remark = resp.GetRemark()
				_regulate.CreatedAt = time.Now()

				var err error

				if this.kind == EnforcerKindForVertical { // 垂直计算

				} else if this.kind == EnforcerKindForHorizontal { // 水平计算
					err = _enforcerCache.saveHorizontalRegulate(gatewayCode, _regulate)
				}
				if err != nil {
					this.logger.Errorf("Aigw-balance cache save：%d error：%v", this.kind, err)
				}
			}(this.gatewayCode, this.archive)
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
