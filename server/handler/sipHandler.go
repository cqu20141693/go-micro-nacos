package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/cqu20141693/sip-server/client"
	"github.com/cqu20141693/sip-server/common"
	"github.com/cqu20141693/sip-server/db"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go-micro.dev/v4/logger"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type sipApi interface {
	health(c *gin.Context)
	GetCameraInfo(c *gin.Context)
	selfHealth(c *gin.Context)
}

type SipService struct {
	sipClient *client.SipClient
}

func (s *SipService) GetCameraInfo(c *gin.Context) {
	panic("implement me")
}

func NewSipService(sipClient *client.SipClient) *SipService {
	return &SipService{sipClient: sipClient}
}

func (s *SipService) InitRouteMapper(router *gin.Engine) {
	router.POST(common.HealthPath, s.health)
	router.POST(common.GetCameraInfoPath, s.GetCameraInfo)
	router.POST(common.SelfHealth, s.selfHealth)
	router.POST(common.CmdTestPath, s.cmdTest)
	router.POST(common.PagePath, s.cmdPaged)
	router.POST(common.RedisPath, s.redisTest)
}

// PingExample godoc
// @Summary health example
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} health
// @Router /health [post]
func (s *SipService) health(c *gin.Context) {
	c.JSON(200, common.ResultUtils.Success(map[string]string{"status": "up"}))
}

func (s *SipService) selfHealth(c *gin.Context) {
	health := s.sipClient.Health()
	c.JSON(200, health)
}

type Command struct {
	Id          int32     `gorm:"primarykey"`
	GmtCreate   time.Time `gorm:"column:gmt_create"`
	GmtModified time.Time `gorm:"column:gmt_modified"`
	GroupKey    string    `gorm:"column:group_key"`
	Sn          string
	CmdTag      string `gorm:"column:cmd_tag"`
	Category    string
	Status      string
	StatusLevel int32 `gorm:"column:status_level"`
	Type        string
	Data        string
	Context     datatypes.JSON
}

func (c Command) TableName() string {
	return "tb_command"
}

func (s *SipService) cmdTest(c *gin.Context) {
	cmd := Command{}
	tx := db.MysqlDB.Model(&Command{}).First(&cmd)
	if tx.Error != nil {
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			c.JSON(200, common.ResultUtils.Success(tx.Error.Error()))
			return
		}
	}
	c.JSON(200, common.ResultUtils.Success(cmd))
}

func (s *SipService) cmdPaged(c *gin.Context) {
	var p db.Page
	err := c.ShouldBindJSON(&p)
	if err != nil {
		logger.Info("binding failed")
		c.JSON(http.StatusOK, common.ResultUtils.Fail("0000", err.Error()))
		return
	}
	if p.Page == 0 {
		p.Page = 1
	}
	if p.PageSize == 0 {
		p.PageSize = 10
	}
	var orderSql string
	if p.Desc {
		orderSql = "id desc"
	} else {
		orderSql = "id asc"
	}
	// total
	var total int64
	tx := db.MysqlDB.Model(&Command{}).Count(&total)
	if tx.Error != nil {
		c.JSON(200, common.ResultUtils.Success(tx.Error.Error()))
		return
	}
	commands := make([]Command, 0)
	if err := db.MysqlDB.Model(&Command{}).Limit(p.PageSize).Offset((p.Page - 1) * p.PageSize).Order(orderSql).Find(&commands).Error; err != nil {
		c.JSON(200, common.ResultUtils.Success(tx.Error.Error()))
		return
	}

	c.JSON(200, common.ResultUtils.Success(db.NewPageInfo(total, p.Page, p.PageSize, commands)))
}

func (s *SipService) redisTest(c *gin.Context) {
	redisBasicTest()
	testRedisNil()

	testPubSub()

	testPipeline()

	testZSet()

	testGeo()
}

