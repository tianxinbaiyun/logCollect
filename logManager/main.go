package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/tianxinbaiyun/logCollect/logManager/manager"
)

func main() {
	var err error

	// 初始化log
	err = manager.InitLog()
	if err != nil {
		logs.Error("init log manager failed, err: %s", err)
	}
	// 初始化mysql
	err = manager.InitMysql()
	if err != nil {
		logs.Error("init mysql failed, err: %s", err)
	}
	// 初始化etcd
	err = manager.InitEtcd()
	if err != nil {
		logs.Error("init etcd failed, err: %s", err)
	}

	beego.Run()
}
