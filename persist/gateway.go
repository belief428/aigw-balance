package persist

type IGateway interface {
	// GetCode 获取编号
	GetCode() string
	// GetName 获取名称
	GetName() string
	// GetBuildCount 获取水平平衡调控数
	GetBuildCount() int
	// GetHouseCount 获取垂直平衡调控数
	GetHouseCount() int
}

type IArchive interface {
	// GetName 获取名称
	GetName() string
	// GetCode 获取编号
	GetCode() string
	// GetRegulate 获取调控状态
	GetRegulate() bool
	// GetWeight 获取权重
	GetWeight() float32
	// GetBuild 获取建筑附加信息
	GetBuild() IArchiveBuild
	// GetDeg 获取开度
	GetDeg() uint8
	// GetRetTemp 获取回温
	GetRetTemp() float32
	// GetRoomTemp 获取室温
	GetRoomTemp() float32
}

type IArchiveBuild interface {
	// GetArea 获取面积
	GetArea() float32
}
