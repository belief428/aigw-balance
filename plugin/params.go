package plugin

type Params struct {
	Mode            int             `json:"mode"` // 调控模式
	Name            string          `json:"name"`
	VerticalTime    int             `json:"vertical_time"`
	VerticalLimit   int             `json:"vertical_limit"`
	HorizontalTime  int             `json:"horizontal_time"`
	HorizontalLimit int             `json:"horizontal_limit"`
	Gateways        []ParamsGateway `json:"gateways"`
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
		Mode:            1,
		Name:            "AI-balance",
		VerticalTime:    10,
		VerticalLimit:   13,
		HorizontalTime:  10,
		HorizontalLimit: 13,
		Gateways:        make([]ParamsGateway, 0),
	}
}
