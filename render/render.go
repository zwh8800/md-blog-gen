package render

import (
	"html/template"
	"net/http"
	"path"
)

const (
	templateDir = "template"
	layout      = "layout.html"
)

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

func writeContentType(w http.ResponseWriter, value []string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = value
	}
}

func (r *Render) Render(w http.ResponseWriter) error {
	writeContentType(w, []string{"text/html; charset=utf-8"})
	t := template.New("")
	if _, err := t.ParseFiles(path.Join(templateDir, layout), r.template); err != nil {
		return err
	}
	return t.ExecuteTemplate(w, r.templateName, r.data)
}
