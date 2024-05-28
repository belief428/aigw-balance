package aibalance

import "github.com/belief428/aigw-balance/persist"

// 执行者

type Enforcer struct {
	mode int // 模式：1-追回温，2-追流量，

	watcher persist.IWatcher
}

type Option func(enforcer *Enforcer)

func WithMode() Option {
	return func(enforcer *Enforcer) {

	}
}

func NewEnforcer() *Enforcer {
	return &Enforcer{
		mode: 1,
	}
}

func (this *Enforcer) Enforcer() error {
	return nil
}
