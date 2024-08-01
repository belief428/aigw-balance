package aibalance

import (
	"encoding/json"
	"fmt"
	"github.com/belief428/aigw-balance/model"
	"github.com/belief428/aigw-balance/persist"
	"html/template"
	"io"
	"net/http"
	"os"
)

type Context struct {
	*http.Request
	http.ResponseWriter
}

type Page struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
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
			enforcer.watcher.GetParamsCallbackFunc()(_params)
		}
		w.Write(resp.Marshal())
	}
}

// getArchive 获取档案信息
func getArchive(enforcer *Enforcer) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := &Response{Code: 0, Message: "ok"}

		_params := &struct {
			Code string `json:"code"`
			Kind int    `json:"kind"`
		}{}
		_bytes, err := io.ReadAll(r.Body)

		if err != nil {
			resp.Code = -1
			resp.Message = err.Error()
			w.Write(resp.Marshal())
			return
		}
		json.Unmarshal(_bytes, &_params)

		if enforcer.watcher == nil || enforcer.watcher.GetArchiveFunc == nil {
			w.Write(resp.Marshal())
			return
		}
		archives := enforcer.watcher.GetArchiveFunc()(&persist.WatcherArchiveParams{
			Code: fmt.Sprintf("%v", _params.Code),
			Kind: _params.Kind,
		})
		resp.Data = archives
		w.Write(resp.Marshal())
	}
}

// getHorizontalHistory 获取水平调控历史信息
func getHorizontalHistory(enforcer *Enforcer) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := &Response{Code: -1, Message: "ok", Data: make([]string, 0)}

		pagination := &Page{Page: 1, Limit: 10}

		_bytes, err := io.ReadAll(r.Body)

		if err != nil {
			resp.Message = err.Error()
			return
		}
		json.Unmarshal(_bytes, pagination)

		if pagination.Page <= 0 {
			pagination.Page = 1
		}
		query := enforcer.engine.Table((&model.RegulateBuild{}).TableName())

		var count int64

		if err = query.Count(&count).Error; err != nil {
			resp.Message = err.Error()
			return
		}
		out := make([]*model.RegulateBuild, 0)

		if err = query.Offset((pagination.Page - 1) * pagination.Limit).
			Limit(pagination.Limit).Find(&out).Error; err != nil {
			resp.Message = err.Error()
			w.Write(resp.Marshal())
			return
		}
		resp.Code = 0
		resp.Data = struct {
			Data  interface{} `json:"data"`
			Count int64       `json:"count"`
		}{Data: out, Count: count}
		w.Write(resp.Marshal())
	}
}

func (this *Enforcer) http() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Aigw-balance http recover error：%v\n", err)
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
