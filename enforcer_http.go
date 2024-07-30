package aibalance

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
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
		if err = json.Unmarshal(_bytes, &_params); err != nil {
			resp.Code = -1
			resp.Message = err.Error()
			goto LOOP
		}
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

// getHistory 获取调控历史信息
func getHistory(enforcer *Enforcer) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := &Response{Code: 0, Message: "ok"}

		_bytes, err := io.ReadAll(r.Body)

		if err != nil {
			resp.Code = -1
			resp.Message = err.Error()
			w.Write(resp.Marshal())
			return
		}
		_params := make(map[string]interface{}, 0)

		if err = json.Unmarshal(_bytes, &_params); err != nil {
			resp.Code = -1
			resp.Message = err.Error()
			w.Write(resp.Marshal())
			return
		}
		code, has := _params["code"]

		if !has {
			resp.Code = -1
			resp.Message = "Please choose code"
			w.Write(resp.Marshal())
			return
		}
		fmt.Println(code)
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

	if err := http.ListenAndServe(fmt.Sprintf(":%d", this.port), mux); err != nil {
		this.errorf("Aigw-balance http listen error：%v", err)
		return
	}
}
