package persist

type IWatcher interface {
	// SetRegulateCallback 设置调控回调函数
	SetRegulateCallback(func(code, archiveCode string, kind int, value uint8) IWatchRegulate)
	// GetRegulateCallback 获取调控回调函数
	GetRegulateCallback() func(code, archiveCode string, kind int, value uint8) IWatchRegulate
	// SetParamsCallback 设置参数回调函数
	SetParamsCallback(func(params map[string]interface{}))
	// GetParamsCallback 获取参数回调函数
	GetParamsCallback() func(params map[string]interface{})
}

type IWatchRegulate interface {
	// GetStatus 获取状态
	GetStatus() int
	// GetRemark 获取备注
	GetRemark() string
}
