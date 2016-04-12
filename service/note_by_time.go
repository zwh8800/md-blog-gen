package service

import (
	"github.com/golang/glog"

	"github.com/zwh8800/md-blog-gen/dao"
	"github.com/zwh8800/md-blog-gen/model"
)

func NotesOrderByTime(page, limit int64) ([]*model.Note, map[int64][]*model.Tag, int64, error) {
	sess := dbConn.NewSession(nil)
	noteCount, err := dao.CountNote(sess)
	if err != nil {
		return nil, nil, 0, err
	}
	maxPage := noteCount/limit + 1

	page-- //数据库层的页数从0开始数
	noteList, err := dao.NotesByPage(sess, page, limit)
	if err != nil {
		return nil, nil, 0, err
	}
	tagListMap := make(map[int64][]*model.Tag)
	for _, note := range noteList {
		tags, err := dao.TagsByNoteId(sess, note.Id)
		if err != nil {
			glog.Warning(err)
			tags = make([]*model.Tag, 0)
		}
		tagListMap[note.Id] = tags
	}
	return noteList, tagListMap, maxPage, nil
}

func NotesWithoutTagOrderByTime(page, limit int64) ([]*model.Note, int64, error) {
	sess := dbConn.NewSession(nil)
	noteCount, err := dao.CountNote(sess)
	if err != nil {
		return nil, 0, err
	}
	maxPage := noteCount/limit + 1

	page-- //数据库层的页数从0开始数
	noteList, err := dao.NotesByPage(sess, page, limit)
	if err != nil {
		return nil, 0, err
	}

	return noteList, maxPage, nil
}

func NoteIdsOrderByTime(page, limit int64) ([]int64, int64, error) {
	sess := dbConn.NewSession(nil)
	noteCount, err := dao.CountNote(sess)
	if err != nil {
		return nil, 0, err
	}
	maxPage := noteCount/limit + 1

	page-- //数据库层的页数从0开始数
	noteIdList, err := dao.NoteIdsByPage(sess, page, limit)
	if err != nil {
		return nil, 0, err
	}

	return noteIdList, maxPage, nil
}
