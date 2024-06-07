package persist

type IGateway interface {
	// SetCode 设置编号
	SetCode(code string)
	// GetCode 获取编号
	GetCode() string
	// SetHorizontalCount 设置水平平衡调控数
	SetHorizontalCount(count int)
	// GetHorizontalCount 获取水平平衡调控数
	GetHorizontalCount() int
	// SetVerticalCount 设置垂直平衡调控数
	SetVerticalCount(count int)
	// GetVerticalCount 获取垂直平衡调控数
	GetVerticalCount() int
}

type IArchive interface {
	// SetName 设置名称
	SetName(name string)
	// GetName 获取名称
	GetName() string
	// SetType 设置类型（1：楼栋，2：房屋）
	SetType(_type int)
	// GetType 获取类型
	GetType() int

	// SetBuild 设置建筑信息
	SetBuild(build IArchiveBuild)
	// GetBuild 获取建筑信息
	GetBuild() IArchiveBuild

	// SetRetTemp 设置回温
	SetRetTemp(value float32)
	// GetRetTemp 获取回温
	GetRetTemp() float32
}

type IArchiveBuild interface {
	// SetArea 设置面积
	SetArea(area float32)
	// GetArea 获取面积
	GetArea() float32
}
