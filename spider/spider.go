package spider

import (
	"runtime/debug"
	"sync"

	"github.com/golang/glog"
	"gopkg.in/go-playground/pool.v1"

	"github.com/zwh8800/md-blog-gen/model"
	"github.com/zwh8800/md-blog-gen/service"
)

var waitFinish sync.WaitGroup

func Go() {
	waitFinish.Add(1)
	defer waitFinish.Done()
	glog.Info("spider started")
	defer func() {
		if err := recover(); err != nil {
			glog.Errorln("panic in spider recovered:", err, string(debug.Stack()))
		}
	}()

	startPage, err := findStartPage()
	if err != nil {
		glog.Error(err)
		return
	}
	noteList := findNoteList(startPage)
	tagListMap := findTagListMap(startPage)

	p := pool.NewPool(4, len(noteList))
	for _, note := range noteList {
		p.Queue(func(job *pool.Job) {
			note := job.Params()[0].(*model.Note)
			FindNoteContent(note)
		}, note)
	}
	for result := range p.Results() {
		err, ok := result.(*pool.ErrRecovery)
		if ok {
			panic(err)
		}
	}

	if err := service.SaveNoteList(noteList, tagListMap); err != nil {
		glog.Error(err)
		return
	}

	glog.Info("spider finished")
}

func WaitFinish() {
	waitFinish.Wait()
}
