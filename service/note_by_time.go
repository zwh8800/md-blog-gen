package service

import (
	"github.com/zwh8800/md-blog-gen/dao"
	"github.com/zwh8800/md-blog-gen/model"
)

func NotesOrderByTime(page, limit int64) ([]*model.Note, int64, error) {
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