/*
geoadd：增加某个位置的坐标。
geopos：获取某个位置的坐标。
geohash：获取某个位置的geohash值。
geodist：获取两个位置的距离。
georadius：根据给定位置坐标获取指定范围内的位置集合。
georadiusbymember：根据给定位置获取指定范围内的位置集合。
*/
func testGeo() {
	key := "geo_key"
	count, err := db.RedisDB.GeoAdd(context.Background(), key, &redis.GeoLocation{
		Name:      "天府广场",
		Longitude: 104.072833,
		Latitude:  30.663422,
	}, &redis.GeoLocation{
		Name:      "四川大剧院",
		Longitude: 104.074378,
		Latitude:  30.664804,
	}, &redis.GeoLocation{
		Name:      "新华文轩",
		Longitude: 104.070084,
		Latitude:  30.664649,
	}, &redis.GeoLocation{
		Name:      "手工茶",
		Longitude: 104.072402,
		Latitude:  30.664121,
	}, &redis.GeoLocation{
		Name:      "宽窄巷子",
		Longitude: 104.059826,
		Latitude:  30.669883,
	}, &redis.GeoLocation{
		Name:      "奶茶",
		Longitude: 104.06085,
		Latitude:  30.670054,
	}, &redis.GeoLocation{
		Name:      "钓鱼台",
		Longitude: 104.058424,
		Latitude:  30.670737,
	}).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("GeoAdd:", count)
	//GeoPos
	//获取某个位置的坐标。
	resPos, err := db.RedisDB.GeoPos(context.Background(), key, "天府广场", "宽窄巷子").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("GeoPos:")
	for _, pos := range resPos {
		fmt.Println(pos)
	}
	// GeoHash
	// 获取某个位置的geohash值。
	resHash, err := db.RedisDB.GeoHash(context.Background(), key, "天府广场").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("GeoHash:", resHash)

	//GeoDist
	//获取两个位置的距离
	resDist, err := db.RedisDB.GeoDist(context.Background(), key, "天府广场", "宽窄巷子", "m").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("GeoDist:", resDist)
	// GeoRadius
	// 获取某个指定经纬度附近的位置
	resRadiu, err := db.RedisDB.GeoRadius(context.Background(), key, 104.072833, 30.663422, &redis.GeoRadiusQuery{
		Radius:      800,   //radius表示范围距离，
		Unit:        "m",   //距离单位是 m|km|ft|mi
		WithCoord:   true,  //传入WITHCOORD参数，则返回结果会带上匹配位置的经纬度
		WithDist:    true,  //传入WITHDIST参数，则返回结果会带上匹配位置与给定地理位置的距离。
		WithGeoHash: true,  //传入WITHHASH参数，则返回结果会带上匹配位置的hash值。
		Count:       4,     //入COUNT参数，可以返回指定数量的结果。
		Sort:        "ASC", //默认结果是未排序的，传入ASC为从近到远排序，传入DESC为从远到近排序。
	}).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("GeoRadius:")
	for _, location := range resRadiu {
		fmt.Println(location)
	}

	resRadiusByMember, err := db.RedisDB.GeoRadiusByMember(context.Background(), key, "天府广场", &redis.GeoRadiusQuery{
		Radius:      800,
		Unit:        "m",
		WithCoord:   true,
		WithDist:    true,
		WithGeoHash: true,
		Count:       4,
		Sort:        "ASC",
	}).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("GeoRadiusByMember:")
	for _, location := range resRadiusByMember {
		fmt.Println(location)
	}

}

