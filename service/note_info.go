package service

import (
	"github.com/zwh8800/md-blog-gen/dao"
	"github.com/zwh8800/md-blog-gen/model"
)

func NoteById(id int64) (*model.Note, error) {
	sess := dbConn.NewSession(nil)
	return dao.NoteById(sess, id)
}

func NoteByNotename(notename string) (*model.Note, error) {
	sess := dbConn.NewSession(nil)
	return dao.NoteByNotename(sess, notename)
}
