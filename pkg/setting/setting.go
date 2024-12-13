package setting

import (
	"time"

	"log"

	"github.com/go-ini/ini"
)

/*
* 相关读写超时的设置需要注意，默认读取得到的单位是ns，需要进行转换，过小的超时时间将导致无法正常连接
 */

var (
	Cfg *ini.File

	AppSetting      = &AppConfig{}
	ServerSetting   = &ServerConfig{}
	DatabaseSetting = &DatabaseConfig{}
	RedisSetting    = &RedisConfig{}
)

// 全局的格外配置
type AppConfig struct {
	PageSize        int
	JwtSecret       string // jwt的密钥
	RuntimeRootPath string

	PrefixUrl      string
	ImageSavePath  string
	ImageMaxSize   int // MB 为单位
	ImageAllowExts []string

	ExportSavePath string // 文件导出存储的路径
	QrCodeSavePath string // 二维码的保存路径
	FontSavePath   string // 二维码的保存路径

	LogSavePath string // log日志保存的位置
	LogSaveName string
	LogFileExt  string
	TimeFormat  string
}

// 服务器的相关配置
type ServerConfig struct {
	RunMode      string
	HTTPPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// 数据库的相关配置
type DatabaseConfig struct {
	Type        string
	DBName      string
	User        string
	Password    string
	Host        string
	TablePrefix string
}

// Redis相关的配置
type RedisConfig struct {
	Addr           string
	DB             int
	Password       string
	MaxIdleConns   int
	MinIdleConns   int
	MaxActiveConns int
	IdleTimeout    time.Duration
	ConnectTimeout time.Duration
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
}

// 初始化加载所有配置信息
func Setup() {
	var err error
	Cfg, err = ini.Load("./conf/app.ini")
	if err != nil {
		log.Fatalf("Fail to parse [conf/app.ini]: %v ", err)
	}

	// 加载本服务器的配置
	err = Cfg.Section("app").MapTo(AppSetting)
	if err != nil {
		log.Fatalf("Fail to get section app: %v", err)
	}
	AppSetting.ImageMaxSize = AppSetting.ImageMaxSize * 1024 * 1024
	log.Printf("Success load ServerSetting: %+v", AppSetting)

	// 加载本服务器的配置
	err = Cfg.Section("server").MapTo(ServerSetting)
	if err != nil {
		log.Fatalf("Fail to get section server: %v", err)
	}
	ServerSetting.ReadTimeout = ServerSetting.ReadTimeout * time.Second
	ServerSetting.WriteTimeout = ServerSetting.ReadTimeout * time.Second
	log.Printf("Success load ServerSetting: %+v", ServerSetting)

	// 加载数据库配置
	err = Cfg.Section("database").MapTo(DatabaseSetting)
	if err != nil {
		log.Fatalf("Fail to get section database: %v", err)
	}
	log.Printf("Success load DatabaseSetting: %+v", DatabaseSetting)

	// 加载redis的配置
	err = Cfg.Section("redis").MapTo(RedisSetting)
	if err != nil {
		log.Fatalf("Fail to get section redis: %v", err)
	}
	RedisSetting.IdleTimeout = RedisSetting.IdleTimeout * time.Second
	RedisSetting.ReadTimeout = RedisSetting.ReadTimeout * time.Second
	RedisSetting.WriteTimeout = RedisSetting.WriteTimeout * time.Second
	RedisSetting.ConnectTimeout = RedisSetting.ConnectTimeout * time.Second
	log.Printf("Success load RedisSetting: %+v", RedisSetting)

}
