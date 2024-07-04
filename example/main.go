package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (this *Response) Marshal() []byte {
	_bytes, _ := json.Marshal(this)
	return _bytes
}

func _gin() {
	app := gin.Default()
	//app.Use(Cors())
	//app.Use(TimeoutHandle(time.Second * 60))

	app.StaticFS("/upload", http.Dir("./upload"))
	app.StaticFS("/dir", http.Dir("/"))
	// 注册路由
	//registerWeb(app)
	//registerApi(app)
	//
	//app.MaxMultipartMemory = 4 << 20

	app.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, &Response{
			Code:    0,
			Message: "ok",
			Data:    "Hello world",
		})
	})
	//return app
	app.Run(":3000")
}

func _http() {
	mux := http.NewServeMux()
	// TODO：注入路由
	// context := &Context{}
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write((&Response{
			Code:    0,
			Message: "ok",
			Data:    "Hello world",
		}).Marshal())
	})
	if err := http.ListenAndServe(fmt.Sprintf(":%d", 3000), mux); err != nil {

		return
	}
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf(" http recover error：%v\n", err)
		}
	}()
	_gin()
	//_http()
}
