package service

import (
	"github.com/gocraft/dbr"
	"github.com/golang/glog"

	"github.com/zwh8800/md-blog-gen/conf"
)

var dbConn *dbr.Connection

func InitDb() (err error) {
	dbConn, err = dbr.Open(conf.Conf.DbConf.Driver, conf.Conf.DbConf.Dsn, nil)
	if err != nil {
		return err
	}
	dbConn.SetMaxOpenConns(conf.Conf.DbConf.MaxConnection)

	return nil
}

func init() {
	glog.Infoln("initilizing database...")

	if err := InitDb(); err != nil {
		glog.Fatalf("error occored: %s", err)
		panic(err)
	}
}
