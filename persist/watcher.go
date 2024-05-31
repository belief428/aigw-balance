package persist

type IWatcher interface {
	SetCalculateCallback(func(kind int, value float32))
	GetCalculateCallback() func(kind int, value float32)

	SetParamsCallback(func(params map[string]interface{}))
	GetParamsCallback() func(params map[string]interface{})
}
