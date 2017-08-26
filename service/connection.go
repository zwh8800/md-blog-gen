package service

import (
	"github.com/go-redis/redis"
	"github.com/gocraft/dbr"
	"github.com/golang/glog"
	"github.com/smartwalle/alipay"
	"github.com/zwh8800/md-blog-gen/conf"
	"gopkg.in/olivere/elastic.v3"
)

var dbConn *dbr.Connection
var esClient *elastic.Client
var redisClient *redis.Client

func newSession() *dbr.Session {
	return dbConn.NewSession(nil)
}

func newAlipayClient() *alipay.AliPay {
	return alipay.New(conf.Conf.Alipay.AppId,
		conf.Conf.Alipay.PartnerId,
		[]byte(conf.Conf.Alipay.PublicKey),
		[]byte(conf.Conf.Alipay.PrivateKey),
		conf.Conf.Alipay.Prod)
}

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

func InitRedis() (err error) {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     conf.Conf.Redis.Addr,
		Password: conf.Conf.Redis.Password,
		DB:       conf.Conf.Redis.DB,
	})
	if err := redisClient.Ping().Err(); err != nil {
		return err
	}
	return nil
}

func init() {
	glog.Infoln("initilizing connections...")

	if err := InitDb(); err != nil {
		glog.Fatalf("error occored: %s", err)
		panic(err)
	}
	if err := InitElasticSearch(); err != nil {
		glog.Fatalf("error occored: %s", err)
		panic(err)
	}
	if err := InitRedis(); err != nil {
		glog.Fatalf("error occored: %s", err)
		panic(err)
	}
}
