package render

import (
	"html/template"
	"net/http"
	"path"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zwh8800/md-blog-gen/conf"
	"github.com/zwh8800/md-blog-gen/util"
)

const (
	templateDir = "template"
	commonDir   = "common"
)

type Render struct {
	templateName string
	template     string
	data         interface{}
}

func NewRender(template string, data interface{}) *Render {
	if conf.Conf.Env.Prod {
		if h, ok := data.(gin.H); ok {
			h["staticVer"] = "/v" + strconv.Itoa(conf.Conf.Env.StaticVersion)
			data = h
		}
	}

	return &Render{
		template,
		path.Join(templateDir, template),
		data,
	}
}

func (r *Render) Render(w http.ResponseWriter) error {
	util.WriteContentType(w, "text/html; charset=utf-8")
	t := template.New("")
	if _, err := t.ParseGlob(path.Join(templateDir, commonDir, "*")); err != nil {
		return err
	}
	if _, err := t.ParseFiles(r.template); err != nil {
		return err
	}
	return t.ExecuteTemplate(w, r.templateName, r.data)
}

func (r *Render) WriteContentType(w http.ResponseWriter) {
	util.WriteContentType(w, "text/html; charset=utf-8")
}
