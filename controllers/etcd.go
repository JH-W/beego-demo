package controllers

import (
	"context"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/coreos/etcd/clientv3"
	"time"
)

type EtcdController struct {
	beego.Controller
}
// @router / [get]
func (this *EtcdController) Get()  {
	cli,err := clientv3.New(clientv3.Config{
		Endpoints:[]string{"localhost:2379"},
		DialTimeout:5 * time.Second,
	})
	if err != nil {
		fmt.Println("connect failed,err:",err)
		return
	}
	fmt.Println("connect success")
	defer cli.Close()



	//设置1秒超时，访问etcd有超时控制
	ctx,cancle := context.WithTimeout(context.Background(), time.Second)
	//操作etcd
	_,err = cli.Put(ctx,"/logagent/conf/", "sample_value")
	//操作完毕，取消etcd
	cancle()
	if err != nil {
		fmt.Println("put failed, err:", err)
		return
	}
	fmt.Println("put success")

	ctx,cancle = context.WithTimeout(context.Background(), time.Second)
	resp,err := cli.Get(ctx,"/logagent/conf/")
	cancle()
	if err != nil {
		fmt.Println("get failed,err:", err)
		return
	}
	fmt.Println("get success")
	for _, ev:= range resp.Kvs {
		fmt.Printf("%s : %s\n", ev.Key, ev.Value)
	}
	this.ServeJSON()
}