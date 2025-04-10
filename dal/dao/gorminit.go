package dao

import (
	"github.com/dog-frame/config"
	"github.com/dog-frame/dal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"time"
)

var _DB *gorm.DB

// DB 返回实例
func DB() *gorm.DB {
	return _DB
}

func init() {
	_DB = initDB(config.Database.Master)
}

func getDialector(t, dsn string) gorm.Dialector {
	//switch t { 项目数据库需要加载多数据源时去掉注释
	//case "postgres":
	//	return postgres.Open(dsn)
	//default:
	//	return mysql.Open(dsn)
	//}
	return mysql.Open(dsn)
}

func initDB(option config.DbConnectOption) *gorm.DB {
	db, err := gorm.Open(
		getDialector(option.Type, option.DSN),
		&gorm.Config{
			Logger: NewGormLogger(),
		},
	)
	if err != nil {
		panic(err)
	}

	// 配置读写分离策略
	err = db.Use(
		dbresolver.Register(dbresolver.Config{
			Sources:  []gorm.Dialector{getDialector(option.Type, option.DSN)},                               // 主库（写库）
			Replicas: []gorm.Dialector{getDialector(config.Database.Slave.Type, config.Database.Slave.DSN)}, // 从库（读库）
			Policy:   dbresolver.RandomPolicy{},                                                             // 负载均衡策略：随机选取一个从库
		}).
			SetConnMaxIdleTime(10 * time.Second). // 设置最大空闲连接时间
			SetConnMaxLifetime(1 * time.Hour).    // 连接最大存活时间
			SetMaxIdleConns(10).                  // 最大空闲连接数
			SetMaxOpenConns(100),                 // 最大打开连接数
	)

	sqlDb, _ := db.DB()

	if err = sqlDb.Ping(); err != nil {
		panic(err)
	}
	if err = db.AutoMigrate(model.Demo{}); err != nil {
		panic(err)
	}
	return db
}

// SetDBConn 设置连接对象 -- 只用在单测中把DB连接改成sqlMock的DB连接
func SetDBConn(conn *gorm.DB) {
	_DB = conn
}
