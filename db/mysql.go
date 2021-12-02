package db

import (
	"encoding/json"
	"errors"
	"fmt"
	ccMicro "github.com/cqu20141693/sip-server/event"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"log"
	"time"
)

var MysqlDB *gorm.DB
var dbErr error

type Page struct {
	Page     int  `json:"page" binding:"required,min=1"`
	PageSize int  `json:"pageSize" binding:"required,max=100"`
	Desc     bool `json:"desc"`
}

type PageInfo struct {
	Total    int64
	Page     int
	PageSize int
	Data     interface{}
}

func NewPageInfo(total int64, page int, pageSize int, data interface{}) *PageInfo {
	return &PageInfo{Total: total, Page: page, PageSize: pageSize, Data: data}
}

type ConnPool struct {
	MaxIdleConn     int           `json:"maxIdleConn"`
	MaxOpenConn     int           `json:"maxOpenConn"`
	ConnMaxIdleTime time.Duration `json:"connMaxIdleTime"`
	ConnMaxLifetime time.Duration `json:"connMaxLifetime"`
}
type DataSource struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Url      string `json:"url"`
}

func (r DataSource) Check() error {
	if r.Username == "" {
		return errors.New("username is empty")
	}
	if r.Password == "" {
		return errors.New("password is empty")
	}
	if r.Url == "" {
		return errors.New("url is empty")
	}
	return nil
}
func (r *DataSource) Dsn() (dsn string) {

	return fmt.Sprintf("%s:%s@%s", r.Username, r.Password, r.Url)
}
func init() {
	ccMicro.RegisterHook(ccMicro.ConfigComplete, initDB)
}
func initDB() {

	defer func() {
		if err := recover(); err != nil {
			log.Fatal(err)
		}
	}()

	master, err := getDataSource("cc.datasource.master")
	if master == nil {
		log.Fatal(errors.New("master datasource must be set"))
	}
	slave, err := getDataSource("cc.datasource.slave")
	//  go-sql-driver作为驱动
	rwSplit := false
	if slave != nil {
		rwSplit = true
	}
	pool := getConnPool("cc.datasource.pool")
	masterOpen := mysql.Open(master.Dsn())
	MysqlDB, dbErr = gorm.Open(masterOpen, &gorm.Config{})
	if dbErr != nil {
		log.Fatal(dbErr)
	}
	if rwSplit {
		dbResolverCfg := dbresolver.Config{
			Sources:  []gorm.Dialector{masterOpen},
			Replicas: []gorm.Dialector{mysql.Open(slave.Dsn())},
			Policy:   dbresolver.RandomPolicy{}}
		readWritePlugin := dbresolver.Register(dbResolverCfg).SetMaxIdleConns(pool.MaxIdleConn).
			SetMaxOpenConns(pool.MaxOpenConn).SetConnMaxIdleTime(pool.ConnMaxIdleTime).SetConnMaxLifetime(pool.ConnMaxLifetime)
		err = MysqlDB.Use(readWritePlugin)
		if err != nil {
			log.Fatal(dbErr)
		}
	} else {
		// 连接池
		sqlDB, err := MysqlDB.DB()
		if err != nil {
			log.Fatal(err)
		}

		// SetMaxIdleConns 设置空闲连接池中连接的最大数量
		sqlDB.SetMaxIdleConns(pool.MaxIdleConn)

		// SetMaxOpenConns 设置打开数据库连接的最大数量。
		sqlDB.SetMaxOpenConns(pool.MaxOpenConn)
		sqlDB.SetConnMaxIdleTime(pool.ConnMaxIdleTime)

		// SetConnMaxLifetime 设置了连接可复用的最大时间。
		sqlDB.SetConnMaxLifetime(pool.ConnMaxLifetime)
	}
}

func getConnPool(key string) *ConnPool {
	sub := viper.Sub(key)
	sub.SetDefault("maxIdleConn", 2)
	sub.SetDefault("maxOpenConn", 20)
	sub.SetDefault("connMaxIdleTime", 30)
	sub.SetDefault("connMaxLifetime", 3600)
	maxIdleConn := sub.GetInt("maxIdleConn")
	maxOpenConn := sub.GetInt("maxOpenConn")
	connMaxIdleTime := sub.GetInt64("connMaxIdleTime")
	connMaxLifetime := sub.GetInt64("connMaxLifetime")
	return &ConnPool{MaxIdleConn: maxIdleConn, MaxOpenConn: maxOpenConn, ConnMaxIdleTime: time.Duration(connMaxIdleTime) * time.Second, ConnMaxLifetime: time.Duration(connMaxLifetime) * time.Second}

}

func getDataSource(key string) (*DataSource, error) {
	datasource := viper.GetStringMap(key)
	if len(datasource) == 0 {
		return nil, errors.New("data source not exist")
	}
	marshal, _ := json.Marshal(datasource)
	source := DataSource{}
	json.Unmarshal(marshal, &source)
	err := source.Check()
	if err != nil {
		log.Fatal(err)
	}
	return &source, err
}
