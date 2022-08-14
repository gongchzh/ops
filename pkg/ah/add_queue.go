package ah

import (
	"ops/pkg/ctl"
	"ops/pkg/al"
	"ops/pkg/auth"
	"ops/pkg/cfg"
	"ops/pkg/db"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type AddQueueData struct {
	Fun        []al.AppFun
	UpEarlTime string
	UpLateTime string
}

func AddQueue(c *gin.Context) {
	var (
		fun AddQueueData
	)
	db.Db.Where("app_zone=?", auth.GetState(c).State).Find(&fun.Fun)
	ctl.Debug(fun)
	fun.UpEarlTime = time.Now().Format("2006-01-02T15:04")
	fun.UpLateTime = time.Now().Add(time.Hour * 4).Format("2006-01-02T15:04")
	c.HTML(200, "plt_add_queue.html", &fun)
}

func AddQueueForm(c *gin.Context) {
	var (
		que  al.AppUpdateQueue
		res  string
		err  error
		quet al.AppUpdateQueue
		fun  al.AppFun
	)
	defer func() {
		if err != nil {
			res = err.Error()
			ctl.Log.Error(err)
		} else {
			ctl.Log.Debug(err)
		}
		if err != nil {
			ctl.Log.Error(err)
		} else {
			ctl.Log.Debug("发送切服钉钉消息成功")
		}
		c.Data(200, "", []byte(res))

	}()
	ctl.Debug(c.PostForm("FunId"))
	que.FunId, _ = strconv.Atoi(c.PostForm("FunId"))
	que.Info = c.PostForm("Info")
	que.Md5 = c.PostForm("Md5")
	que.UpEarlTime = c.PostForm("UpEarlTime")
	que.UpEarlTime = strings.Replace(que.UpEarlTime, "T", " ", -1)
	que.UpLateTime = c.PostForm("UpLateTime")
	que.UpLateTime = strings.Replace(que.UpLateTime, "T", " ", -1)
	ctl.Debug(que)
	que.Status = 3
	if que.FunId == 0 {
		err = ctl.Errorf("应用功能ID为空")
		return
	}
	db.Db.Where("fun_id=?", que.FunId).Get(&fun)
	if fun.FunId == 0 {
		err = ctl.Errorf("获取数据库应用功能ID异常")
		return
	}
	if len(que.Md5) != 32 {
		err = ctl.Errorf("MD5长度不对")
		return
	}

	ctl.Debug(cfg.RegMinute.FindString(que.UpLateTime))
	ctl.Debug(cfg.RegMinute.FindString(que.UpEarlTime + ":00"))
	if cfg.RegMinute.MatchString(que.UpEarlTime) {
		que.UpEarlTime = que.UpEarlTime + ":00"
	}
	if cfg.RegMinute.MatchString(que.UpLateTime) {
		que.UpLateTime = que.UpLateTime + ":00"
	}
	if que.UpEarlTime == "" {
		err = ctl.Errorf("最早更新时间为空")
		return
	}
	if que.UpLateTime == "" {
		err = ctl.Errorf("最迟更新时间为空")
		return
	}
	//	db.Db.Where("up_earl_time>=? and up_late_time<=?",que.UpEarlTime,que.UpEarlTime).Or("up_earl_time>=? and up_late_time<=?")
	db.Db.Where("((up_earl_time>=? and up_earl_time<=?) or (up_late_time>=? and up_late_time<=?)) and fun_id=? and md5=?", que.UpEarlTime, que.UpLateTime, que.UpEarlTime, que.UpLateTime, que.FunId, que.Md5).Get(&quet)
	//	db.Db.Where("up_earl_time>=? and up_late_time<=? and fun_id=? and md5=?", que.UpEarlTime, que.UpLateTime, que.FunId, que.Md5).Get(&quet)
	if quet.UpdateQueueId != 0 {
		err = ctl.Errorf("相同时间段队列己存在")
		return
	}
	ctl.Debug(quet)
	_, err = db.Db.InsertOne(&que)
	if err != nil {
		return
	}
	res += "应用功能名:" + fun.FunName + "\n"
	res += "MD5:" + que.Md5 + "\n"
	res += "最早更新时间:" + que.UpEarlTime + "\n"
	res += "最迟更新时间:" + que.UpLateTime + "\n"
	res += "更新内容:" + que.Info + "\n添加更新队列成功"
}
