package ah

import (
	"ops/pkg/ctl"
	"ops/pkg/al"
	"ops/pkg/auth"
	"ops/pkg/cfg"
	"ops/pkg/db"
	"ops/pkg/op"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type AppQueueData struct {
	UpdateQueueId int
	FunId         int
	AppId         int
	Md5           string
	Status        int
	UpEarlTime    string
	UpLateTime    string
}

type AppFunM struct {
	AppId       int
	FunId       int
	FunName     string
	ParentFunId int
	OpsEarlTime string
	OpsLateTime string
	Status      int
}

func getFunGroup(memsk map[int][]int, gm map[int]int) map[int]int {
	for k, _ := range gm {
		if _, ok := memsk[k]; !ok {
			continue
		}
		for s, t := range memsk[k] {
			if _, ok := gm[t]; !ok {
				gm[t] = 0
				ctl.Debug(t, s)
				gm = getFunGroupM(memsk, t, gm)
			}
		}

	}
	return gm
}

func getFunGroupM(memsk map[int][]int, gid int, gm map[int]int) map[int]int {
	if _, ok := memsk[gid]; !ok {
		return gm
	}
	for _, v := range memsk[gid] {
		if _, ok := gm[v]; !ok {
			gm[v] = 0
			ctl.Debug(gid, v)
			gm = getFunGroupM(memsk, v, gm)
		}
	}
	return gm
}

func QueueData(c *gin.Context) {
	var (
		apps  al.Apps
		appl  []al.AppList
		info  []op.QueueInfo
		aid   []int
		que   []AppQueueData
		infot op.QueueInfo
		fun   []al.AppFun
		fid   []int
		funm  = make(map[int]al.AppFun)
		grp   []al.AppGroup
		grpm  = make(map[int]al.AppGroup)
		mem   []al.AppGroupMember
		mems  []al.AppGroupSubMember
		memsk = make(map[int][]int)
		quem  = make(map[int]AppQueueData)
		fm    = make(map[int]map[int]int)
		fmOps = make(map[int]int)
	)

	db.Db.SQL(`select a.update_queue_id,a.fun_id,c.app_id,a.md5,a.status,a.up_earl_time,a.up_late_time from app_update_queue a,
app c,app_fun d
where a.fun_id=c.fun_id
and a.fun_id=d.fun_id
and a.status in(1,2,3)
and (a.fun_id,a.up_earl_time) in(select fun_id,max(up_earl_time) up_earl_time from app_update_queue 
group by fun_id)
and d.app_zone=?`, auth.GetState(c).State).Find(&que)

	db.Db.SQL(`select a.fun_id,b.app_group_id from app_fun a,app_group b,app_group_member c
where a.fun_id=c.fun_id
and c.app_group_id=b.app_group_id
and a.fun_id in(select a.fun_id from app_update_queue a,app_fun b
where 
a.fun_id=b.fun_id
and a.status in(1,2,3)
and b.app_zone=?
group by a.fun_id)`, auth.GetState(c).State).Find(&mem)
	db.Db.Find(&grp)
	ctl.Debug(que)
	ctl.Debug(grp)
	ctl.Debug(mem)
	for _, v := range grp {
		grpm[v.AppGroupId] = v
	}

	for _, v := range mem {
		if _, ok := fm[v.FunId]; !ok {
			fm[v.FunId] = make(map[int]int)
			fmOps[v.FunId] = 0
			fid = append(fid, v.FunId)
		}
		fm[v.FunId][v.AppGroupId] = 0
	}
	db.Db.Find(&mems)
	db.Db.In("fun_id", fid).Find(&fun)
	ctl.Debug(fmOps)
	for _, v := range fun {
		funm[v.FunId] = v
	}

	for _, v := range mems {
		//memsk[v.SubGroupId] = v.AppGroupId
		memsk[v.SubGroupId] = append(memsk[v.SubGroupId], v.AppGroupId)
	}
	for _, v := range mem {
		fm[v.FunId] = getFunGroup(memsk, fm[v.FunId])
	}
	ctl.Debug(memsk)
	ctl.Debug(fm)
	ctl.Debug(funm)
	for k, _ := range fmOps {
		ctl.Debug(ctl.TimeFormat)
		te, err := time.ParseInLocation(ctl.TimeFormatT, funm[k].OpsEarlTime, time.Local)
		ctl.Debug(funm[k].OpsEarlTime)
		if err != nil {
			continue
		}
		tl, err := time.ParseInLocation(ctl.TimeFormatT, funm[k].OpsLateTime, time.Local)
		ctl.Debug(funm[k].OpsLateTime)
		if err != nil {
			continue
		}
		if time.Now().Sub(te) > time.Second && tl.Sub(time.Now()) > time.Second {
			ctl.Debug(funm[k].FunName)
			fmOps[k] = 1
			continue
		}
		ctl.Debug(fm[k])
		for s, _ := range fm[k] {
			ctl.Debug(s, grpm[s].AppGroupName, grpm[s].OpsEarlTime, grpm[s].OpsLateTime)
			te1, err := time.ParseInLocation(ctl.TimeFormatT, grpm[s].OpsEarlTime, time.Local)
			if err != nil {
				continue
			}
			tl1, err := time.ParseInLocation(ctl.TimeFormatT, grpm[s].OpsLateTime, time.Local)
			if err != nil {
				continue
			}
			ctl.Debug(grpm[s].AppGroupName, te1, time.Now().Sub(te1))
			ctl.Debug(grpm[s].AppGroupName, tl1, tl1.Sub(time.Now()))
			if time.Now().Sub(te1) > time.Second && tl1.Sub(time.Now()) > time.Second {
				fmOps[k] = 1
				ctl.Debug(funm[k].FunName, grpm[s].AppGroupName)
				continue
			}
		}
	}
	for _, v := range que {
		quem[v.AppId] = v
		aid = append(aid, v.AppId)
	}
	ctl.Debug(que)
	ctl.Debug(aid)
	db.Db.In("app_id", aid).OrderBy("app_name").Find(&apps)
	ctl.Debug(apps)
	appl = apps.GetInfo()
	for _, v := range appl {
		if _, ok := quem[v.AppId]; !ok {
			infot.NewMd5 = "应用队列异常"
			continue
		}
		infot = op.QueueInfo{}
		infot.AppName = v.AppName
		infot.AppServerId = v.AppServerId
		infot.CurMd5 = v.CurMd5
		infot.UpdateTime = v.UpdateTime
		infot.HostName = v.HostName
		infot.NewTime = v.NewTime
		infot.Port = v.Port
		infot.RunState = v.RunState
		infot.Pid = v.Pid
		infot.UpdateProgram = v.UpdateProgram
		if v.NewMd5 != quem[v.AppId].Md5 {
			infot.NewMd5 = "应用服新MD5与队列不一致"
		}
		//	ctl.Debug(infot.AppServerId, infot.CurMd5)
		//	ctl.Debug(infot.AppServerId, v.CurMd5)
		if infot.NewMd5 == "" && v.NewMd5 == quem[v.AppId].Md5 {
			infot.NewMd5 = v.NewMd5
		}
		if infot.CurMd5 == infot.NewMd5 {
			infot.UpdateStatus = 0
		} else {
			infot.UpdateStatus = 1
		}
		if len(infot.NewMd5) != 32 {
			infot.UpdateStatus = 3
		}
		infot.UpEarlTime = quem[v.AppId].UpEarlTime
		infot.UpLateTime = quem[v.AppId].UpLateTime
		if _, ok := fmOps[v.FunId]; ok {
			infot.OpsStatus = fmOps[v.FunId]
		}

		info = append(info, infot)
	}
	c.JSON(200, &info)
}

type Template struct {
	Md5      string
	Port     int
	HostName string
	AppName  string
}

func Update(c *gin.Context) {
	var (
		res string
		err error
		opt op.UpdateOpt
	)
	ctl.Debug("start update context")
	cfg.AppUpdateLock.Lock()
	ctl.Debug("start update context1")
	defer func() {
		cfg.AppUpdateLock.Unlock()
		if err != nil {
			//	res = err.Error()
			ctl.Log.Error(err)
		} else {
			ctl.Log.Debug(err)
		}
		if err != nil {
			ctl.Log.Error(err)
		} else {
			ctl.Log.Debug("发送切服钉钉消息成功")
		}
		ctl.Log.Debug(res)
		c.Data(200, "", []byte(res))
	}()
	err = c.BindJSON(&opt.Info)
	ctl.Debug(opt.Info)
	if err != nil {
		res = err.Error()
		return
	}
	ctl.Debug(opt.Info.AutoSwitch)
	url := c.Request.URL.String()
	if strings.Contains(url, "update_single") {
		opt.IsSingle = true
	} else if strings.Contains(url, "update_multi") {
		opt.IsSingle = false
	} else {
		err = ctl.Errorf("参数异常")
		return
	}
	err = opt.Check(c)
	ctl.Debug(err)
	if err != nil {
		res = err.Error()
		return
	}
	ctl.Debug("start update")
	err = opt.Update()
	res = opt.Res
	ctl.Debug(err)
	ctl.Debug(res)
}
