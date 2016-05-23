package service

import (
	"github.com/gocraft/dbr"
	"github.com/golang/glog"
	"gopkg.in/olivere/elastic.v3"

	"github.com/zwh8800/md-blog-gen/conf"
)

var dbConn *dbr.Connection
var esClient *elastic.Client

func InitDb() (err error) {
	dbConn, err = dbr.Open(conf.Conf.DbConf.Driver, conf.Conf.DbConf.Dsn, nil)
	if err != nil {
		return err
	}
	dbConn.SetMaxOpenConns(conf.Conf.DbConf.MaxConnection)

	return nil
}

func InitElasticSearch() (err error) {
	esClient, err = elastic.NewClient(elastic.SetURL(conf.Conf.ElasticSearchConf.Url))
	if err != nil {
		return err
	}
	return nil
}

func init() {
	glog.Infoln("initilizing database...")

	if err := InitDb(); err != nil {
		glog.Fatalf("error occored: %s", err)
		panic(err)
	}
	if err := InitElasticSearch(); err != nil {
		glog.Fatalf("error occored: %s", err)
		panic(err)
	}
}
