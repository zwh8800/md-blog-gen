package service

import (
	"github.com/zwh8800/md-blog-gen/dao"
	"github.com/zwh8800/md-blog-gen/model"
)

func NotesOrderByTime(page, limit int64) ([]*model.Note, map[int64][]*model.Tag, int64, error) {
	sess := dbConn.NewSession(nil)
	noteCount, err := dao.CountNote(sess)
	if err != nil {
		return nil, nil, 0, err
	}
	maxPage := (noteCount-1)/limit + 1

	page-- //数据库层的页数从0开始数
	noteList, err := dao.NotesByPage(sess, page, limit)
	if err != nil {
		return nil, nil, 0, err
	}

	noteIdList := make([]int64, 0)
	for _, note := range noteList {
		noteIdList = append(noteIdList, note.Id)
	}
	tagListMap, err := dao.TagsByNoteIds(sess, noteIdList)
	if err != nil {
		return nil, nil, 0, err
	}
	return noteList, tagListMap, maxPage, nil
}

func NotesWithoutTagOrderByTime(page, limit int64) ([]*model.Note, int64, error) {
	sess := dbConn.NewSession(nil)
	noteCount, err := dao.CountNote(sess)
	if err != nil {
		return nil, 0, err
	}
	maxPage := (noteCount-1)/limit + 1

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
	maxPage := (noteCount-1)/limit + 1

	page-- //数据库层的页数从0开始数
	noteIdList, err := dao.NoteIdsByPage(sess, page, limit)
	if err != nil {
		return nil, 0, err
	}

	return noteIdList, maxPage, nil
}

func NoteGroupByMonth() ([]*model.YearMonth, map[*model.YearMonth][]*model.Note, error) {
	sess := dbConn.NewSession(nil)
	return dao.NoteGroupByMonth(sess)
}

func NotesByMonth(month *model.YearMonth) ([]*model.Note, error) {
	sess := dbConn.NewSession(nil)
	return dao.NotesMonth(sess, month)
}
