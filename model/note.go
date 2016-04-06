package model

import (
	"html/template"
	"regexp"
	"strings"
	"time"
)

const NoteTableName = "Note"

type Note struct {
	Id        int64     `db:"id" json:"id" struct:"id"`
	UniqueId  int64     `db:"unique_id" json:"uniqueId" struct:"unique_id"`
	Title     string    `db:"title" json:"title" struct:"title"`
	Url       string    `db:"url" json:"url" struct:"url"`
	Content   string    `db:"content" json:"content" struct:"content"`
	Timestamp time.Time `db:"timestamp" json:"timestamp" struct:"timestamp"`
	Removed   bool      `db:"removed" json:"removed" struct:"removed"`
}

func NewNote(uniqueId int64, title string, url string, timestamp time.Time) *Note {
	return &Note{
		UniqueId:  uniqueId,
		Title:     title,
		Url:       url,
		Timestamp: timestamp,
		Removed:   false,
	}
}

func (obj *Note) UnescapedContent() template.HTML {
	return template.HTML(obj.Content)
}

func (obj *Note) FillContent(content string) {
	obj.Content = content
}

func (obj *Note) Preview() string {
	const maxLength = 200
	src := obj.Content
	//将HTML标签全转换成小写
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllStringFunc(src, strings.ToLower)

	//去除STYLE
	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	src = re.ReplaceAllString(src, "")

	//去除SCRIPT
	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	src = re.ReplaceAllString(src, "")

	//去除所有尖括号内的HTML代码，并换成换行符
	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllString(src, "\n")

	//去除连续的换行符
	re, _ = regexp.Compile("\\s{2,}")
	src = re.ReplaceAllString(src, "\n")

	runes := ([]rune)(src)
	if len(runes) < maxLength {
		return string(runes)
	} else {
		return string(runes[:maxLength])
	}
}

func (obj *Note) FormattedTimestamp() string {
	return obj.Timestamp.Local().Format("2006-01-02 15:04")
}
