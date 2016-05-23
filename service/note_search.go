package service

import (
	"encoding/json"
	"strconv"

	"gopkg.in/olivere/elastic.v3"

	"github.com/golang/glog"
	"github.com/zwh8800/md-blog-gen/dao"
	"github.com/zwh8800/md-blog-gen/model"
	"github.com/zwh8800/md-blog-gen/util"
)

func NoteByKeyword(keyword string, page, limit int64) ([]*model.Note, map[int64][]*model.Tag, int64, error) {
	sess := dbConn.NewSession(nil)

	noteCount, err := dao.CountNoteByKeyword(sess, keyword)
	if err != nil {
		return nil, nil, 0, err
	}
	maxPage := (noteCount-1)/limit + 1

	page-- //数据库层的页数从0开始数
	noteList, err := dao.NoteByKeyword(sess, keyword, page, limit)
	if err != nil {
		return nil, nil, 0, err
	}
	if len(noteList) == 0 {
		return noteList, nil, maxPage, nil
	}

	noteIdList := make([]int64, 0, len(noteList))
	for _, note := range noteList {
		noteIdList = append(noteIdList, note.Id)
	}
	tagListMap, err := dao.TagsByNoteIds(sess, noteIdList)
	if err != nil {
		return nil, nil, 0, err
	}
	return noteList, tagListMap, maxPage, nil
}

func SearchNoteByKeyword(keyword string, page, limit int64) ([]*model.SearchedNote, int64, error) {
	page-- //数据库层的页数从0开始数
	offset := page * limit

	query := elastic.NewMultiMatchQuery(keyword,
		"notename", "title", "content", "tagList")
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
	glog.Infoln(util.JsonStringify(result, true))

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

func IsNoteIndexExist(uniqueId int64) (bool, error) {
	return esClient.Exists().
		Index("mdblog").
		Type("note").
		Id(strconv.FormatInt(uniqueId, 10)).
		Do()
}

func InsertOrUpdateNoteIndex(note *model.Note, tagList []*model.Tag) error {
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
