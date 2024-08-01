package persist

// WatcherArchiveParams 档案参数
type WatcherArchiveParams struct {
	Code string `json:"code"` // 网关编号
	Kind int    `json:"kind"` //  1-垂直平衡、户阀信息，2-水平平衡、楼阀信息
}

// WatcherRegulateParams 调控参数
type WatcherRegulateParams struct {
	Code        string `json:"code"`         // 网关编号
	ArchiveCode string `json:"archive_code"` // 设备编号
	Kind        int    `json:"kind"`         //  1-垂直平衡、户阀信息，2-水平平衡、楼阀信息
	Value       uint8  `json:"value"`        // 反馈开度
}

type IWatcher interface {
	// GetArchiveFunc 获取档案函数
	GetArchiveFunc() func(*WatcherArchiveParams) []IArchive
	// SetArchiveFunc 设置档案函数
	SetArchiveFunc(func(*WatcherArchiveParams) []IArchive)
	// SetRegulateCallbackFunc 设置调控回调函数
	SetRegulateCallbackFunc(func(*WatcherRegulateParams) IWatchRegulate)
	// GetRegulateCallbackFunc 获取调控回调函数
	GetRegulateCallbackFunc() func(*WatcherRegulateParams) IWatchRegulate
	// SetParamsCallbackFunc 设置参数回调函数
	SetParamsCallbackFunc(func(params map[string]interface{}))
	// GetParamsCallbackFunc 获取参数回调函数
	GetParamsCallbackFunc() func(params map[string]interface{})
}

type IWatchRegulate interface {
	// GetStatus 获取状态
	GetStatus() int
	// GetRemark 获取备注
	GetRemark() string
}
