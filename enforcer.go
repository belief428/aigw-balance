package aibalance

import (
	"github.com/belief428/aigw-balance/model"
	"github.com/belief428/aigw-balance/persist"
	"github.com/belief428/aigw-balance/utils"
	"sync"
	"time"
)

// Enforcer 执行者
type Enforcer struct {
	port    int // 端口
	mode    int // 模式：1-追回温，2-追流量，
	watcher persist.IWatcher
	logger  persist.Logger

	maxCycle int

	params *model.Params
	data   []EnforcerData // 信息
	time   time.Time
}

// EnforcerData 执行者网关参数
type EnforcerData struct {
	persist.IGateway
	horizontalArchive []persist.IArchive // 水平平衡档案
	verticalArchive   []persist.IArchive // 垂直平衡档案
}

var once sync.Once

type Option func(enforcer *Enforcer)

func WithPort(port int) Option {
	return func(enforcer *Enforcer) {
		enforcer.port = port
	}
}

func WithMode(mode int) Option {
	return func(enforcer *Enforcer) {
		enforcer.mode = mode
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

func NewEnforcer(options ...Option) *Enforcer {
	_enforcer := &Enforcer{
		mode:    EnforcerModeForZHW,
		watcher: NewWatcher(),
		params:  new(model.Params),
		data:    make([]EnforcerData, 0),
		time:    time.Now(),
	}
	for _, option := range options {
		option(_enforcer)
	}
	return _enforcer
}

func (this *Enforcer) Enforcer() error {
	once.Do(func() {
		// 载入配置
		utils.LoadConfig(this.params.Filepath(), this.params)
		// 载入Http
		if this.port > 0 {
			go this.http()
		}
		// 载入进程
		go this.process()
	})
	return nil
}
