package main

import (
	"github.com/tianxinbaiyun/logCollect/logAgent/agent"
	"github.com/tianxinbaiyun/logCollect/logAgent/kafka"
	"github.com/tianxinbaiyun/logCollect/logAgent/tailf"

	"github.com/astaxie/beego/logs"
)

var (
	configType = "ini"
	configPath string
)

func main() {
	var err error
	// 加载配置文件
	config, err := agent.LoadConfig()
	if err != nil {
		logs.Error("Start logAgent [init loadConfig] failed, err: %s", err)
		return
	}
	logs.Debug("load Agent [config] success")

	// 初始化日志
	err = agent.InitAgentLog()
	if err != nil {
		logs.Error("Start logAgent [init agentLog] failed, err: %s", err)
		return
	}
	logs.Debug("Init Agent [log] success")

	// 初始化Etcd
	err = agent.InitEtcd(config.EtcdAddress, config.CollectKey)
	if err != nil {
		logs.Error("Start logAgent [init etcd] failed, err:", err)
		return
	}
	logs.Debug("Init Agent [etcd] success")

	// 初始化tailf
	err = tailf.InitTailf(config.Collects, config.Chansize, config.Ip)
	if err != nil {
		logs.Error("Start logAgent [init tailf] failed, err:", err)
		return
	}
	logs.Debug("Init Agent [tailf] success")

	// 初始化kafka
	err = kafka.InitKafka(config.KafkaAddress)
	if err != nil {
		logs.Error("Start logAgent [init kafka] failed, err:", err)
		return
	}
	logs.Debug("Init Agent [kafka] success")
	// 启动logagent服务
	err = agent.ServerRun()
	if err != nil {
		logs.Error("Start logAgent [init serverRun] failed, err:", err)
		return
	}
	logs.Info("Log Agent exit")
}
