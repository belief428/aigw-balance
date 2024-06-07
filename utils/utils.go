package utils

import (
	"bytes"
	"encoding/json"
)

func StructToMap(src interface{}) map[string]interface{} {
	out := make(map[string]interface{}, 0)

	marshalContent, err := json.Marshal(src)

	if err != nil {
		return out
	}
	d := json.NewDecoder(bytes.NewReader(marshalContent))

	//d.UseNumber() // 设置将float64转为一个number
	if err = d.Decode(&out); err != nil {
		return out
	}
	return out
}
