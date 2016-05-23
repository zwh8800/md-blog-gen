package service

import (
	"strconv"

	"github.com/zwh8800/md-blog-gen/dao"
	"github.com/zwh8800/md-blog-gen/model"
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
