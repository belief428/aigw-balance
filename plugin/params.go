package plugin

type Params struct {
	Mode              int             `json:"mode"` // 调控模式
	Name              string          `json:"name"`
	VerticalTime      int             `json:"vertical_time"`
	VerticalLimit     int             `json:"vertical_limit"`
	HorizontalTime    int             `json:"horizontal_time"`
	HorizontalLimit   int             `json:"horizontal_limit"`
	RegulateSaveCycle int             `json:"regulate_save_cycle"` // 调控日志保存周期、天
	Gateways          []ParamsGateway `json:"gateways"`
}

type ParamsGateway struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

func (this *Params) Filepath() string {
	return "data/params.json"
}

func NewParams() *Params {
	return &Params{
		Mode:              1,
		Name:              "AI-balance",
		VerticalTime:      60,
		VerticalLimit:     10,
		HorizontalTime:    60,
		HorizontalLimit:   10,
		RegulateSaveCycle: 3,
		Gateways:          make([]ParamsGateway, 0),
	}
}
