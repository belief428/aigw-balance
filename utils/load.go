package utils

import (
	"encoding/json"
	"errors"
	"os"
	"path"
	//"gopkg.in/yaml.v2"
)

type callback func(interface{})

var loaders = map[string]func([]byte, interface{}) error{
	".json": LoadConfigFormJsonBytes,
	//".yaml": LoadConfigFromYamlBytes,
}

func LoadConfigFormJsonBytes(content []byte, obj interface{}) error {
	return json.Unmarshal(content, obj)
}

//func LoadConfigFromYamlBytes(content []byte, obj interface{}) error {
//	return yaml.Unmarshal(content, obj)
//}

func LoadConfig(file string, v interface{}) error {
	content, err := os.ReadFile(file)

	if err != nil {
		return err
	}
	loader, ok := loaders[path.Ext(file)]

	if !ok {
		return errors.New("Unknown File Type：" + path.Ext(file))
	}
	return loader(content, v)
}
