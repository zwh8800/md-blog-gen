package conf

import (
	"path/filepath"

	"gopkg.in/gcfg.v1"
)

type config struct {
	DbConf struct {
		Driver string
		Dsn    string
	}
	Env struct {
		Prod       bool
		ServerPort string
		StaticDir  string
	}
	Spider struct {
		StartUrl  string
		SpiderTag string
	}
	UrlPush struct {
		Baidu string
	}
	Site struct {
		Name    string
		BaseUrl string
		NoteUrl string
		TagUrl  string
		PageUrl string
	}
}

var Conf config

func ReadConf(filename string) error {
	absFile := filename
	if !filepath.IsAbs(filename) {
		var err error
		absFile, err = filepath.Abs(filename)
		if err != nil {
			return err
		}
	}
	return gcfg.ReadFileInto(&Conf, absFile)
}
