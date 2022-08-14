package ah

import (
	"ops/pkg/ctl"
	"ops/pkg/al"
	"ops/pkg/auth"

	"ops/pkg/db"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type FunInfo struct {
	FunId       int
	FunName     string
	OpsEarlTime string
	AppZone     int
	OpsLateTime string
	NewEarlTime string
	NewLateTime string
	Status      int
}

func SetFunData(c *gin.Context) {
	var (
		fun   []al.AppFun
		info  []FunInfo
		infot FunInfo
	)

	db.Db.Where("app_zone=? and app_zone!=0", auth.GetState(c).State).Find(&fun)
	for k, v := range fun {
		infot = FunInfo{}
		infot.AppZone = v.AppZone
		infot.FunId = v.FunId
		infot.FunName = v.FunName
		infot.OpsEarlTime = v.OpsEarlTime
		infot.OpsLateTime = v.OpsLateTime
		ctl.Debug(v.OpsLateTime)
		ctl.Debug(v.OpsEarlTime)

		te, err := time.ParseInLocation(ctl.TimeFormatT, v.OpsEarlTime, time.Local)
		ctl.Debug(err)

		if err != nil {
			//	fun[k].OpsEarlTime = time.Now().Format("2006-01-02 15:04:05")
			infot.NewEarlTime = time.Now().Format("2006-01-02 15:04:05")
		} else {
			if v.OpsEarlTime == "" {
				//fun[k].OpsEarlTime = time.Now().Format("2006-01-02 15:04:05")
				infot.NewEarlTime = time.Now().Format("2006-01-02 15:04:05")
			} else {
				ctl.Debug(time.Now())
				ctl.Debug(te)
				ctl.Debug(time.Now().Sub(te))
				if time.Now().Sub(te) > time.Hour*5 {
					//fun[k].OpsEarlTime = time.Now().Format("2006-01-02 15:04:05")
					infot.NewEarlTime = time.Now().Format("2006-01-02 15:04:05")
				}
			}
		}
		tt, err := time.ParseInLocation(ctl.TimeFormatT, v.OpsLateTime, time.Local)
		ctl.Debug(ctl.TimeFormat)
		ctl.Debug(v.OpsLateTime)
		ctl.Debug(v.OpsEarlTime)
		if err != nil {
			//	fun[k].OpsLateTime = time.Now().Add(time.Hour * 6).Format("2006-01-02 15:04:05")
			ctl.Debug(err)
			infot.NewLateTime = time.Now().Add(time.Hour * 6).Format("2006-01-02 15:04:05")
		} else {
			if v.OpsEarlTime == "" {
				//fun[k].OpsLateTime = time.Now().Add(time.Hour * 6).Format("2006-01-02 15:04:05")
				infot.NewLateTime = time.Now().Add(time.Hour * 6).Format("2006-01-02 15:04:05")
			} else {
				if tt.Sub(time.Now()) < time.Hour {
					//fun[k].OpsLateTime = time.Now().Add(time.Hour * 6).Format("2006-01-02 15:04:05")
					infot.NewLateTime = time.Now().Add(time.Hour * 6).Format("2006-01-02 15:04:05")
				}
			}
		}
		ctl.Debug(fun[k].OpsEarlTime)
		ctl.Debug(fun[k].OpsLateTime)
		/*if infot.NewEarlTime == "" && infot.NewLateTime == "" {
			infot.Status = 1
		}*/
		ctl.Debug(te)
		ctl.Debug(time.Now().Sub(te))
		ctl.Debug(tt)
		ctl.Debug(tt.Sub(time.Now()))

		if time.Now().Sub(te) > time.Second/10 && tt.Sub(time.Now()) > time.Minute {
			infot.Status = 1
		}
		if infot.NewEarlTime == "" {
			infot.NewEarlTime = v.OpsEarlTime
		}
		if infot.NewLateTime == "" {
			infot.NewLateTime = v.OpsLateTime
		}
		info = append(info, infot)
	}

	c.JSON(200, &info)
	ctl.Debug("abc")
}
func SetFunEdit(c *gin.Context) {
	ctl.Debug(c)
	var (
		info al.AppFun
		res  string
		err  error
	)
	defer func() {
		if err != nil {
			res = err.Error()
			ctl.Log.Error(err)
		} else {
			ctl.Log.Debug(err)
		}
		c.Data(200, "", []byte(res))
	}()
	info.FunId, _ = strconv.Atoi(c.PostForm("FunId"))
	info.FunName = c.PostForm("FunName")
	info.OpsEarlTime = c.PostForm("OpsEarlTime")
	info.OpsLateTime = c.PostForm("OpsLateTime")
	if info.FunId == 0 || info.FunName == "" || info.OpsEarlTime == "" || info.OpsLateTime == "" {
		err = ctl.Errorf("参数有空值")
		return
	}
	//	te, err := time.ParseInLocation(ctl.TimeFormat, v.OpsLateTime, time.Local)
	//	tt, err := time.ParseInLocation(ctl.TimeFormat, v.OpsLateTime, time.Local)
	//	ctl.Debug(tt)
	//	ctl.Debug(te)
	ctl.Debug(info.OpsEarlTime)
	ctl.Debug(info.OpsLateTime)
	ln, err := db.Db.Where("fun_id=? and fun_name=?", info.FunId, info.FunName).Cols("ops_earl_time,ops_late_time").Update(&al.AppFun{FunId: info.FunId, OpsEarlTime: info.OpsEarlTime, OpsLateTime: info.OpsLateTime})
	if err != nil {
		return
	}
	if ln != 1 {
		err = ctl.Errorf("修改维护时间的行数不对,受影响的行 %d", ln)
		return
	}
	res = "修改" + info.FunName + "维护时间成功"
}
