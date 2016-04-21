package util

import (
	"html/template"
	"io/ioutil"

	"github.com/golang/glog"
)

func DumpCss(filename string) template.CSS {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		glog.Errorln(err)
		return template.CSS("")
	}
	return template.CSS(string(data))
}

func DumpJs(filename string) template.JS {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		glog.Errorln(err)
		return template.JS("")
	}
	return template.JS(string(data))
}
