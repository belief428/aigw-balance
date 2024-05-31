package aibalance

import (
	"fmt"
	"github.com/belief428/aigw-balance/persist"
	"html/template"
	"log"
	"net/http"
	"os"
)

// Enforcer 执行者
type Enforcer struct {
	port    int                // 端口
	mode    int                // 模式：1-追回温，2-追流量，
	archive []persist.IArchive // 档案信息
	watcher persist.IWatcher
}

type Option func(enforcer *Enforcer)

func WithPort(port int) Option {
	return func(enforcer *Enforcer) {
		enforcer.port = port
	}
}

func WithMode(mode int) Option {
	return func(enforcer *Enforcer) {
		enforcer.mode = mode
	}
}

func WithWatcher(watcher persist.IWatcher) Option {
	return func(enforcer *Enforcer) {
		enforcer.watcher = watcher
	}
}

func NewEnforcer(options ...Option) *Enforcer {
	_enforcer := &Enforcer{
		mode:    EnforcerModeForZHW,
		archive: make([]persist.IArchive, 0),
	}
	for _, option := range options {
		option(_enforcer)
	}
	return _enforcer
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./dist/index.html")

	if err != nil {
		log.Println(err)
	}

	err = t.Execute(w, nil)

	if err != nil {
		log.Println(err)
	}
}

func (this *Enforcer) http() {
	mux := http.NewServeMux()
	// 启动静态文件服务
	mux.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./dist/css"))))
	mux.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./dist/js"))))
	mux.Handle("/favicon.ico", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		favicon, err := os.ReadFile("./dist/favicon.ico")

		if err != nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
		w.Write(favicon)
	}))
	mux.HandleFunc("/admin", IndexHandler)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", this.port), mux); err != nil {
		fmt.Println(err)
	}
}

func (this *Enforcer) crond() {

}

func (this *Enforcer) Enforcer() error {
	if this.port > 0 {
		go this.http()
	}
	return nil
}
