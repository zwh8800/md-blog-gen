package service

import (
	"github.com/golang/glog"

	"github.com/zwh8800/md-blog-gen/dao"
	"github.com/zwh8800/md-blog-gen/model"
)

func SaveNoteList(noteList []*model.Note, tagListMap map[int64][]*model.Tag) error {
	sess := dbConn.NewSession(nil)
	for _, note := range noteList {
		tagList, _ := tagListMap[note.UniqueId]

		indexExist, err := IsNoteDocumentExist(note.UniqueId)
		if err != nil {
			glog.Error(err)
			continue
		}
		if indexExist {
			modified, err := dao.IsNoteModified(sess, note)
			if err != nil {
				glog.Error(err)
				continue
			}
			if !modified {
				continue
			}
			glog.Infoln("note", note.Title, "modified, updating")
		} else {
			glog.Infoln("note", note.Title, "index not exist, indexing")
		}

		if err := SaveNote(note, tagList); err != nil {
			glog.Error(err)
			continue
		}
		if err := InsertOrUpdateNoteDocument(note, tagList); err != nil {
			glog.Error(err)
			continue
		}
	}

	tx, err := sess.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if tx != nil {
			tx.RollbackUnlessCommitted()
		}
	}()
	dao.RemoveUnpublishedNote(tx, noteList)

	return tx.Commit()
}

func SaveNote(note *model.Note, tagList []*model.Tag) error {
	sess := dbConn.NewSession(nil)
	tx, err := sess.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if tx != nil {
			tx.RollbackUnlessCommitted()
		}
	}()

	if err := dao.InsertOrUpdateNote(tx, note); err != nil {
		return err
	}

	for _, tag := range tagList {
		realTag, err := dao.SelectTagOrInsertIfNotExists(tx, tag)
		if err != nil {
			return err
		}
		*tag = *realTag
		if err := dao.InsertNoteTag(tx, note, tag); err != nil {
			// tag already exists
			continue
		}
	}
	if err := dao.DeleteNoteTagsNotExist(tx, note, tagList); err != nil {
		return err
	}

	return tx.Commit()
}
