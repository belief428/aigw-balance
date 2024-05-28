package aibalance

import "github.com/belief428/aigw-balance/persist"

type Archives []persist.Archive

// Vertical 垂直计算
func (this Archives) Vertical(mode int) (bool, float32) {
	for _, v := range this {
		//v.GetBuild().GetArea()
		v.GetRetTemp()
	}
	return false, 0
}

// Horizontal 水平计算
func (this Archives) Horizontal(mode int) (bool, float32) {
	return false, 0
}
