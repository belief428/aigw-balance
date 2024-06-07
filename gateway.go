package aibalance

import "github.com/belief428/aigw-balance/persist"

type Gateway struct {
	code            string
	horizontalCount int
	verticalCount   int
}

// SetCode 设置编码
func (this *Gateway) SetCode(code string) {
	this.code = code
}

// GetCode 获取编码
func (this *Gateway) GetCode() string {
	return this.code
}

// SetHorizontalCount 设置水平平衡调控数
func (this *Gateway) SetHorizontalCount(count int) {
	this.horizontalCount = count
}

// GetHorizontalCount 获取水平平衡调控数
func (this *Gateway) GetHorizontalCount() int {
	return this.horizontalCount
}

// SetVerticalCount 设置垂直平衡调控数
func (this *Gateway) SetVerticalCount(count int) {
	this.verticalCount = count
}

// GetVerticalCount 获取垂直平衡调控数
func (this *Gateway) GetVerticalCount() int {
	return this.verticalCount
}

func NewGateway() *Gateway {
	return &Gateway{}
}

// Archive 档案信息
type Archive struct {
	name  string
	_type int
	build persist.IArchiveBuild

	retTemp float32
}

func (this *Archive) SetName(name string) {
	this.name = name
}

func (this *Archive) GetName() string {
	return this.name
}

func (this *Archive) SetType(_type int) {
	this._type = _type
}

func (this *Archive) GetType() int {
	return this._type
}

func (this *Archive) SetBuild(build persist.IArchiveBuild) {
	this.build = build
}

func (this *Archive) GetBuild() persist.IArchiveBuild {
	return this.build
}

func (this *Archive) SetRetTemp(value float32) {
	this.retTemp = value
}

func (this *Archive) GetRetTemp() float32 {
	return this.retTemp
}

func NewArchive() *Archive {
	return &Archive{}
}

// ArchiveBuild 档案建筑信息
type ArchiveBuild struct {
	area float32
}

func (this *ArchiveBuild) SetArea(area float32) {
	this.area = area
}

func (this *ArchiveBuild) GetArea() float32 {
	return this.area
}

func NewArchiveBuild() *ArchiveBuild {
	return &ArchiveBuild{}
}
