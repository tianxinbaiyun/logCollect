package manager

import (
	"errors"
	"github.com/tianxinbaiyun/logCollect/logManager/models"
	_ "github.com/tianxinbaiyun/logCollect/logManager/routers"
	"strings"

	"github.com/astaxie/beego"
)

// 初始化mysq
func InitMysql() (err error) {
	mysqlInfo := models.MysqlInfo{
		Address:  beego.AppConfig.String("mysqladdr"),
		Port:     beego.AppConfig.String("mysqlport"),
		DbName:   beego.AppConfig.String("mysqldbname"),
		User:     beego.AppConfig.String("mysqluser"),
		Password: beego.AppConfig.String("mysqlpwd"),
	}
	err = models.InitMySql(mysqlInfo)
	if err != nil {
		return
	}
	return
}

// 初始化etcd
func InitEtcd() (err error) {
	address := beego.AppConfig.String("etcdaddr")
	if len(address) == 0 {
		errors.New("etcd addresss config error")
		return
	}
	etcdAddress := strings.Split(address, ",")
	err = models.InitEtcd(etcdAddress)
	if err != nil {
		return
	}
	return
}
