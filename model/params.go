package model

type Params struct {
	Name           string `json:"name"`
	VerticalTime   int    `json:"vertical_time"`
	HorizontalTime int    `json:"horizontal_time"`
}

func (this *Params) Filepath() string {
	return "data/params.json"
}

func NewParams() *Params {
	return &Params{
		Name:           "AI-balance",
		VerticalTime:   10,
		HorizontalTime: 10,
	}
}
