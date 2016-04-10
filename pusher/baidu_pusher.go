package pusher

import (
	"github.com/golang/glog"

	"github.com/zwh8800/md-blog-gen/service"
	"github.com/zwh8800/md-blog-gen/util"
)

func BaiduPush() {
	glog.Infoln("Baidu push start")
	noteIdList, _, err := service.NoteIdsOrderByTime(0, 10)
	if err != nil {
		glog.Error(err)
		return
	}
	urls := make([]string, 0, len(noteIdList))
	for _, noteId := range noteIdList {
		urls = append(urls, util.GetNoteUrl(noteId))
	}
	respData, err := util.PushUrlToBaidu(urls)
	if err != nil {
		glog.Error(err)
		return
	}
	glog.Infoln("Baidu response: ", string(respData))
	glog.Infoln("Baidu push finish")
}
