package aibalance

import (
	"encoding/json"
	"github.com/belief428/aigw-balance/lib/orm"
	"github.com/belief428/aigw-balance/lib/queue"
	"github.com/belief428/aigw-balance/model"
	"github.com/belief428/aigw-balance/persist"
	"github.com/belief428/aigw-balance/plugin"
	"github.com/belief428/aigw-balance/utils"
	"gorm.io/gorm"
	"sync"
	"time"
)

// Enforcer 执行者
type Enforcer struct {
	port int // 端口

	watcher persist.IWatcher
	logger  persist.Logger // 日志模块

	params *plugin.Params
	//data   []EnforcerData[persist.IArchive] // 信息

	queue *queue.Instance

	archives map[string]map[string]model.ArchiveAttribute

	engine *gorm.DB

	time time.Time
}

// EnforcerData 执行者网关参数
type EnforcerData[T persist.IArchive] struct {
	persist.IGateway
	build []T // 水平平衡档案
	house []T // 垂直平衡档案
}

func (this *EnforcerData[T]) MarshalJSON() ([]byte, error) {
	data := map[string]interface{}{
		"name":        this.GetName(),
		"code":        this.GetCode(),
		"build_count": this.GetBuildCount(),
		"house_count": this.GetHouseCount(),

		"builds": this.build,
		"houses": this.house,
	}
	return json.Marshal(data)
}

var once sync.Once

type Option func(enforcer *Enforcer)

func WithPort(port int) Option {
	return func(enforcer *Enforcer) {
		enforcer.port = port
	}
}

func WithWatcher(watcher persist.IWatcher) Option {
	return func(enforcer *Enforcer) {
		enforcer.watcher = watcher
	}
}

func WithLogger(logger persist.Logger) Option {
	return func(enforcer *Enforcer) {
		enforcer.logger = logger
	}
}

func NewEnforcer(options ...Option) *Enforcer {
	_enforcer := &Enforcer{
		params:  plugin.NewParams(),
		watcher: NewWatcher(),
		//data:    make([]EnforcerData[persist.IArchive], 0),
		queue:    queue.NewInstance(),
		archives: map[string]map[string]model.ArchiveAttribute{},
		time:     time.Now(),
	}
	for _, option := range options {
		option(_enforcer)
	}
	return _enforcer
}

func (this *Enforcer) info(args ...interface{}) {
	if this.logger != nil {
		this.logger.Info(args...)
	}
}

func (this *Enforcer) infof(template string, args ...interface{}) {
	if this.logger != nil {
		//this.logger.Infof("Aigw balance "+template, args...)
		this.logger.Infof(template, args...)
	}
}

func (this *Enforcer) error(args ...interface{}) {
	if this.logger != nil {
		this.logger.Error(args...)
	}
}

func (this *Enforcer) errorf(template string, args ...interface{}) {
	if this.logger != nil {
		this.logger.Errorf(template, args...)
	}
}

func (this *Enforcer) sync(model interface{}, engine string, objs interface{}) {
	if !this.engine.Migrator().HasTable(model) {
		if err := this.engine.Set("gorm:table_options", "ENGINE="+engine).AutoMigrate(model); err != nil {
			this.errorf("Aigw-balance sync model error：%v", err)
		} else {
			_ = this.engine.Model(model).Create(objs).Error
		}
	}
}

func (this *Enforcer) load() {
	//{
	//	_params := new(model.Params)
	//	this.engine.Model(&model.Params{}).Where("`key` = ?", "params").First(_params)
	//
	//	if _params.ID > 0 {
	//		_bytes, _ := json.Marshal(_params.Value)
	//		json.Unmarshal(_bytes, this.params)
	//	}
	//}
	{
		_archives := make([]*model.Archive, 0)

		this.engine.Table((&model.Archive{}).TableName()).Find(&_archives)

		for _, v := range _archives {
			_, has := this.archives[v.GatewayCode]

			if !has {
				this.archives[v.GatewayCode] = map[string]model.ArchiveAttribute{
					v.GatewayCode: v.Attribute,
				}
			} else {
				this.archives[v.GatewayCode][v.Code] = v.Attribute
			}
		}
	}
}

func (this *Enforcer) Enforcer() error {
	//now := time.Now()

	once.Do(func() {
		utils.LoadConfig(this.params.Filepath(), this.params)
		// 获取数据引擎
		this.engine = orm.NewInstance().GetEngine()
		{
			//_bytes, _ := json.Marshal(this.params)
			////this.params
			//this.sync(&model.Params{}, "InnoDB", map[string]interface{}{
			//	"`key`":      "params",
			//	"`value`":    string(_bytes),
			//	"created_at": now,
			//	"updated_at": now,
			//})
			this.sync(&model.Archive{}, "InnoDB", nil)
			this.sync(&model.RegulateBuild{}, "Archive", nil)
			this.sync(&model.RegulateHouse{}, "Archive", nil)
		}
		this.load()

		_enforcerCache.engine = this.engine
		// 载入Http
		if this.port > 0 {
			go this.http()
		}
		// 载入进程
		go this.process()
		// 载入队列
		go this.consume()
	})
	return nil
}
