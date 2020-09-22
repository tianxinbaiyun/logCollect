package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/tianxinbaiyun/logCollect/logAgent/tailf"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/coreos/etcd/mvcc/mvccpb"
	etcd "github.com/etcd-io/etcd/clientv3"
)

// etcd客户端对象
type EtcdClient struct {
	client *etcd.Client
	// 存储日志收集的key
	collectKeys []string
}

var (
	etcdClient *EtcdClient
)

// 初始化etcd
func initEtcd(etcdAddress []string, collectKey string) (err error) {
	client, err := etcd.New(etcd.Config{
		Endpoints:   etcdAddress,
		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		return
	}
	etcdClient = &EtcdClient{
		client: client,
	}

	if strings.HasSuffix(collectKey, "/") == false {
		collectKey = fmt.Sprintf("%s/", collectKey)
	}

	// 通过本地ip和配置文件中的前缀值获取etcd中真正的数据值
	for _, ip := range localIpArray {
		etcdKey := fmt.Sprintf("%s%s", collectKey, ip)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		resp, err := etcdClient.client.Get(ctx, etcdKey)
		cancel()

		if err != nil {
			logs.Warn("get key: %s from etcd failed, err: %s", etcdKey, err)
			continue
		}
		etcdClient.collectKeys = append(etcdClient.collectKeys, etcdKey)

		for _, v := range resp.Kvs {
			if string(v.Key) == etcdKey {
				err = json.Unmarshal(v.Value, &agentConfig.Collects)
				if err != nil {
					logs.Warn("json Unmarshal key: %s failed, err: %s", v.Key, err)
					continue
				}
				// 设置生效的ip地址
				agentConfig.Ip = ip
			}
		}
		logs.Debug("log agent collect is: %v", agentConfig.Collects)
	}
	// 初始化etcd key监控
	initEtcdWatch()
	return
}

// 初始化etcd key监控
func initEtcdWatch() {
	for _, key := range etcdClient.collectKeys {
		go etcdWatch(key)
	}
}

// 	etcd key监控处理
func etcdWatch(key string) {
	logs.Debug("start watch key: %s", key)
	for true {
		rech := etcdClient.client.Watch(context.Background(), key)
		var colConfig []tailf.Collect
		var getConfStatus = true
		for wresp := range rech {
			for _, ev := range wresp.Events {
				// key 删除
				if int(ev.Type) == int(mvccpb.DELETE) {
					logs.Warn("key [%s] is deleted", key)
					continue
				}
				// key 更新
				if int(ev.Type) == int(mvccpb.PUT) && string(ev.Kv.Key) == key {
					err := json.Unmarshal(ev.Kv.Value, &colConfig)
					if err != nil {
						logs.Error("key [%s], Unmarshal[%s], err:%s", key, err)
						getConfStatus = false
						continue
					}
				}
				logs.Debug("get etcd config, %s %q : %q", ev.Type, ev.Kv.Key, ev.Kv.Value)
			}
			if getConfStatus {
				break
			}
		}
		logs.Info("Update tailf config")
		// 更新tailf任务
		err := tailf.UpdateTailfTask(colConfig)
		if err != nil {
			logs.Error("Update tailf task failed, connect: %s, err: %s", colConfig, err)
		}
	}
}
