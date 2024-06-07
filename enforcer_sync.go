package aibalance

import (
	"github.com/belief428/aigw-balance/persist"
	"time"
)

const (
	// EnforcerModeForZHW 追回温
	EnforcerModeForZHW = iota + 1
)

// vertical 垂直计算
func (this *Enforcer) vertical() {
	for _, v := range this.data {

		var value float32

		_length := len(v.verticalArchive)

		//this.infof()

		// 根据总数，上报数量限制
		if (_length/v.GetHorizontalCount())*100 < 100-13 {
			goto LOOP
		}
		// 获取数值
		for _, archive := range v.verticalArchive {
			if this.mode == EnforcerModeForZHW {
				value += archive.GetRetTemp()
			}
		}
		if this.watcher.GetCalculateCallback() != nil {
			this.watcher.GetCalculateCallback()(v.GetCode(), 1, value/float32(_length))
		}
	LOOP:
		// 清空
		v.verticalArchive = make([]persist.IArchive, 0)
	}
	return
}

// horizontal 水平计算
func (this *Enforcer) horizontal() {
	for _, v := range this.data {

		var value float32

		_length := len(v.horizontalArchive)

		//this.infof()

		// 根据总数，上报数量限制
		if (_length/v.GetHorizontalCount())*100 < 100-13 {
			goto LOOP
		}
		// 获取数值
		for _, archive := range v.horizontalArchive {
			if this.mode == EnforcerModeForZHW {
				value += archive.GetRetTemp()
			}
		}
		if this.watcher.GetCalculateCallback() != nil {
			this.watcher.GetCalculateCallback()(v.GetCode(), 2, value/float32(_length))
		}
	LOOP:
		// 清空
		v.horizontalArchive = make([]persist.IArchive, 0)
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
				if status := this.rule(startTime, nowTime, this.params.VerticalTime); !status {
					continue
				}
				this.vertical()
			}
			// 水平平衡模式
			if this.params.HorizontalTime > 0 {
				if status := this.rule(startTime, nowTime, this.params.HorizontalTime); !status {
					continue
				}
				this.horizontal()
			}
		}
	}
}
