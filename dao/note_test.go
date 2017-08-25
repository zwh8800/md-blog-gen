package dao

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gocraft/dbr"

	"github.com/zwh8800/md-blog-gen/model"
)

func TestRemoveUnpublishedNote(t *testing.T) {
	dbConn, err := dbr.Open("mysql", "root:9310308800@tcp(127.0.0.1:3306)/mdblog?charset=utf8mb4&parseTime=true", nil)
	if err != nil {
		t.Error(err)
		return
	}
	noteList := [...]*model.Note{
		&model.Note{
			UniqueId: 332154,
		},
	}
	sess := newSession()
	tx, err := sess.Begin()
	if err != nil {
		t.Error(err)
		return
	}
	defer func() {
		if tx != nil {
			tx.RollbackUnlessCommitted()
		}
	}()
	RemoveUnpublishedNote(tx, noteList[:])

	if err := tx.Commit(); err != nil {
		t.Error(err)
		return
	}
}
