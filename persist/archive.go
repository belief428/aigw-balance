package persist

type Archive interface {
	// SetName 设置名称
	SetName(name string)
	// GetName 获取名称
	GetName() string

	// SetBuild 设置建筑信息
	SetBuild(build ArchiveBuild)
	// GetBuild 获取建筑信息
	GetBuild() ArchiveBuild

	// SetRetTemp 设置回温
	SetRetTemp(value float32)
	// GetRetTemp 获取回温
	GetRetTemp() float32
}

type ArchiveBuild interface {
	// SetArea 设置面积
	SetArea(area float32)
	// GetArea 获取面积
	GetArea() float32
}
