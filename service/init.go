package service

import (
	"encoding/json"
	"io/ioutil"
	"mj/model"
	"os/exec"

	"sync"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/go-redis/redis/v8"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type config struct {
	Mysql string `yaml:"mysql"`
}

var (
	err                  error
	Config               = config{}
	DB                   = new(gorm.DB)
	ImagineRequestModel  = model.Imagine{}
	BlendRequestModel    = model.Blend{}
	DescribeRequestModel = model.Describe{}
	Redis                = new(redis.Client)
	MJListRunMap         sync.Map       // 看队列中是否有任务了
	SDListRunMap         sync.Map       //看队列中是否有任务了
	MjNumberMap          sync.Map       //判断每个账号下剩余并发数
	ChannelStateMap      sync.Map       //判断哪一个频道现在没人
	MjInfoMod            []model.MjInfo //方便快速查看账号下的频道
	DiscribeNumber       sync.Map       //discribe的时候看返回了几次
	DisChanTime          sync.Map       //定时
	DisChanUser          sync.Map       //每个频道 是哪个用户Id在用
	MQ                   rocketmq.Producer
	MJMQ                 rocketmq.Producer
	SDMQ                 = make(chan model.SDPromptConfig, 1000)
)

// 初始化两个并发Map
func initMap() {

	DB.Model(&model.MjAccount{}).Find(&MjInfoMod)
	//MjNumberMap[bot token]=并发数
	//ChannelStateMap[bot token||||channel_id]=false  没人用
	for i, j := range MjInfoMod {
		DB.Model(&model.RequestInfo{}).Where("mj_id=?", j.ID).Find(&MjInfoMod[i].RequestInfo)
		MjNumberMap.Store(j.BotToken, len(MjInfoMod[i].RequestInfo)) //这个账号下有多少频道就是多少个并发
		for _, jj := range MjInfoMod[i].RequestInfo {
			key := jj.ChannelID
			ChannelStateMap.Store(key, 1) //初始化都没人用
			DisChanTime.Store(key, false)
		}

	}

}

// 解析配置文件
func getConfig() {
	by, err := ioutil.ReadFile("config/yaml.yaml")
	if err != nil {
		panic("read config error")
	}
	err = yaml.Unmarshal(by, &Config)
	if err != nil {
		panic("yaml Unmarshal error")
	}
}

// 连接数据库
func setSqlDB() {

	DB, err = gorm.Open(mysql.Open(Config.Mysql), &gorm.Config{
		PrepareStmt: true, // 启用 prepared statements 缓存，可以提高查询性能
	})
	if err != nil {
		panic("failed to connect database")
	}
	sqlDB, err := DB.DB()
	if err != nil {
		panic("failed to get DB instance")
	}
	sqlDB.SetMaxIdleConns(10)   // 最大闲置连接数
	sqlDB.SetMaxOpenConns(100)  // 最大打开连接数
	sqlDB.SetConnMaxLifetime(0) // 连接的最大存活时间（0 表示不限制）
}

// 初始化文生图模块
func imagine() {

	by, err := ioutil.ReadFile("config/config/imagine.json")
	if err != nil {
		panic("read config/config/imagine.json error")
	}
	json.Unmarshal(by, &ImagineRequestModel)
	for i, _ := range ImagineRequestModel.Data.Options {
		ImagineRequestModel.Data.Options[i].Value = ""
	}
}
func describe() {
	by, _ := ioutil.ReadFile("config/config/describe.json")
	if err != nil {
		panic("read config/config/describe.json error")
	}
	json.Unmarshal(by, &DescribeRequestModel)
}
func redisSet() {
	Redis = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis 服务器地址
		Password: "",               // 设置密码
		DB:       0,                // 使用默认的数据库
	})
}
func clearSDMQ() {
	// 设置 NameServer 地址
	nameServer := "*.*.*.*:9876"
	// 设置消费者组
	consumerGroup := "resp"
	// 设置主题名称
	topicName := "threezto-test"

	// 构建 mqadmin 命令来设置消费进度到最新位置，间接清空消息
	cmd := exec.Command("./mqadmin", "resetOffsetByTime",
		"-n", nameServer,
		"-g", consumerGroup,
		"-t", topicName,
		"-s", "now",
	)

	// 设置工作目录为 `mqadmin` 所在的目录
	cmd.Dir = "/rocketmq/rocketmq-all-5.1.4-bin-release/bin"

	// 运行命令并捕获输出
	_, err := cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}

}
func rocketMq() {
	//清空MQ
	clearSDMQ()
	MQ, err = rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{"*.*.*.*:*"})),
		producer.WithRetry(2), //指定重试次数
	)
	if err != nil {
		panic(err)
	}
	if err = MQ.Start(); err != nil {
		panic("启动producer失败")
	}

}

func init() {

	getConfig()

	setSqlDB()

	imagine()

	describe()

	redisSet()

	initMap()

	// rocketMq()
	for _, j := range MjInfoMod {
		go message(j.BotToken)
	}

}
