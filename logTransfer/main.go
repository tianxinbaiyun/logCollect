package main

import (
	"github.com/tianxinbaiyun/logCollect/logTransfer/elasticsearch"
	"github.com/tianxinbaiyun/logCollect/logTransfer/kafka"
	"github.com/tianxinbaiyun/logCollect/logTransfer/transfer"

	"github.com/astaxie/beego/logs"
)

var (
	configType = "ini"
	configPath string
)

func main() {
	var err error
	// 加载配置文件
	transferConfig, err := transfer.LoadConfig()
	if err != nil {
		logs.Error("load config failed, err:%s", err)
		return
	}
	logs.Debug("load log transfer success")

	// 初始化日志
	err = transfer.InitTransgerLog()
	if err != nil {
		logs.Error("init log failed, err:%s", err)
		return
	}
	logs.Debug("init transfer log success")

	// 初始化etcd
	err = transfer.InitEtcd(transferConfig.EtcdAddress)
	if err != nil {
		logs.Error("init etcd failed, err:%s", err)
		return
	}
	logs.Debug("init transfer etcd success")

	// 初始化elasticsearch
	err = elasticsearch.InitElastic(transferConfig.EsAddress, transferConfig.Chansize)
	if err != nil {
		logs.Error("init elasticsearch failed, err: %s", err)
		return
	}
	logs.Debug("int transger elasticsearch success")

	// 初始化kafka
	err = kafka.InitKafka(transferConfig.KafkaAddress, transferConfig.Topics)
	if err != nil {
		logs.Error("init kafka failed, err: %s", err)
		return
	}
	logs.Debug("init transger kafka success")

	// 启动服务
	err = transfer.ServerRun()
	if err != nil {
		logs.Error("server run failed, err: %s", err)
		return
	}
	logs.Info("Log Transger exited")

}
