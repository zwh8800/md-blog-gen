package crontab

import (
	"github.com/robfig/cron"

	"github.com/zwh8800/md-blog-gen/spider"
)

func Go() {
	crontab := cron.New()
	crontab.AddFunc("@every 5m", spider.Go)
	crontab.Start()
}
