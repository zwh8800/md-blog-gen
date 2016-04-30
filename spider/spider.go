package spider

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/PuerkitoBio/goquery"
	"github.com/golang/glog"
	"github.com/mozillazg/go-pinyin"

	"github.com/zwh8800/md-blog-gen/conf"
	"github.com/zwh8800/md-blog-gen/model"
	"github.com/zwh8800/md-blog-gen/service"
	"github.com/zwh8800/md-blog-gen/util"
)

func findBlogTagUl(doc *goquery.Document) *goquery.Selection {
	var result *goquery.Selection
	doc.Find("#file-list").Children().Each(func(i int, s *goquery.Selection) {
		ul := s.Find("ul")
		tagName, _ := ul.Find(".tag-item.item").Attr("tag-name")
		if tagName == conf.Conf.Spider.SpiderTag {
			result = ul
		}
	})
	return result
}

func findAllBlogTagNotes(doc *goquery.Document) []*model.Note {
	ul := findBlogTagUl(doc)
	liList := ul.Find(".file-item.item")
	result := make([]*model.Note, 0, liList.Length())
	liList.Each(func(i int, s *goquery.Selection) {
		timestampStr, _ := s.Attr("file-created-date")
		timestamp, err := time.Parse(time.RFC3339Nano, timestampStr)
		if err != nil {
			glog.Warning(err)
			return
		}
		a := s.Find("a")
		url, _ := a.Attr("href")
		span := a.Find("span")
		idStr, _ := span.Attr("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			glog.Warning(err)
			return
		}
		title := span.Text()
		result = append(result, model.NewNote(id, title, url, timestamp))
	})
	return result
}

func findTagListMap(doc *goquery.Document) map[int64][]*model.Tag {
	tagListMap := make(map[int64][]*model.Tag)
	doc.Find("#file-list .tag-list").Each(func(i int, s *goquery.Selection) {
		tagSel := s.Find(".tag-item.item")
		notesSel := s.Find(".file-item.item")
		tagName, _ := tagSel.Attr("tag-name")
		tag := model.NewTag(tagName)
		if tagName == conf.Conf.Spider.SpiderTag {
			return
		}

		notesSel.Each(func(i int, s *goquery.Selection) {
			span := s.Find("a span")
			idStr, _ := span.Attr("id")
			id, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				glog.Warning(err)
				return
			}
			tagList, ok := tagListMap[id]
			if !ok {
				tagList = make([]*model.Tag, 0)
			}
			tagList = append(tagList, tag)
			tagListMap[id] = tagList
		})

	})

	return tagListMap
}

func downloadImg(src string) (string, error) {
	resp, err := http.Get(src)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	outFilename := path.Join(conf.Conf.Site.StaticUrl, "img", util.MD5(src))

	outFile, err := os.OpenFile(outFilename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
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

func transNotename(notename string) string {
	notename = strings.TrimSpace(notename)
	notename = strings.ToLower(notename)
	notename = strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return '-'
		} else {
			return r
		}
	}, notename)

	return url.QueryEscape(notename)
}

// 有道API: http://fanyi.youdao.com/openapi?path=data-mode
type youdaoResponse struct {
	Translation []string `json:"translation"`
	Query       string   `json:"query"`
	ErrorCode   int      `json:"errorCode"`
}

func translateTitleToNotename(title string) string {
	if conf.Conf.Youdao.ApiUrl == "" {
		return ""
	}

	u, err := url.Parse(conf.Conf.Youdao.ApiUrl)
	if err != nil {
		glog.Errorln(err)
		return ""
	}
	q := u.Query()
	q.Set("q", title)
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		glog.Errorln(err)
		return ""
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Errorln(err)
		return ""
	}
	var youdaoData youdaoResponse

	if err := json.Unmarshal(data, &youdaoData); err != nil {
		glog.Errorln(err)
		return ""
	}
	if youdaoData.ErrorCode != 0 {
		glog.Errorln("youdaoData.ErrorCode is ", youdaoData.ErrorCode)
		return ""
	}
	if len(youdaoData.Translation) == 0 {
		glog.Errorln("youdaoData.Translation length is 0")
		return ""
	}

	return transNotename(youdaoData.Translation[0])
}

func pinyinNotename(title string) string {
	py := pinyin.Pinyin(title, pinyin.NewArgs())
	sb := &bytes.Buffer{}
	for _, p := range py {
		sb.WriteString(p[0])
		sb.WriteRune(' ')
	}
	return transNotename(sb.String())
}

func handleNotename(content *goquery.Selection, note *model.Note) {
	a := content.Find("p a[href='/notename/']")
	attr, ok := a.Attr("title")
	if ok {
		notename := transNotename(attr)
		note.Notename.Valid = true
		note.Notename.String = notename

		a.SetAttr("href", util.GetNoteUrlByNotename(notename))
	} else if notename := translateTitleToNotename(note.Title); notename != "" {
		note.Notename.Valid = true
		note.Notename.String = notename
	} else if notename := pinyinNotename(note.Title); notename != "" {
		note.Notename.Valid = true
		note.Notename.String = notename
	} else {
		note.Notename.Valid = false
	}
}

func findNoteContent(note *model.Note) {
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
	handleTime(content, note.Timestamp.Local())
	handleNotename(content, note)

	html, err := content.Html()
	if err != nil {
		glog.Warning(err)
		return
	}
	note.FillContent(html)
}

func Go() {
	glog.Info("spider started")
	doc, err := goquery.NewDocument(conf.Conf.Spider.StartUrl)
	if err != nil {
		glog.Errorln(err)
		return
	}
	noteList := findAllBlogTagNotes(doc)
	tagListMap := findTagListMap(doc)
	for _, note := range noteList {
		findNoteContent(note)
	}

	if err := service.SaveNoteList(noteList, tagListMap); err != nil {
		glog.Error(err)
	}

	glog.Info("spider finished")
}
