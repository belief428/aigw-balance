package persist

type IWatcher interface {
	SetCalculateCallback(func(code string, kind int, value float32))
	GetCalculateCallback() func(code string, kind int, value float32)

	SetParamsCallback(func(params map[string]interface{}))
	GetParamsCallback() func(params map[string]interface{})
}
