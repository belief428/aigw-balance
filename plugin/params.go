package plugin

type Params struct {
	Mode           int             `json:"mode"` // 调控模式
	Name           string          `json:"name"`
	VerticalTime   int             `json:"vertical_time"`
	HorizontalTime int             `json:"horizontal_time"`
	Gateways       []ParamsGateway `json:"gateways"`
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
		Mode:           1,
		Name:           "AI-balance",
		VerticalTime:   10,
		HorizontalTime: 10,
		Gateways:       make([]ParamsGateway, 0),
	}
}
