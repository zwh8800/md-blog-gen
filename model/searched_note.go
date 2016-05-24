package model

import (
	"html/template"
	"strings"
)

type searchNoteHighlight struct {
	Title   []template.HTML
	Content []template.HTML
	TagList []template.HTML
}

type SearchedNote struct {
	NoteDetail
	Highlight searchNoteHighlight
}

func NewSearchedNote() *SearchedNote {
	return &SearchedNote{
		Highlight: searchNoteHighlight{
			Title:   make([]template.HTML, 0),
			Content: make([]template.HTML, 0),
			TagList: make([]template.HTML, 0),
		},
	}
}

func (note *SearchedNote) HighlightTitle() template.HTML {
	if len(note.Highlight.Title) != 0 {
		return note.Highlight.Title[0]
	} else {
		return template.HTML(note.Title)
	}
}

type tagForRender struct {
	Html template.HTML
	Name string
}

func highlightTagToTag(highlight template.HTML) string {
	h := string(highlight)
	h = strings.Replace(h, "<em>", "", -1)
	h = strings.Replace(h, "</em>", "", -1)
	return h
}

func (note *SearchedNote) HighlightTagList() []*tagForRender {
	filter := make(map[string]bool)
	tagList := make([]*tagForRender, 0, len(note.Highlight.TagList))
	for _, tag := range note.Highlight.TagList {
		name := highlightTagToTag(tag)
		if _, ok := filter[name]; !ok {
			filter[name] = true
			tagList = append(tagList, &tagForRender{
				tag,
				name,
			})
		}
	}
	for _, tag := range note.TagList {
		if _, ok := filter[tag]; !ok {
			filter[tag] = true
			tagList = append(tagList, &tagForRender{
				template.HTML(tag),
				tag,
			})
		}
	}

	return tagList
}

func (note *SearchedNote) HighlightContent() []template.HTML {
	if len(note.Highlight.Content) != 0 {
		return note.Highlight.Content
	} else {
		return []template.HTML{
			template.HTML(string([]rune(note.Content)[:160])),
		}
	}
}

func (note *SearchedNote) FillHighlight(highlight map[string][]string) {
	for key, value := range highlight {
		switch key {
		case "title":
			for _, h := range value {
				note.Highlight.Title = append(note.Highlight.Title, template.HTML(h))
			}
		case "content":
			for _, h := range value {
				note.Highlight.Content = append(note.Highlight.Content, template.HTML(h))
			}
		case "tagList":
			for _, h := range value {
				note.Highlight.TagList = append(note.Highlight.TagList, template.HTML(h))
			}
		}
	}
}
