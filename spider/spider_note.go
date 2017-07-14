package spider

import (
	"bytes"
	"html/template"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
	"unicode"

	"github.com/PuerkitoBio/goquery"
	"github.com/golang/glog"

	"github.com/zwh8800/md-blog-gen/conf"
	"github.com/zwh8800/md-blog-gen/model"
	"github.com/zwh8800/md-blog-gen/util"
)

func downloadImg(src string) (string, error) {
	resp, err := http.Get(src)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	outFilename := path.Join(conf.Conf.Site.StaticUrl, "img", util.MD5(src))
	if ext := path.Ext(src); ext == ".jpg" ||
		ext == ".png" ||
		ext == ".gif" ||
		ext == ".jpeg" ||
		ext == ".tif" ||
		ext == ".tiff" {
		outFilename = outFilename + ext
	}

	outFile, err := os.OpenFile(outFilename, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		if os.IsExist(err) {
			return outFilename, nil
		}
		return "", err
	}
	if _, err := io.Copy(outFile, resp.Body); err != nil {
		return "", err
	}

	return outFilename, nil
}

func handleTag(content *goquery.Selection) {
	content.Find("h1").Next().Find("code").Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			s.BeforeHtml(`<i class="icon-tags"></i>`)
		}
		if s.Text() == conf.Conf.Spider.SpiderTag {
			s.Remove()
		}
		s.WrapAllHtml(`<a href="/tag/` + s.Text() + `"></a>`)
	})
}

func handleTime(content *goquery.Selection, timestamp time.Time) {
	t := template.Must(template.New("time").Parse(`<h6>
        <i class="icon-time"></i>
        {{ .Format "2006-01-02 15:04 PM" }}
    </h6>`))
	output := &bytes.Buffer{}

	if err := t.Execute(output, timestamp); err != nil {
		panic(err)
	}

	content.Find("h1").Next().AfterHtml(output.String())
}

func transformToNotename(notename string) string {
	notename = strings.TrimSpace(notename)
	notename = strings.ToLower(notename)
	notename = strings.Join(strings.FieldsFunc(notename, func(r rune) bool {
		return !(unicode.IsLetter(r) ||
			unicode.IsNumber(r))
	}), "-")

	return strings.ToLower(notename)
}

func getArchiveDate(notename string) (*time.Time, bool) {
	if strings.Index(notename, "archive-") != 0 {
		return nil, false
	}
	timestamp, err := time.ParseInLocation("20060102", notename[8:], time.Local)
	if err != nil {
		glog.Error(err)
		timestamp, _ = time.ParseInLocation("20060102", "20160101", time.Local) // default date
	}
	return &timestamp, true
}

func handleNotename(content *goquery.Selection, note *model.Note) {
	a := content.Find("p a[href='/notename/']")
	attr, ok := a.Attr("title")
	if ok {
		notename := transformToNotename(attr)
		note.Notename.Valid = true
		note.Notename.String = notename
		if timestamp, ok := getArchiveDate(notename); ok {
			note.Timestamp = *timestamp
		}

		a.SetAttr("href", util.GetNoteUrlByNotename(notename))
	} else if notename := transformToNotename(util.YoudaoTranslate(note.Title)); notename != "" {
		note.Notename.Valid = true
		note.Notename.String = notename
	} else if notename := transformToNotename(util.Pinyin(note.Title)); notename != "" {
		note.Notename.Valid = true
		note.Notename.String = notename
	} else {
		note.Notename.Valid = false
	}
}

func FindNoteContent(note *model.Note) {
	doc, err := goquery.NewDocument(note.Url)
	if err != nil {
		glog.Errorln(err)
		return
	}
	content := doc.Find("#wmd-preview")
	content.Find("img").Each(func(i int, s *goquery.Selection) {
		src, ok := s.Attr("src")
		if !ok {
			return
		}
		dest, err := downloadImg(src)
		if err != nil {
			glog.Errorln(err)
			return
		}

		s.SetAttr("src", "/"+dest)
	})
	handleTag(content)
	handleNotename(content, note)
	handleTime(content, note.Timestamp.Local())

	html, err := content.Html()
	if err != nil {
		glog.Warning(err)
		return
	}
	note.FillContent(html)
}
