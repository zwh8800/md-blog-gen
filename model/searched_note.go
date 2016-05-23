package model

import "html/template"

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

func (note *SearchedNote) HighlightTagList() []template.HTML {
	if len(note.Highlight.TagList) != 0 {
		return note.Highlight.TagList
	} else {
		tagList := make([]template.HTML, 0, len(note.TagList))
		for _, tag := range note.TagList {
			tagList = append(tagList, template.HTML(tag))
		}
		return tagList
	}
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
				note.Highlight.Content = append(note.Highlight.Title, template.HTML(h))
			}
		case "tagList":
			for _, h := range value {
				note.Highlight.TagList = append(note.Highlight.Title, template.HTML(h))
			}
		}
	}
}
