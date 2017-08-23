package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"

	"github.com/zwh8800/md-blog-gen/conf"
	"github.com/zwh8800/md-blog-gen/controller"
	"github.com/zwh8800/md-blog-gen/crontab"
	"github.com/zwh8800/md-blog-gen/service"
	"github.com/zwh8800/md-blog-gen/spider"
)

func main() {
	defer glog.Flush()
	startServer()
	handleSignal()
	glog.Infoln("signal received")
	stopServer()
	glog.Infoln("gracefully shutdown")
}

func startServer() {
	glog.Infoln("starting...")

	if conf.Conf.Env.Prod {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	controller.Route(r)

	go func() {
		err := r.Run(fmt.Sprintf("%v:%v", "", conf.Conf.Env.ServerPort))
		if err != nil {
			glog.Fatalf("error occored: %s", err)
			panic(err)
		}
	}()
	if err := service.CreateIndexAndMappingIfNotExist(); err != nil {
		glog.Fatalf("error occored: %s", err)
		panic(err)
	}
	crontab.Go()
	go spider.Go()
	service.RemoveAllCache()
	glog.Infoln("server started")
}

func handleSignal() {
	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, os.Kill, os.Interrupt, syscall.SIGTERM)
	<-signalChan
}

func stopServer() {
	spider.WaitFinish()
}