func testZSet() {
	client := db.RedisDB
	/*
		type Z struct {
			Score  float64 // 分数
			Member interface{} // 元素名
		}
	*/
	// 添加一个集合元素到集合中， 这个元素的分数是2.5，元素名是tizi
	key := "z-set-key"
	err := client.ZAdd(context.Background(), key, &redis.Z{Score: 2.5, Member: "tizi"}).Err()
	if err != nil {
		panic(err)
	}

	// 查询集合元素tizi的分数
	score, _ := client.ZScore(context.Background(), key, "tizi").Result()
	fmt.Println(score)

	//返回集合元素个数
	size, err := client.ZCard(context.Background(), key).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(size)

	// 返回： 1<=分数<=5 的元素个数, 注意："1", "5"两个参数是字符串
	size, err = client.ZCount(context.Background(), key, "1", "5").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(size)

	// 返回： 1<分数<=5 的元素个数
	// 说明：默认第二，第三个参数是大于等于和小于等于的关系。
	// 如果加上（ 则表示大于或者小于，相当于去掉了等于关系。
	size, err = client.ZCount(context.Background(), key, "(2.5", "5").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(size)

	//增加元素的分数

	// 给元素5，加上2分
	client.ZIncrBy(context.Background(), key, 2, "5")

	// ZRange
	//返回集合中某个索引范围的元素，根据分数从小到大排序
	// 返回从0到-1位置的集合元素， 元素按分数从小到大排序
	// 0到-1代表则返回全部数据
	vals, err := client.ZRange(context.Background(), key, 0, -1).Result()
	if err != nil {
		panic(err)
	}

	for _, val := range vals {
		fmt.Println(val)
	}
	//ZRevRange
	//用法跟ZRange一样，区别是ZRevRange的结果是按分数从大到小排序。

	// ZRangeByScore
	//根据分数范围返回集合元素，元素根据分数从小到大排序，支持分页。
	// 初始化查询条件， Offset和Count用于分页
	op := redis.ZRangeBy{
		Min:    "2",  // 最小分数
		Max:    "10", // 最大分数
		Offset: 0,    // 类似sql的limit, 表示开始偏移量
		Count:  5,    // 一次返回多少数据
	}

	vals, err = client.ZRangeByScore(context.Background(), key, &op).Result()
	if err != nil {
		panic(err)
	}
	for _, val := range vals {
		fmt.Println(val)
	}
	// .ZRevRangeByScore
	//用法类似ZRangeByScore，区别是元素根据分数从大到小排序。

	//ZRangeByScoreWithScores
	//用法跟ZRangeByScore一样，区别是除了返回集合元素，同时也返回元素对应的分数
	// 初始化查询条件， Offset和Count用于分页
	op = redis.ZRangeBy{
		Min:    "2",  // 最小分数
		Max:    "10", // 最大分数
		Offset: 0,    // 类似sql的limit, 表示开始偏移量
		Count:  5,    // 一次返回多少数据
	}

	valsWithScore, err := client.ZRangeByScoreWithScores(context.Background(), key, &op).Result()
	if err != nil {
		panic(err)
	}

	for _, val := range valsWithScore {
		fmt.Println(val.Member, " -> ", val.Score) // 集合元素

	}

	//ZRank
	//根据元素名，查询集合元素在集合中的排名，从0开始算，集合元素按分数从小到大排序
	rk, _ := client.ZRank(context.Background(), key, "tizi").Result()
	fmt.Println(rk)
	//ZRevRank的作用跟ZRank一样，区别是ZRevRank是按分数从大到小排序。
	rk, _ = client.ZRevRank(context.Background(), key, "tizi").Result()
	fmt.Println(rk)
	// 删除集合中的元素tizi
	client.ZRem(context.Background(), key, "tizi")

	// 删除集合中的元素tizi和xiaoli
	// 支持一次删除多个元素
	client.ZRem(context.Background(), key, "tizi", "xiaoli")
	// ZRemRangeByRank
	//根据索引范围删除元素
	// 集合元素按分数排序，从最低分到高分，删除第0个元素到第5个元素。
	// 这里相当于删除最低分的几个元素
	client.ZRemRangeByRank(context.Background(), key, 0, 5)

	// 位置参数写成负数，代表从高分开始删除。
	// 这个例子，删除最高分数的两个元素，-1代表最高分数的位置，-2第二高分，以此类推。
	client.ZRemRangeByRank(context.Background(), key, -1, -2)

	//ZRemRangeByScore
	//根据分数范围删除元素
	// 删除范围： 2<=分数<=5 的元素
	client.ZRemRangeByScore(context.Background(), key, "2", "5")

	// 删除范围： 2<=分数<5 的元素
	client.ZRemRangeByScore(context.Background(), key, "2", "(5")
}

