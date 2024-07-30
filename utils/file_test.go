package utils

import (
	"encoding/csv"
	"log"
	"os"
	"testing"
)

func TestCreate(t *testing.T) {
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
