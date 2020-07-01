package main

import (
	"community-cloud/test"
	"community-cloud/web"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {

}

func main() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	//启动web服务
	go web.Run()

	//测试函数
	test.Test()

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
