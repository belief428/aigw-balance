package utils

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"testing"
)

func TestCreate(t *testing.T) {
	_file, err := Create("data/1/2/3/name.txt")
	defer _file.Close()
	t.Log(err)

	err = os.WriteFile("data/1/2/3/name.txt", nil, 0666)
	t.Log(err)
	return
	file, err := os.OpenFile("example.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// 创建csv writer
	writer := csv.NewWriter(file)

	// 追加一行数据
	err = writer.Write([]string{"Name", "Age", "City"})
	if err != nil {
		log.Fatal(err)
	}

	// 将缓冲的数据flush到CSV文件中
	writer.Flush()
}

func TestOpenFile(t *testing.T) {
	file, err := os.OpenFile("ip.txt", os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		t.Log(err)
		return
	}
	defer file.Close()

	_bytes := make([]byte, 0)

	if _bytes, err = io.ReadAll(file); err != nil {
		t.Log(err)
		return
	}
	t.Log(string(_bytes))

}
