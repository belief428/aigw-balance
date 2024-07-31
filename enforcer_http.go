package aibalance

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/belief428/aigw-balance/utils"
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type Context struct {
	*http.Request
	http.ResponseWriter
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (this *Response) Marshal() []byte {
	_bytes, _ := json.Marshal(this)
	return _bytes
}

func (this *Response) Write() {

}

// getParams 获取参数信息
func getParams(enforcer *Enforcer) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write((&Response{
			Code:    0,
			Message: "ok",
			Data:    enforcer.GetParams(),
		}).Marshal())
	}
}

// setParams 设置参数信息
func setParams(enforcer *Enforcer) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := &Response{Code: 0, Message: "ok"}

		_params := make(map[string]interface{}, 0)

		_bytes, err := io.ReadAll(r.Body)

		if err != nil {
			resp.Code = -1
			resp.Message = err.Error()
			goto LOOP
		}
		json.Unmarshal(_bytes, &_params)

		if err = enforcer.SetParams(_params); err != nil {
			resp.Code = -1
			resp.Message = err.Error()
		}
	LOOP:
		if enforcer.watcher != nil {
			enforcer.watcher.GetParamsCallback()(_params)
		}
		w.Write(resp.Marshal())
	}
}

// getArchive 获取档案信息
func getArchive(enforcer *Enforcer) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := &Response{Code: 0, Message: "ok"}
		//_bytes, err := io.ReadAll(r.Body)
		//
		//if err != nil {
		//	resp.Code = -1
		//	resp.Message = err.Error()
		//	w.Write(resp.Marshal())
		//	return
		//}
		resp.Data = enforcer.data
		w.Write(resp.Marshal())
	}
}

// getHorizontalHistory 获取水平调控历史信息
func getHorizontalHistory(enforcer *Enforcer) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := &Response{Code: -1, Message: "ok", Data: make([]string, 0)}

		_bytes, err := io.ReadAll(r.Body)

		if err != nil {
			resp.Message = err.Error()
			w.Write(resp.Marshal())
			return
		}
		//utils.PathExists()
		_params := make(map[string]interface{}, 0)
		json.Unmarshal(_bytes, &_params)

		date, has := _params["date"]

		if !has {
			date = time.Now().Format("20060102")
		}
		date = strings.ReplaceAll(fmt.Sprintf("%v", date), "-", "")

		filepath := fmt.Sprintf("data/regulate/horizontal/%s.csv", date)

		isExist, _ := utils.FileExists(filepath)

		if !isExist {
			resp.Code = 0
			w.Write(resp.Marshal())
			return
		}
		var file *os.File

		if file, err = os.OpenFile(filepath, os.O_RDONLY, 0644); err != nil {
			resp.Message = err.Error()
			w.Write(resp.Marshal())
			return
		}
		reader := csv.NewReader(file)

		records := make([][]string, 0)

		if records, err = reader.ReadAll(); err != nil {
			resp.Message = err.Error()
			w.Write(resp.Marshal())
			return
		}
		resp.Code = 0

		if len(records) > 1 {
			resp.Data = records[1:]
		}
		w.Write(resp.Marshal())
	}
}

func (this *Enforcer) http() {
	defer func() {
		if err := recover(); err != nil {
			this.logger.Errorf("Aigw-balance http recover error：%v", err)
		}
		go this.http()
	}()
	mux := http.NewServeMux()
	// 启动静态文件服务
	mux.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("dist/css"))))
	mux.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("dist/js"))))
	mux.Handle("/favicon.ico", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		favicon, err := os.ReadFile("dist/favicon.ico")

		if err != nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
		_, _ = w.Write(favicon)
	}))
	mux.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) {
		_template, err := template.ParseFiles("dist/index.html")

		if err != nil {
			this.errorf("Aigw-balance http template parsefile error：%v", err)
			return
		}
		if err = _template.Execute(w, nil); err != nil {
			this.errorf("Aigw-balance http template execute error：%v", err)
			return
		}
	})
	// TODO：注入路由
	// context := &Context{}
	mux.HandleFunc("/api/v1/params", getParams(this))
	mux.HandleFunc("/api/v1/params/set", setParams(this))
	mux.HandleFunc("/api/v1/archive", getArchive(this))
	mux.HandleFunc("/api/v1/horizontal/history", getHorizontalHistory(this))

	if err := http.ListenAndServe(fmt.Sprintf(":%d", this.port), mux); err != nil {
		this.errorf("Aigw-balance http listen error：%v", err)
		return
	}
}
