package model

import "time"

// Regulate 调控数据模型
type Regulate struct {
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	Mode      string    `json:"mode"`          // 自动；手动
	RetTemp   float32   `json:"perv_ret_temp"` // 回温
	PrevDeg   uint8     `json:"prev_deg"`      // 上一个开度
	NextDeg   uint8     `json:"next_deg"`      // 下一个开度，即调控完的开度
	Status    int       `json:"status"`        // 状态（-1-失败，1-成功）
	Remark    string    `json:"remark"`        // 备注信息
	CreatedAt time.Time `json:"created_at"`    // 调控时间
}

func (this *Regulate) Filepath() string {
	return "data/regulate.csv"
}

func NewRegulate() *Regulate {
	return &Regulate{}
}
