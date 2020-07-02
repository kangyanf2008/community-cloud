package main

import (
	"community-cloud/config"
	"community-cloud/db"
	"community-cloud/logging"
	"community-cloud/web"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	//加载初始化.toml默认配置
	if err := config.LoadConfigAndSetDefault(); err != nil {
		panic(err.Error())
	}

	//初始化日志配置
	if err := logging.InitZap(&config.GetConf().LogConf); err != nil {
		panic("InitLogger:" + err.Error())
	}
	//初始化mysql库连接
	if err := db.InitDbConn(config.GetConf()); err != nil {
		panic("init db " + err.Error())
	}

}

func main() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	//启动web服务
	go web.Run()

	//测试函数
	//test.Test()

	//监听退出指令
	for s := range c {
		switch s {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			fmt.Println("退出", s)
			//服务退出
			web.Shutdown()
			time.Sleep(time.Second * 3) //等待三秒
			os.Exit(0)
		default:
			fmt.Println("other", s)
		}
	}

}
