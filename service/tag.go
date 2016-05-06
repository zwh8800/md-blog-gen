package service

import (
	"github.com/zwh8800/md-blog-gen/dao"
	"github.com/zwh8800/md-blog-gen/model"
)

func Tags() ([]*model.Tag, error) {
	sess := dbConn.NewSession(nil)
	tagList, err := dao.Tags(sess)
	if err != nil {
		return nil, err
	}

	return tagList, nil
}

func TagsByNoteId(noteId int64) ([]*model.Tag, error) {
	sess := dbConn.NewSession(nil)
	return dao.TagsByNoteId(sess, noteId)
}
