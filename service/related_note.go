package service

import (
	"github.com/zwh8800/md-blog-gen/conf"
	"github.com/zwh8800/md-blog-gen/dao"
	"github.com/zwh8800/md-blog-gen/model"
)

func RelatedNote(id int64) ([]*model.Note, error) {
	sess := dbConn.NewSession(nil)
	tagIds, err := dao.TagIdsByNoteId(sess, id)
	if err != nil {
		return nil, err
	}
	noteIds, err := dao.NoteIdsByTagIds(sess, tagIds)
	if err != nil {
		return nil, err
	}
	noteIds = removeInt64Slice(noteIds, id)

	return dao.NotesByIds(sess, noteIds, conf.Conf.Site.NotePerPage)
}

func removeInt64Slice(s []int64, n int64) []int64 {
	for i := 0; i < len(s); i++ {
		if s[i] == n {
			s = append(s[:i], s[i+1:]...)
			i--
		}
	}
	return s
}
