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
	"github.com/zwh8800/md-blog-gen/crontab"
	"github.com/zwh8800/md-blog-gen/route"
	"github.com/zwh8800/md-blog-gen/spider"
)

func main() {
	defer glog.Flush()
	startServer()
	handleSignal()
	glog.Infoln("gracefully shutdown")
}

func startServer() {
	glog.Infoln("starting...")

	if conf.Conf.Env.Prod {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	route.Route(r)

	go func() {
		err := r.Run(fmt.Sprintf("%v:%v", "", conf.Conf.Env.ServerPort))
		if err != nil {
			glog.Fatalf("error occored: %s", err)
			panic(err)
		}
	}()
	crontab.Go()
	spider.Go()
	glog.Infoln("server started")
}

func handleSignal() {
	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, os.Kill, os.Interrupt, syscall.SIGTERM)
	<-signalChan
	glog.Infoln("signal received")
}
