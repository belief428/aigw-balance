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

type Archives []persist.IArchive

func (this Archives) HandleCalc(mode, limit int) (bool, uint8) {
	var value float32

	_length := len(this)

	if _length <= 0 {
		return false, 0
	}
	report := _length

	for _, v := range this {
		if !v.GetRegulate() || v.GetCode() == "" {
			_length -= 1
			goto LOOP
		}
		if mode == EnforcerModeForZHW {
			if v.GetRetTemp() > 0 {
				value += v.GetRetTemp()
				continue
			}
		} else if mode == EnforcerModeForZLL {

		} else {
			return false, 0
		}
	LOOP:
		report--
	}
	if report <= 0 {
		return false, 0
	}
	if (report/_length)*100 < 100-limit {
		return false, 0
	}
	if value <= 0 {
		return false, 0
	}
	return true, uint8(value / float32(_length))
}

func (this *Enforcer) fillCalc(archive persist.IArchive, value uint8) uint8 {
	if archive.GetWeight() > 0 {
		value = uint8(float32(value) * archive.GetWeight())
	}
	return value
}

// vertical 垂直计算
func (this *Enforcer) vertical() {
	if this.watcher == nil || this.watcher.GetArchiveFunc() == nil {
		return
	}
	for _, v := range this.params.Gateways {
		archives := this.watcher.GetArchiveFunc()(&persist.WatcherArchiveParams{
			Code: v.Code, Kind: EnforcerKindForVertical,
		})
		for _, val := range archives {
			attribute := EnforcerArchive(this.archives).filter(v.Code, val.GetCode())
			val.SetRegulate(attribute.Regulate > 0)
			val.SetWeight(attribute.Weight)
		}
		valid, value := Archives(archives).HandleCalc(this.params.Mode, this.params.VerticalLimit)
		//valid, value := calc(this.params.Mode, archives, 13)
		if !valid {
			return
		}
		for _, val := range archives {
			if !val.GetRegulate() || val.GetCode() == "" {
				continue
			}
			_value := this.fillCalc(val, value)

			this.queue.RPush(&EnforcerQueueData[persist.IArchive]{
				gatewayCode: v.Code,
				archive:     val, kind: EnforcerKindForVertical, value: _value,
				watcher: this.watcher, logger: this.logger,
			})
		}
	}
	return
}

// horizontal 水平计算
func (this *Enforcer) horizontal() {
	if this.watcher == nil || this.watcher.GetArchiveFunc() == nil {
		return
	}
	builds := make([]persist.IArchive, 0)
	buildCodes := make(map[string][]persist.IArchive, 0)

	for _, v := range this.params.Gateways {
		archives := this.watcher.GetArchiveFunc()(&persist.WatcherArchiveParams{
			Code: v.Code, Kind: EnforcerKindForHorizontal,
		})
		for _, val := range archives {
			attribute := EnforcerArchive(this.archives).filter(v.Code, val.GetCode())
			val.SetRegulate(attribute.Regulate > 0)
			val.SetWeight(attribute.Weight)
		}
		builds = append(builds, archives...)

		buildCodes[v.Code] = archives
	}
	valid, value := Archives(builds).HandleCalc(this.params.Mode, this.params.HorizontalLimit)
	//valid, value := calc(this.params.Mode, builds, 13)
	if !valid {
		return
	}
	for key, val := range buildCodes {
		for _, _val := range val {
			if !_val.GetRegulate() || _val.GetCode() == "" {
				continue
			}
			_value := this.fillCalc(_val, value)

			this.queue.RPush(&EnforcerQueueData[persist.IArchive]{
				gatewayCode: key,
				archive:     _val, kind: EnforcerKindForHorizontal, value: _value,
				watcher: this.watcher, logger: this.logger,
			})
		}
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
				if valid := this.rule(startTime, nowTime, this.params.VerticalTime); valid {
					go this.vertical()
				}
			}
			// 水平平衡模式
			if this.params.HorizontalTime > 0 {
				if valid := this.rule(startTime, nowTime, this.params.HorizontalTime); valid {
					go this.horizontal()
				}
			}
		}
	}
}
