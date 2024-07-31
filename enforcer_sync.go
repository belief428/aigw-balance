package aibalance

import (
	"github.com/belief428/aigw-balance/persist"
	"time"
)

const (
	// EnforcerModeForZHW 追回温
	EnforcerModeForZHW = iota + 1
	// EnforcerModeForZLL 追流量
	EnforcerModeForZLL
)

const (
	// EnforcerKindForVertical 垂直计算
	EnforcerKindForVertical int = iota + 1
	// EnforcerKindForHorizontal 水平计算
	EnforcerKindForHorizontal
)

// calc
// @Description:
// @param mode
// @param data
// @param report
// @param limit
// @return []persist.IArchive
// @return uint8
func calc(mode int, data []persist.IArchive, limit int) ([]persist.IArchive, uint8) {
	var value float32

	_length := len(data)

	out := make([]persist.IArchive, 0)

	if _length <= 0 {
		return out, 0
	}
	report := _length
	//|| (_length/report)*100 < 100-limit
	// 获取数值
	for _, v := range data {
		if mode == EnforcerModeForZHW {
			if value > 0 {
				value += v.GetRetTemp()
				continue
			}
			report--
		} else if mode == EnforcerModeForZLL {

		} else {
			return out, 0
		}
	}
	if (report/_length)*100 < 100-limit {
		return out, 0
	}
	return out, uint8(value / float32(_length))
}

// horizontal 水平计算
func (this *Enforcer) horizontal() {
	//builds := make([]persist.IArchive, 0)
	//buildCount := 0
	//
	//for _, v := range this.data {
	//	buildCount += v.GetBuildCount()
	//
	//	for _, val := range v.build {
	//		if !val.GetRegulate() {
	//			buildCount--
	//			continue
	//		}
	//		builds = append(builds, val)
	//	}
	//}
	//notices, value := calc(this.params.Mode, builds, buildCount, 13)
	//
	//if len(notices) <= 0 {
	//	return
	//}
	//for _, v := range this.data {
	//	for _, val := range v.build {
	//		if !val.GetRegulate() {
	//			continue
	//		}
	//		this.queue.RPush(&EnforcerQueueData[persist.IGateway, persist.IArchive]{
	//			gateway: v, archive: val, kind: EnforcerKindForHorizontal, value: value,
	//			watcher: this.watcher, logger: this.logger,
	//		})
	//	}
	//	//// 清空
	//	//v.build = make([]persist.IArchive, 0)
	//}
	return
}

// vertical 垂直计算
func (this *Enforcer) vertical() {
	// 获取档案信息
	if this.watcher == nil || this.watcher.GetArchiveFunc == nil {
		return
	}
	notices, value := calc(this.params.Mode, this.watcher.GetArchiveFunc()(&persist.WatcherArchiveParams{
		Code: "",
		Kind: EnforcerKindForVertical,
	}), 13)

	if len(notices) <= 0 {
		return
	}
	for _, val := range notices {
		this.queue.RPush(&EnforcerQueueData[persist.IGateway, persist.IArchive]{
			//gateway: this.data[0].IGateway,
			archive: val, kind: EnforcerKindForVertical, value: value,
			watcher: this.watcher, logger: this.logger,
		})
	}
	return
}

// rule 规则
// TODO：调控周期规则
// 一：分钟为单位，计算相差的小时*60+当前时间所在的分钟 % 周期
func (this *Enforcer) rule(startTime, nowTime time.Time, cycle int) bool {
	if startTime.After(nowTime) {
		return false
	}
	sub := nowTime.Sub(startTime)

	_hours := sub.Hours()

	return (int(_hours)*60+nowTime.Minute())%cycle == 0
}

// process 计算进程
func (this *Enforcer) process() {
	defer func() {
		if err := recover(); err != nil {
			this.errorf("Aigw-balance process recover error：%v", err)
		}
		go this.process()
	}()
	ticket := time.NewTicker(time.Second)

	defer ticket.Stop()

	local, _ := time.LoadLocation("Local")

	startTime := time.Date(this.time.Year(), this.time.Month(), this.time.Day(), this.time.Hour(), this.time.Minute(), 0, 0, local)

	var flagTime time.Time

	for {
		select {
		case <-ticket.C:
			now := time.Now()

			if flagTime.IsZero() {
				flagTime = now
				continue
			} else if flagTime.Minute() == now.Minute() { // 判断与当前分钟是否相同，相同就不处理
				continue
			}
			flagTime = now

			nowTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, local)

			// 垂直平衡模式
			if this.params.VerticalTime > 0 {
				if status := this.rule(startTime, nowTime, this.params.VerticalTime); status {
					this.vertical()
				}
			}
			// 水平平衡模式
			if this.params.HorizontalTime > 0 {
				if status := this.rule(startTime, nowTime, this.params.HorizontalTime); status {
					this.horizontal()
				}
			}
		}
	}
}
