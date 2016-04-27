package util

import (
	"math/rand"
	"time"

	"github.com/zwh8800/md-blog-gen/conf"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func HahaGenarate() string {
	return conf.Conf.Haha.Data[rand.Intn(len(conf.Conf.Haha.Data))]
}
