package transfer

import (
	"errors"
	"github.com/astaxie/beego/logs"
	"os"
	"strings"

	"github.com/astaxie/beego/config"
)

const configType = "ini"

// 配置文件结构体
type Config struct {
	LogLevel     string
	LogPath      string
	Chansize     int
	KafkaAddress []string
	EtcdAddress  []string
	EsAddress    []string
	Topics       []string
}

// 加载配置文件配置
func LoadConfig() (transferConfig *Config, err error) {
	// 获取配置文件路径
	configPath, err := GetConfigPath()
	if err != nil {
		logs.Error("get config file failed, err: %s", err)
		return
	}
	logs.Debug("get config file success, file: %s", configPath)

	filePath := GetExecpath() + "/" + configPath

	//配置文件不存在，从配置文件指定的目录找
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		filePath = GetCurrentPath() + "/../" + configPath
	}
	conf, err := config.NewConfig(configType, filePath)
	if err != nil {
		return
	}

	transferConfig = &Config{}

	// 获取基础配置
	err = getTransgerConfig(conf, transferConfig)
	if err != nil {
		return
	}
	return
}

// 根据传参的方式获取配置文件路径
func GetConfigPath() (configPath string, err error) {
	cmdArgs := os.Args
	if len(cmdArgs) < 2 {
		configPath = "config/app.ini"
		//err = fmt.Errorf("USAGE: %v  <agent config file>, go to start the log agent.", cmdArgs[0])
		//return
	} else {
		configPath = cmdArgs[1]
	}

	return
}

func getTransgerConfig(conf config.Configer, transferConfig *Config) (err error) {
	// 获取日志级别
	logLevel := conf.String("base::log_level")
	if len(logLevel) == 0 {
		logLevel = "debug"
	}
	transferConfig.LogLevel = logLevel

	// 获取日志路径
	logPath := conf.String("base::log_path")
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		logPath = GetExecpath() + "/logs/logtransfer.log"
		//配置文件不存在，从配置文件指定的目录找
		if _, err := os.Stat(logPath); os.IsNotExist(err) {
			logPath = GetCurrentPath() + "/../logs/logtransfer.log"
		}
	}
	transferConfig.LogPath = logPath

	// 日志收集开启chan大小
	chanSize, chanStatus := conf.Int("base::queue_size")
	if chanStatus != nil {
		chanSize = 200
	}
	transferConfig.Chansize = chanSize

	// etcd 地址
	etcdAddress := conf.String("etcd::etcd_address")
	if len(etcdAddress) == 0 {
		err = errors.New("Transger config etcd address error")
		return
	}
	transferConfig.EtcdAddress = strings.Split(etcdAddress, ",")

	// kafka 地址
	kafkaAddress := conf.String("kafka::kafka_address")
	if len(kafkaAddress) == 0 {
		err = errors.New("Transger config kafka address error")
		return
	}
	transferConfig.KafkaAddress = strings.Split(kafkaAddress, ",")

	// elasticsearch 地址
	esAddress := conf.String("elasticsearch::es_address")
	if len(kafkaAddress) == 0 {
		err = errors.New("Transger config elasticsearch address error")
		return
	}
	transferConfig.EsAddress = strings.Split(esAddress, ",")

	return
}
