package db

import (
	ccMicro "github.com/cqu20141693/sip-server/event"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"log"
	"time"
)

/*
redis常用函数列表：

Set - 设置一个key的值
Get - 查询key的值
GetSet - 设置一个key的值，并返回这个key的旧值
SetNX - 如果key不存在，则设置这个key的值
MGet - 批量查询key的值
MSet - 批量设置key的值
Incr,IncrBy,IncrByFloat - 针对一个key的数值进行递增操作
Decr,DecrBy - 针对一个key的数值进行递减操作
Del - 删除key操作，可以批量删除
Expire - 设置key的过期时间



*/
var RedisDB *redis.Client
var ClusterDB *redis.ClusterClient

func init() {
	ccMicro.RegisterHook(ccMicro.ConfigComplete, initRedisDB)
}

func initRedisDB() {

	sub := viper.Sub("cc.redis")
	addr := sub.GetString("addr")
	if addr == "" {
		sentinelNodes := sub.GetStringSlice("sentinel.nodes")
		if sentinelNodes == nil || len(sentinelNodes) == 0 {
			nodes := sub.GetStringSlice("cluster.nodes")
			if nodes == nil || len(nodes) == 0 {
				log.Fatal("redis addr not config")
			} else { //cluster
				ClusterDB = redis.NewClusterClient(configClusterOptions())
			}
		} else { // sentinel
			RedisDB = redis.NewFailoverClient(configSentinelOptions())
		}
	} else { // redis server
		RedisDB = redis.NewClient(configRedisOptions())
	}
}

func configClusterOptions() *redis.ClusterOptions {
	sub := viper.Sub("cc.redis")
	options := redis.ClusterOptions{}
	options.Addrs = sub.GetStringSlice("cluster.nodes")
	options.Password = sub.GetString("password")
	options.Username = sub.GetString("username")
	options.DialTimeout = sub.GetDuration("conn-timeout") * time.Second
	options.ReadTimeout = sub.GetDuration("read-timeout") * time.Second
	options.PoolTimeout = sub.GetDuration("pool-timeout") * time.Second
	options.IdleTimeout = sub.GetDuration("idle-timeout") * time.Second
	options.MaxRetries = sub.GetInt("retry")
	options.PoolSize = sub.GetInt("pool-size")
	options.MinIdleConns = sub.GetInt("min-idle-conn")
	return &options
}

func configSentinelOptions() *redis.FailoverOptions {
	sub := viper.Sub("cc.redis")
	options := redis.FailoverOptions{}
	options.MasterName = sub.GetString("sentinel.master")
	options.SentinelAddrs = sub.GetStringSlice("sentinel.nodes")
	options.SentinelUsername = sub.GetString("sentinel.username")
	options.SentinelPassword = sub.GetString("sentinel.password")
	options.Password = sub.GetString("password")
	options.DB = sub.GetInt("database")
	options.Username = sub.GetString("username")
	options.DialTimeout = sub.GetDuration("conn-timeout") * time.Second
	options.ReadTimeout = sub.GetDuration("read-timeout") * time.Second
	options.PoolTimeout = sub.GetDuration("pool-timeout") * time.Second
	options.IdleTimeout = sub.GetDuration("idle-timeout") * time.Second
	options.MaxRetries = sub.GetInt("retry")
	options.PoolSize = sub.GetInt("pool-size")
	options.MinIdleConns = sub.GetInt("min-idle-conn")
	return &options
}

func configRedisOptions() *redis.Options {
	sub := viper.Sub("cc.redis")
	options := redis.Options{}
	options.Addr = sub.GetString("addr")
	options.Password = sub.GetString("password")
	options.DB = sub.GetInt("database")
	options.Username = sub.GetString("username")
	options.DialTimeout = sub.GetDuration("conn-timeout") * time.Second
	options.ReadTimeout = sub.GetDuration("read-timeout") * time.Second
	options.PoolTimeout = sub.GetDuration("pool-timeout") * time.Second
	options.IdleTimeout = sub.GetDuration("idle-timeout") * time.Second
	options.MaxRetries = sub.GetInt("retry")
	options.PoolSize = sub.GetInt("pool-size")
	options.MinIdleConns = sub.GetInt("min-idle-conn")
	return &options
}
