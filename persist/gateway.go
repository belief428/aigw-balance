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
	// SetName 设置名称
	SetName(name string)
	// GetName 获取名称
	GetName() string
	// SetCode 设置编号
	SetCode(code string)
	// GetCode 获取编号
	GetCode() string
	// SetRegulate 设置调控状态
	SetRegulate(regulate bool)
	// GetRegulate 获取调控状态
	GetRegulate() bool
	// SetWeight 设置权重
	SetWeight(weight float32)
	// GetWeight 获取权重
	GetWeight() float32
	// SetBuild 设置建筑附加信息
	SetBuild(build IArchiveBuild)
	// GetBuild 获取建筑附加信息
	GetBuild() IArchiveBuild
	// SetDeg 设置开度
	SetDeg(deg uint8)
	// GetDeg 获取开度
	GetDeg() uint8
	// SetSupTemp 设置供温
	SetSupTemp(supTemp float32)
	// GetSupTemp 获取供温
	GetSupTemp() float32
	// SetRetTemp 设置回温
	SetRetTemp(retTemp float32)
	// GetRetTemp 获取回温
	GetRetTemp() float32
	// SetRoomTemp 设置室温
	SetRoomTemp(roomTemp float32)
	// GetRoomTemp 获取室温
	GetRoomTemp() float32
	// SetLsl 设置瞬流
	SetLsl(lsl float32)
	// GetLsl 获取瞬流
	GetLsl() float32
	// SetRgl 设置瞬热
	SetRgl(rgl float32)
	// GetRgl 获取瞬热
	GetRgl() float32
}

type IArchiveBuild interface {
	// GetArea 获取面积
	GetArea() float32
}
