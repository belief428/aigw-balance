package aibalance

type Watcher struct {
	calculateCallbackFunc func(kind int, value float32)
	setParamsCallbackFunc func(params map[string]interface{})
}

func (this *Watcher) SetCalculateCallback(function func(kind int, value float32)) {
	this.calculateCallbackFunc = function
}

func (this *Watcher) GetCalculateCallback() func(kind int, value float32) {
	return this.calculateCallbackFunc
}

func (this *Watcher) SetParamsCallback(function func(params map[string]interface{})) {
	this.setParamsCallbackFunc = function
}

func (this *Watcher) GetParamsCallback() func(params map[string]interface{}) {
	return this.setParamsCallbackFunc
}

func NewWatcher() *Watcher {
	return &Watcher{}
}
