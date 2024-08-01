package model

type Archive struct {
}

func (*Archive) TableName() string {
	return "archive"
}
