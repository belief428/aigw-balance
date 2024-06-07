package aibalance

import "fmt"

type Watcher struct {
	calculateCallbackFunc func(code string, kind int, value float32)
	setParamsCallbackFunc func(params map[string]interface{})
}

func (this *Watcher) SetCalculateCallback(function func(code string, kind int, value float32)) {
	this.calculateCallbackFunc = function
}

func (this *Watcher) GetCalculateCallback() func(code string, kind int, value float32) {
	return this.calculateCallbackFunc
}

func (this *Watcher) SetParamsCallback(function func(params map[string]interface{})) {
	this.setParamsCallbackFunc = function
}

func (this *Watcher) GetParamsCallback() func(params map[string]interface{}) {
	return this.setParamsCallbackFunc
}

func NewWatcher() *Watcher {
	return &Watcher{
		calculateCallbackFunc: func(code string, kind int, value float32) {
			fmt.Println("calculateCallbackFunc", " code：", code, " kind：", kind, " value：", value)
		},
		setParamsCallbackFunc: func(params map[string]interface{}) {
			fmt.Println("setParamsCallbackFunc params：", params)
		},
	}
}
