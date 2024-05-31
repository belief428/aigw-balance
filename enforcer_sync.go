package aibalance

import "github.com/belief428/aigw-balance/persist"

type Archives []persist.IArchive

const (
	// EnforcerModeForZHW 追回温
	EnforcerModeForZHW = iota + 1
)

// vertical 垂直计算
func (this *Enforcer) vertical() {
	for _, v := range this.archive {
		//v.GetBuild().GetArea()
		v.GetRetTemp()
	}
	return
}

// horizontal 水平计算
func (this *Enforcer) horizontal() {
	return
}
