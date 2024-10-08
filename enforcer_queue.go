package aibalance

import (
	"github.com/belief428/aigw-balance/lib/queue"
	"github.com/belief428/aigw-balance/persist"
	"github.com/belief428/aigw-balance/plugin"
	"sync"
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

type EnforcerQueues struct {
	queue  *queue.Instance
	caches map[string]uint8

	logger persist.Logger
	locker *sync.RWMutex
}

var triggerCount = 2

var queues = make(map[string]*EnforcerQueues, 0)

func (this *EnforcerQueueData[T]) Priority() int {
	return 0
}

func (this *EnforcerQueueData[T]) Delay() int {
	return 0
}

func (this *EnforcerQueueData[T]) Call(args ...interface{}) {
	if this.watcher == nil || this.watcher.GetRegulateCallbackFunc() == nil {
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
		if (resp != nil && resp.GetStatus() == 1) || i == triggerCount {
			func(gatewayCode string, archive persist.IArchive) {
				_regulate := plugin.NewRegulate()
				_regulate.Name = archive.GetName()
				_regulate.Code = archive.GetCode()
				_regulate.RetTemp = archive.GetRetTemp()
				_regulate.PrevDeg = archive.GetDeg()
				_regulate.NextDeg = this.value
				_regulate.Status = -1
				_regulate.Remark = "通讯失败"

				if resp != nil {
					_regulate.Status = resp.GetStatus()
					_regulate.Remark = resp.GetRemark()
				}
				_regulate.CreatedAt = time.Now()

				var err error

				if this.kind == EnforcerKindForVertical { // 垂直计算
					err = _enforcerCache.saveVerticalRegulate(gatewayCode, _regulate)
				} else if this.kind == EnforcerKindForHorizontal { // 水平计算
					err = _enforcerCache.saveHorizontalRegulate(gatewayCode, _regulate)
				}
				if err != nil && this.logger != nil {
					this.logger.Errorf("Aigw-balance cache save：%d error：%v", this.kind, err)
				}
			}(this.gatewayCode, this.archive)
			break
		}
	}
}

func (this *EnforcerQueues) consume() {
	defer func() {
		if err := recover(); err != nil {
			if this.logger != nil {
				this.logger.Errorf("Aigw-balance consume recover error：%v", err)
			}
		}
		go this.consume()
	}()
	for {
		_data := this.queue.LPop()

		if _data == nil {
			time.Sleep(500 * time.Millisecond)
			continue
		}
		_data.Call()
	}
}

func (this *Enforcer) push(data *EnforcerQueueData[persist.IArchive]) {
	obj, has := queues[data.gatewayCode]

	if !has {
		queues[data.gatewayCode] = &EnforcerQueues{
			queue:  queue.NewInstance(),
			caches: make(map[string]uint8, 0),
			logger: this.logger,
			locker: new(sync.RWMutex),
		}
		go queues[data.gatewayCode].consume()
		queues[data.gatewayCode].queue.RPush(data)
		return
	}
	obj.queue.RPush(data)
}
