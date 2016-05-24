package service

import (
	"encoding/json"
	"strconv"

	"gopkg.in/olivere/elastic.v3"

	"github.com/zwh8800/md-blog-gen/model"
)

func SearchNoteByKeyword(keyword string, page, limit int64) ([]*model.SearchedNote, int64, error) {
	page-- //数据库层的页数从0开始数
	offset := page * limit

	query := elastic.NewMultiMatchQuery(keyword).
		FieldWithBoost("notename", 1).
		FieldWithBoost("tagList", 2).
		FieldWithBoost("content", 4).
		FieldWithBoost("title", 4)
	highlight := elastic.NewHighlight().
		Field("content").
		Field("title").
		Field("tagList")

	result, err := esClient.Search().
		Index("mdblog").
		Type("note").
		Query(query).
		From(int(offset)).
		Size(int(limit)).
		Highlight(highlight).Do()
	if err != nil {
		return nil, 0, err
	}

	if result.Hits == nil {
		return nil, 0, nil
	}
	maxPage := (result.TotalHits()-1)/limit + 1

	noteList := make([]*model.SearchedNote, 0, len(result.Hits.Hits))
	for _, hit := range result.Hits.Hits {
		note := model.NewSearchedNote()
		err := json.Unmarshal(*hit.Source, note)
		if err != nil {
			return nil, 0, err
		}
		note.FillHighlight(hit.Highlight)

		noteList = append(noteList, note)
	}

	return noteList, maxPage, nil
}

func CreateIndexAndMappingIfNotExist() error {
	exist, err := IsMdBlogIndexExist()
	if err != nil {
		return err
	}
	if !exist {
		if err := CreateIndex(); err != nil {
			return err
		}
	}
	if err := CreateNoteMapping(); err != nil {
		return err
	}
	return nil
}

func IsMdBlogIndexExist() (bool, error) {
	return esClient.IndexExists("mdblog").Do()
}

func CreateIndex() error {
	_, err := esClient.CreateIndex("mdblog").Do()
	return err
}

func CreateNoteMapping() error {
	_, err := esClient.PutMapping().
		Index("mdblog").
		Type("note").
		BodyString(`{
			"note": {
				"properties": {
					"id": {
						"type": "long"
					},
					"title": {
						"type": "string",
						"term_vector": "with_positions_offsets",
						"analyzer": "ik_syno",
						"search_analyzer": "ik_syno"
					},
					"content": {
						"type": "string",
						"term_vector": "with_positions_offsets",
						"analyzer": "ik_syno",
						"search_analyzer": "ik_syno"
					},
					"notename": {
						"type": "string"
					},
					"tagList": {
						"type": "string",
						"term_vector": "with_positions_offsets",
						"analyzer": "ik_syno",
						"search_analyzer": "ik_syno"
					},
					"timestamp": {
						"type": "date",
						"index": "not_analyzed"
					},
					"lastModified": {
						"type": "date",
						"index": "not_analyzed"
					}
				}
			}
		}`).
		Do()
	return err
}

func IsNoteDocumentExist(uniqueId int64) (bool, error) {
	return esClient.Exists().
		Index("mdblog").
		Type("note").
		Id(strconv.FormatInt(uniqueId, 10)).
		Do()
}

func InsertOrUpdateNoteDocument(note *model.Note, tagList []*model.Tag) error {
	tagNameList := make([]string, 0, len(tagList))
	for _, tag := range tagList {
		tagNameList = append(tagNameList, tag.Name)
	}
	noteDetail := model.NoteDetail{
		Id:           note.Id,
		Notename:     note.Notename,
		Title:        note.Title,
		Content:      note.ContentText(),
		Timestamp:    note.Timestamp,
		LastModified: note.LastModified,
		TagList:      tagNameList,
	}

	_, err := esClient.Index().
		Index("mdblog").
		Type("note").
		Id(strconv.FormatInt(note.UniqueId, 10)).
		BodyJson(noteDetail).
		Do()
	if err != nil {
		return err
	}

	return nil
}
