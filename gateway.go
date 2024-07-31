package aibalance

import (
	"encoding/json"
	"github.com/belief428/aigw-balance/persist"
)

type Gateway struct {
	code       string
	name       string
	buildCount int
	houseCount int
}

// SetCode 设置编码
func (this *Gateway) SetCode(code string) {
	this.code = code
}

// GetCode 获取编码
func (this *Gateway) GetCode() string {
	return this.code
}

// SetName 设置名称
func (this *Gateway) SetName(name string) {
	this.name = name
}

// GetName 获取名称
func (this *Gateway) GetName() string {
	return this.name
}

// SetBuildCount 设置楼栋数，即调控参数数
func (this *Gateway) SetBuildCount(count int) {
	this.buildCount = count
}

// GetBuildCount 获取楼栋数，即调控参数数
func (this *Gateway) GetBuildCount() int {
	return this.buildCount
}

// SetHouseCount 设置户数，即调控参数数
func (this *Gateway) SetHouseCount(count int) {
	this.houseCount = count
}

// GetHouseCount 获取户数，即调控参数数
func (this *Gateway) GetHouseCount() int {
	return this.houseCount
}

//func (this *Gateway) MarshalJSON() ([]byte, error) {
//	data := map[string]interface{}{
//		"code": this.code, "name": this.name,
//		"build_count": this.buildCount, "house_count": this.houseCount,
//	}
//	return json.Marshal(data)
//}

func NewGateway() *Gateway {
	return &Gateway{
		buildCount: 1, houseCount: 1,
	}
}

// Archive 档案信息
type Archive struct {
	name     string
	code     string
	regulate bool
	weight   float32

	build persist.IArchiveBuild

	deg               uint8
	retTemp, roomTemp float32
}

func (this *Archive) SetName(name string) {
	this.name = name
}

func (this *Archive) GetName() string {
	return this.name
}

func (this *Archive) SetCode(code string) {
	this.code = code
}

func (this *Archive) GetCode() string {
	return this.code
}

func (this *Archive) SetRegulate(value bool) {
	this.regulate = value
}

func (this *Archive) GetRegulate() bool {
	return this.regulate
}

func (this *Archive) SetWeight(weight float32) {
	this.weight = weight
}

func (this *Archive) GetWeight() float32 {
	return this.weight
}

func (this *Archive) SetBuild(build persist.IArchiveBuild) {
	this.build = build
}

func (this *Archive) GetBuild() persist.IArchiveBuild {
	return this.build
}

func (this *Archive) SetDeg(value uint8) {
	this.deg = value
}

func (this *Archive) GetDeg() uint8 {
	return this.deg
}

func (this *Archive) SetRetTemp(value float32) {
	this.retTemp = value
}

func (this *Archive) GetRetTemp() float32 {
	return this.retTemp
}

func (this *Archive) SetRoomTemp(value float32) {
	this.roomTemp = value
}

func (this *Archive) GetRoomTemp() float32 {
	return this.roomTemp
}

func (this *Archive) MarshalJSON() ([]byte, error) {
	data := map[string]interface{}{
		"name": this.name, "code": this.code,
		"regulate": this.regulate, "weight": this.weight,
		"build":    this.build,
		"ret_temp": this.retTemp,
	}
	return json.Marshal(data)
}

func NewArchive() *Archive {
	return &Archive{
		regulate: true,
	}
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

func (this *ArchiveBuild) MarshalJSON() ([]byte, error) {
	data := map[string]interface{}{
		"area": this.area,
	}
	return json.Marshal(data)
}

func NewArchiveBuild() *ArchiveBuild {
	return &ArchiveBuild{}
}
