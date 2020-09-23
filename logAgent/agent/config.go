package agent

import (
	"errors"
	"github.com/astaxie/beego/logs"
	"github.com/tianxinbaiyun/logCollect/logAgent/tailf"
	"os"
	"strings"

	"github.com/astaxie/beego/config"
)

const configType = "ini"

// 存储logAgent配置信息
type Config struct {
	LogLevel     string
	LogPath      string
	Chansize     int
	KafkaAddress []string
	EtcdAddress  []string
	CollectKey   string
	Collects     []tailf.Collect
	Ip           string
}

var (
	// 配置信息对象
	agentConfig *Config
)

// 加载配置信息
func LoadConfig() (agentConfig *Config, err error) {
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
	agentConfig = &Config{}

	// 获取基础配置
	err = getAgentConfig(conf, agentConfig)
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

func getAgentConfig(conf config.Configer, agentConfig *Config) (err error) {
	// 获取日志级别
	logLevel := conf.String("base::log_level")
	if len(logLevel) == 0 {
		logLevel = "debug"
	}
	agentConfig.LogLevel = logLevel

	// 获取日志路径
	logPath := conf.String("base::log_path")
	if len(logPath) == 0 {
		logPath = "/Users/aery/Data/code/Go/go_old/logs/logagent.log"
	}
	agentConfig.LogPath = logPath

	// 日志收集开启chan大小
	chanSize, chanStatus := conf.Int("base::queue_size")
	if chanStatus != nil {
		chanSize = 200
	}
	agentConfig.Chansize = chanSize

	// etcd 地址
	etcdAddress := conf.String("etcd::etcd_address")
	if len(etcdAddress) == 0 {
		err = errors.New("Agent config etcd address error")
		return
	}
	agentConfig.EtcdAddress = strings.Split(etcdAddress, ",")

	// kafka 地址
	kafkaAddress := conf.String("kafka::kafka_address")
	if len(kafkaAddress) == 0 {
		err = errors.New("Agent config kafka address error")
		return
	}
	agentConfig.KafkaAddress = strings.Split(kafkaAddress, ",")

	// 获取日志收集前缀key
	collectKey := conf.String("collect::collectKey")
	if len(collectKey) == 0 {
		err = errors.New("Agent config collectKey error")
		return
	}
	agentConfig.CollectKey = collectKey

	return
}
