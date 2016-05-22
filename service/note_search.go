package service

import (
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

func InsertOrUpdateNoteIndex(note *model.Note, tagList []*model.Tag) error {

	return nil
}
