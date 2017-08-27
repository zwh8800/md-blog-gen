package conf

import (
	"flag"

	"github.com/BurntSushi/toml"
	"github.com/golang/glog"
)

type config struct {
	DbConf struct {
		Driver        string
		Dsn           string
		MaxConnection int
	}
	ElasticSearchConf struct {
		Url string
	}
	Redis struct {
		Addr     string
		Password string
		DB       int
	}
	Env struct {
		Prod          bool
		ServerPort    int
		StaticVersion int
	}
	Spider struct {
		StartUrl  string
		SpiderTag string
	}
	UrlPush struct {
		Baidu string
	}
	Site struct {
		Name        string
		FuckIcp     string
		Description string
		BaseUrl     string
		NoteUrl     string
		TagUrl      string
		PageUrl     string
		RssUrl      string
		ArchiveUrl  string
		SearchUrl   string
		StaticUrl   string
		AuthorName  string
		AuthorEmail string
		NotePerPage int64
		Language    string
		LicenseName string
		LicenseUrl  string
		ICP         string
	}
	Social struct {
		Github   string
		Weibo    string
		Twitter  string
		Linkedin string
		About    string
	}
	Youdao struct {
		ApiUrl string
	}
	Haha struct {
		Data []string
	}
	Alipay struct {
		Prod            bool
		AppId           string
		PartnerId       string
		PublicKey       string
		PrivateKey      string
		AlipayPublicKey string
	}
	Crypto struct {
		AesKey string
	}
}

var Conf config

func init() {
	configFilename := flag.String("config", "md-blog-gen.toml", "specify a config file")
	flag.Parse()
	glog.Infoln("configuring...")

	if _, err := toml.DecodeFile(*configFilename, &Conf); err != nil {
		panic(err)
	}
}