func redisBasicTest() {
	// 第三个参数代表key的过期时间，0代表不会过期。
	client := db.RedisDB
	err := client.Set(context.Background(), "key", "value", time.Minute).Err()
	if err != nil {
		panic(err)
	}
	// Result函数返回两个值，第一个是key的值，第二个是错误信息
	val, err := client.Get(context.Background(), "key").Result()
	// 判断查询是否出错
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)
	// Result函数返回两个值，第一个是key的值，第二个是错误信息
	oldVal, err := client.GetSet(context.Background(), "key", "new value").Result()

	if err != nil {
		panic(err)
	}
	// 打印key的旧值
	fmt.Println("key", oldVal)
	// 第三个参数代表key的过期时间，0代表不会过期。
	result, err := client.SetNX(context.Background(), "setNX", "new", time.Minute).Result()
	if err != nil {
		return
	}
	fmt.Println("setNX=", result)

	_, err = client.MSet(context.Background(), "key1", "value1", "key2", "value2", "key3", "value3").Result()
	if err != nil {
		panic(err)
	}

	vals, err := client.MGet(context.Background(), "key1", "key2", "key3").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(vals)

	// Incr函数每次加一
	intVal, err := client.Incr(context.Background(), "int-key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("最新值", intVal)

	// IncrBy函数，可以指定每次递增多少
	intVal, err = client.IncrBy(context.Background(), "int-key", 2).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("最新值", intVal)

	// IncrByFloat函数，可以指定每次递增多少，跟IncrBy的区别是累加的是浮点数
	floatVal, err := client.IncrByFloat(context.Background(), "int-key", 2).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("最新值", floatVal)

	// Decr函数每次减一
	intVal, err = client.Decr(context.Background(), "int-key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("最新值", intVal)

	// DecrBy函数，可以指定每次递减多少
	intVal, err = client.DecrBy(context.Background(), "int-key", 2).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("最新值", intVal)

	client.Expire(context.Background(), "key", 3)

	// 删除key
	client.Del(context.Background(), "key")

	// 删除多个key, Del函数支持删除多个key
	_, err = client.Del(context.Background(), "key1", "key2", "key3").Result()
	if err != nil {
		panic(err)
	}
}

func testPipeline() {
	pipe := db.RedisDB.Pipeline()

	incr := pipe.Incr(context.Background(), "pipeline_counter")
	pipe.Expire(context.Background(), "pipeline_counter", time.Minute)
	pipe.Get(context.Background(), "pipeline_counter")
	cmds, err := pipe.Exec(context.Background())
	if err != nil {
		panic(err)
	}
	for _, cmd := range cmds {
		switch cmd.(type) {
		case *redis.IntCmd:
			fmt.Println("incr ", cmd.(*redis.IntCmd).Val())
		case *redis.StringCmd:
			fmt.Println("get ", cmd.(*redis.StringCmd).Val())
		case *redis.BoolCmd:
			fmt.Println("expire ", cmd.(*redis.BoolCmd).Val())
		}
	}

	// The value is available only after Exec is called.
	fmt.Println(incr.Val())

	cmds, err = db.RedisDB.Pipelined(context.Background(), func(pipe redis.Pipeliner) error {
		incr = pipe.Incr(context.Background(), "pipeline_counter")
		pipe.Expire(context.Background(), "pipeline_counter", time.Minute)
		return nil
	})
	if err != nil {
		panic(err)
	}
	for _, cmd := range cmds {
		switch cmd.(type) {
		case *redis.IntCmd:
			fmt.Println("incr ", cmd.(*redis.IntCmd).Val())
		case *redis.StringCmd:
			fmt.Println("get ", cmd.(*redis.StringCmd).Val())
		case *redis.BoolCmd:
			fmt.Println("expire ", cmd.(*redis.BoolCmd).Val())
		}
	}

}

func testPubSub() {
	channel := "my-channel"
	go func(counter int) {
		for i := 0; i < counter; i++ {
			err := db.RedisDB.Publish(context.Background(), channel, "payload").Err()
			if err != nil {
				panic(err)
			}
			time.Sleep(time.Millisecond)
		}
	}(10)
	// There is no error because go-redis automatically reconnects on error.
	pubsub := db.RedisDB.Subscribe(context.Background(), channel)
	// Close the subscription when we are done.

	//To receive a message:
	go func() {
		defer pubsub.Close()
		for {
			msg, err := pubsub.ReceiveMessage(context.Background())
			if err != nil {
				panic(err)
			}

			fmt.Println("for", msg.Channel, msg.Payload)
		}
	}()
	//simple way
	go func() {
		ch := pubsub.Channel()
		for msg := range ch {
			fmt.Println("simple", msg.Channel, msg.Payload)
		}
	}()

	subscribe := db.RedisDB.PSubscribe(context.Background(), "my_*")
	//simple way
	go func() {
		ch := subscribe.Channel()
		for msg := range ch {
			fmt.Println("simple", msg.Channel, msg.Payload)
		}
	}()
}

/*
	GET is not the only command that returns nil reply, for example, BLPOP and ZSCORE can also return redis.Nil
*/
func testRedisNil() {
	val, err := db.RedisDB.Get(context.Background(), "key").Result()
	switch {
	case err == redis.Nil:
		fmt.Println("key does not exist")
	case err != nil:
		fmt.Println("Get failed", err)
	case val == "":
		fmt.Println("value is empty")
	}

	result, err1 := db.RedisDB.BLPop(context.Background(), time.Second, "key").Result()
	switch {
	case err1 == redis.Nil:
		fmt.Println("key does not exist")
	case err1 != nil:
		fmt.Println("Get failed", err1)
	case len(result) == 0:
		fmt.Println("value is empty")
	}
	score, err2 := db.RedisDB.ZScore(context.Background(), "key", "member").Result()
	switch {
	case err2 == redis.Nil:
		fmt.Println("key does not exist")
	case err2 != nil:
		fmt.Println("Get failed", err2)
	case score == 0:
		fmt.Println("value is empty")
	}
}
