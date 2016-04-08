package render

import (
	"bytes"
	"html/template"
	"net/http"
	"path"

	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/html"

	"github.com/zwh8800/md-blog-gen/util"
)

const (
	templateDir = "template"
	commonDir   = "common"
)

var m *minify.M

func init() {
	m = minify.New()
	m.AddFunc("text/html", html.Minify)
}

type Render struct {
	templateName string
	template     string
	data         interface{}
}

func NewRender(template string, data interface{}) *Render {
	return &Render{
		template,
		path.Join(templateDir, template),
		data,
	}
}

func (r *Render) Render(w http.ResponseWriter) error {
	util.WriteContentType(w, []string{"text/html; charset=utf-8"})
	t := template.New("")
	if _, err := t.ParseGlob(path.Join(templateDir, commonDir, "*")); err != nil {
		return err
	}
	if _, err := t.ParseFiles(r.template); err != nil {
		return err
	}
	buf := &bytes.Buffer{}

	if err := t.ExecuteTemplate(buf, r.templateName, r.data); err != nil {
		return err
	}

	return m.Minify("text/html", w, buf)
}
