package aibalance

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
)

func (this *Enforcer) http() {
	defer func() {
		if err := recover(); err != nil {
			this.logger.Errorf("Aigw-balance http recover error：%v", err)
		}
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
	if err := http.ListenAndServe(fmt.Sprintf(":%d", this.port), mux); err != nil {
		this.errorf("Aigw-balance http listen error：%v", err)
		return
	}
}
