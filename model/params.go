package model

type Params struct {
	VerticalTime   int `json:"vertical_time"`
	HorizontalTime int `json:"horizontal_time"`
}

func (this *Params) Filepath() string {
	return "data/params.json"
}
