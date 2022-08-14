package al

import (
	"ops/pkg/ctl"
	"ops/pkg/db"
)

type Srvs []AppServer

func (srvs *Srvs) GetInfo() map[int]AppList {
	return srvs.GetApps().GetInfoM()
}

func (srvs Srvs) GetApps() Apps {
	var (
		apps Apps
	)
	db.Db.In("app_id", srvs.GetAppIds()).Find(&apps)
	ctl.Debug(apps)
	ctl.Debug(srvs.GetAppIds())
	if apps == nil || apps[0].AppId == 0 {
		return nil
	}
	return apps
}

type SrvFun struct {
	AppServerId  int
	AppId        int
	FunId        int
	FunName      string
	ParentFunId  int
	OpsEarlTime  string
	AppZone      int
	SucForce     int
	UpdateState  int
	CheckUsedCmd string
	CheckUsedNot string
	OpsLateTime  string
}

func (srvs Srvs) GetFunM() map[int]AppFun {
	var (
		fun  []SrvFun
		funm = make(map[int]AppFun)
	)
	db.Db.In("a.app_server_id", srvs.GetSrvIds()).SQL(`select a.app_server_id,a.app_id,c.* from app_server a,app b,app_fun c
where a.app_id=b.app_id
and b.fun_id=c.fun_id`).Find(&fun)
	for _, v := range fun {
		funm[v.AppServerId] = AppFun{
			FunId:        v.FunId,
			FunName:      v.FunName,
			OpsEarlTime:  v.OpsEarlTime,
			OpsLateTime:  v.OpsLateTime,
			AppZone:      v.AppZone,
			SucForce:     v.SucForce,
			UpdateState:  v.UpdateState,
			CheckUsedCmd: v.CheckUsedCmd,
			CheckUsedNot: v.CheckUsedNot,
		}
	}
	return funm
}

func (srvs Srvs) GetSrvIds() []int {
	var (
		ids []int
	)
	for _, v := range srvs {
		ids = append(ids, v.AppServerId)
	}
	return ids
}

func (srvs Srvs) GetAppIds() []int {
	var (
		ids []int
		idm = make(map[int]int)
	)
	for _, v := range srvs {
		idm[v.AppId] = 0
	}
	for _, v := range srvs {
		if idm[v.AppId] != 0 {
			continue
		}
		ids = append(ids, v.AppId)
		idm[v.AppId] = v.AppServerId
	}
	return ids
}

type SrvState struct {
	Error string
	Pid   int
}

func (srvs Srvs) CheckState() map[int]*SrvState {
	ctl.Debug("abc")
	var (
		funm map[int]AppFun
		stem = make(map[int]*SrvState)
	)
	ctl.Debug(srvs)
	list := srvs.GetInfo()
	funm = srvs.GetFunM()
	ctl.Debug(list)
	for _, v := range srvs {
		stem[v.AppServerId] = &SrvState{}
		if _, ok := list[v.AppServerId]; !ok {
			stem[v.AppServerId].Error = "获取应用信息异常"
			return stem
		}
		if _, ok := funm[v.AppServerId]; !ok {
			stem[v.AppServerId].Error = "获取应用功能异常"
			return stem
		}
		if funm[v.AppServerId].SucForce == 0 {
			if !(list[v.AppServerId].Pid != 0 && list[v.AppServerId].RunState < 2) {
				if list[v.AppServerId].Pid == 0 {
					stem[v.AppServerId].Error += "应用PID不存在,"
				}
				if list[v.AppServerId].RunState < 2 {
					stem[v.AppServerId].Error += "应用状态异常,状态为:" + ParseStateStr(list[v.AppServerId].RunState)
				}
			}
		}
		stem[v.AppServerId].Pid = list[v.AppServerId].Pid
	}
	return stem
}
