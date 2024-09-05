package aibalance

import (
	"github.com/belief428/aigw-balance/model"
	"github.com/belief428/aigw-balance/utils"
	"time"
)

func (this *Enforcer) crontab() {
	defer func() {
		if err := recover(); err != nil {
			this.errorf("Aigw-balance crontab recover error：%v", err)
		}
		go this.process()
	}()
	ticket := time.NewTicker(time.Second)

	defer ticket.Stop()

	for {
		select {
		case <-ticket.C:
			now := time.Now()

			if utils.IsDawn(now) {
				// 删除多余的日志
				_time := now.AddDate(0, 0, -this.params.RegulateSaveCycle)

				var err error

				if err = this.engine.Where("date < ?", _time).Delete(&model.RegulateBuild{}).Error; err != nil {
					this.errorf("Aigw-balance crontab regulate build error：%v", err)
				}
				if err = this.engine.Where("date < ?", _time).Delete(&model.RegulateHouse{}).Error; err != nil {
					this.errorf("Aigw-balance crontab regulate house error：%v", err)
				}
				if this.engine.Mode == "Sqlite" {
					if err = this.engine.Exec("VACUUM").Error; err != nil {
						this.errorf("Aigw-balance crontab exec vacuum error：%v", err)
					}
				}
			}
		}
	}
}
