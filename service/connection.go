package service

import (
	"github.com/gocraft/dbr"

	"github.com/zwh8800/md-blog-gen/conf"
)

var dbConn *dbr.Connection

func InitDb() (err error) {
	dbConn, err = dbr.Open(conf.Conf.DbConf.Driver, conf.Conf.DbConf.Dsn, nil)
	return err
}
