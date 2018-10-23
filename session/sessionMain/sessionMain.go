package sessionMain

import (
	"fmt"
	"github.com/learn_to_build_web_application_with_go/session"
	"net/http"
	"time"
)

var globalSessions *session.Manager

func init() {
	//初始化全局session管理器，一般放在main包中，应用启动时初始化
	globalSessions, _ = session.NewManger("memory", "gosessionid", 3600)
	//启动Session自动GC
	//fmt.Println(globalSessions)
	//go globalSessions.GC()
}

//操作Session
func Count(w http.ResponseWriter, r *http.Request) {
	sess := globalSessions.SessionStart(w, r)
	createtime := sess.Get("createtime")
	if createtime == nil {
		sess.Set("createtime", time.Now().Unix())
	} else if (createtime.(int64) + 360) < (time.Now().Unix()) {
		//超时重置
		globalSessions.SessionDestroy(w, r)
		sess = globalSessions.SessionStart(w, r)
	}
	ct := sess.Get("countnum")
	if ct == nil {
		sess.Set("countnum", 1)
	} else {
		sess.Set("countnum", ct.(int)+1)
	}
	fmt.Fprint(w, sess.Get("countnum"))
}
