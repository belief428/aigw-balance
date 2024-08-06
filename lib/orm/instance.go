package orm

import (
	"fmt"
	"github.com/belief428/aigw-balance/utils"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"io"
	"log"
	"os"
	"path"
	"time"

	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"gorm.io/gorm"
)

type Mysql struct {
	Address, Username, Password string
	Database, Parameters        string
}

func (this *Mysql) DSN() gorm.Dialector {
	return mysql.Open(fmt.Sprintf("%s:%s@tcp(%s)/%s?%s",
		this.Username, this.Password, this.Address, this.Database, this.Parameters,
	))
}

func (this *Mysql) Mode() string {
	return "Mysql"
}

type Sqlite struct {
	Path, Name string
}

func (this *Sqlite) DSN() gorm.Dialector {
	if isExist, _ := utils.PathExists(this.Path); !isExist {
		_ = utils.MkdirAll(this.Path)
	}
	url := path.Join(this.Path, this.Name)

	return sqlite.Open(url)
}

func (this *Sqlite) Mode() string {
	return "Sqlite"
}

type IEngine interface {
	DSN() gorm.Dialector
	Mode() string
}

type Engine struct {
	*gorm.DB
	Mode string
}

type Instance struct {
	engine *Engine

	debug                                   bool
	tablePrefix                             string
	singularTable                           bool
	maxIdleConns, maxOpenConns, maxLifetime int

	handler IEngine
}

type Option func(instance *Instance)

func WithTablePrefix(tablePrefix string) Option {
	return func(instance *Instance) {
		instance.tablePrefix = tablePrefix
	}
}

func WithSingularTable(singularTable bool) Option {
	return func(instance *Instance) {
		instance.singularTable = singularTable
	}
}

func WithMaxIdleConns(maxIdleConns int) Option {
	return func(instance *Instance) {
		instance.maxIdleConns = maxIdleConns
	}
}

func WithMaxOpenConns(maxOpenConns int) Option {
	return func(instance *Instance) {
		instance.maxOpenConns = maxOpenConns
	}
}

func WithMaxLifetime(maxLifetime int) Option {
	return func(instance *Instance) {
		instance.maxLifetime = maxLifetime
	}
}

func (this *Instance) init() *Instance {
	option := &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: this.tablePrefix,
			//SingularTable: this.singularTable,
		},
		Logger: logger.New(
			log.New(io.MultiWriter(), "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             time.Second,
				LogLevel:                  logger.Warn,
				Colorful:                  false,
				IgnoreRecordNotFoundError: true,
			},
		),
	}
	if this.debug {
		option.Logger = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             time.Second,
				LogLevel:                  logger.Info,
				Colorful:                  false,
				IgnoreRecordNotFoundError: true,
			},
		)
	}
	db, err := gorm.Open(this.handler.DSN(), option)

	if err != nil {
		panic("Orm Init Error：" + err.Error())
	}
	_db, _ := db.DB()
	_db.SetMaxIdleConns(this.maxIdleConns)
	_db.SetMaxOpenConns(this.maxOpenConns)
	_db.SetConnMaxLifetime(time.Duration(this.maxLifetime) * time.Second)
	this.engine = &Engine{
		DB:   db,
		Mode: this.handler.Mode(),
	}
	return this
}

func (this *Instance) GetEngine() *Engine {
	return this.engine
}

func NewInstance(option ...Option) *Instance {
	instance := &Instance{
		debug: true,
		handler: &Sqlite{
			Path: "data", Name: "app.db",
		},
		//handler: &Mysql{
		//	Address:    "127.0.0.1:3306",
		//	Username:   "appuser",
		//	Password:   "ABCabc01",
		//	Database:   "aigw_balance",
		//	Parameters: "charset=utf8&loc=Local&parseTime=true",
		//},
		maxLifetime:  3600,
		maxOpenConns: 200,
		maxIdleConns: 100,
	}
	for _, v := range option {
		v(instance)
	}
	return instance.init()
}
