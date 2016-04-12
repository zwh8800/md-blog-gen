package util

import (
	"net/url"
	"path"
	"strconv"

	"github.com/golang/glog"
	"github.com/zwh8800/md-blog-gen/conf"
)

func GetNoteBase() string {
	return "/" + conf.Conf.Site.NoteUrl
}

func GetNoteUrl(id int64) string {
	baseUrl, err := url.Parse(conf.Conf.Site.BaseUrl)
	if err != nil {
		glog.Error(err)
		return ""
	}
	noteChild := path.Join(GetNoteBase(), strconv.FormatInt(id, 10))
	u, err := url.Parse(noteChild)
	if err != nil {
		glog.Error(err)
		return ""
	}
	return baseUrl.ResolveReference(u).String()
}

func GetTagBase() string {
	return "/" + conf.Conf.Site.TagUrl
}

func GetTagUrl(id int64) string {
	baseUrl, err := url.Parse(conf.Conf.Site.BaseUrl)
	if err != nil {
		glog.Error(err)
		return ""
	}
	noteChild := path.Join(GetTagBase(), strconv.FormatInt(id, 10))
	u, err := url.Parse(noteChild)
	if err != nil {
		glog.Error(err)
		return ""
	}
	return baseUrl.ResolveReference(u).String()
}

func GetTagNameUrl(name string) string {
	baseUrl, err := url.Parse(conf.Conf.Site.BaseUrl)
	if err != nil {
		glog.Error(err)
		return ""
	}
	noteChild := path.Join(GetTagBase(), name)
	u, err := url.Parse(noteChild)
	if err != nil {
		glog.Error(err)
		return ""
	}
	return baseUrl.ResolveReference(u).String()
}

func GetPageBase() string {
	return "/" + conf.Conf.Site.PageUrl
}

func GetPageUrl(id int64) string {
	baseUrl, err := url.Parse(conf.Conf.Site.BaseUrl)
	if err != nil {
		glog.Error(err)
		return ""
	}
	noteChild := path.Join(GetPageBase(), strconv.FormatInt(id, 10))
	u, err := url.Parse(noteChild)
	if err != nil {
		glog.Error(err)
		return ""
	}
	return baseUrl.ResolveReference(u).String()
}
