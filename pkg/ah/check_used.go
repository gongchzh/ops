package ah

import (
	"ops/pkg/ctl"
	"ops/pkg/al"
	"ops/pkg/db"
	"strconv"
)

func CheckNotUsed(list []al.AppList, srvNum bool) (bool, error) {
	var (
		aid  []int
		aidm = make(map[int]int)
		app  al.Apps
		sid  []int
		fun  al.AppFun
		srv  []al.AppServer
	)
	if list == nil {
		return false, ctl.Errorf("获取应用服列表为空")
	}

	for _, v := range list {
		fun = al.AppFun{}

		db.Db.Where("fun_id=?", v.FunId).Get(&fun)
		if fun.FunId == 0 {
			ctl.Debug(fun)
			ctl.Debug(v.FunId)
			return false, ctl.Errorf("获取应用功能数据异常")
		}
		if v.RunState > fun.UpdateState {
			return false, ctl.Errorf(strconv.Itoa(v.AppServerId) + v.AppName + "当前运行状态禁止操作")
		}
		aidm[v.AppId] = 0
		sid = append(sid, v.AppServerId)
	}
	for k, _ := range aidm {
		aid = append(aid, k)
	}
	if len(aid) != 1 {
		ctl.Debug(len(aid))
		ctl.Debug(list[0])
		return false, ctl.Errorf(strconv.Itoa(list[0].AppServerId) + list[0].AppName + "选中的应用数不为1")
	}
	db.Db.In("app_id", aid).Find(&app)
	ctl.Debug(app)
	ctl.Debug(len(app))
	ctl.Debug(aid)
	ctl.Debug(aidm)
	if len(app) != 1 || app[0].AppId == 0 {
		return false, ctl.Errorf(strconv.Itoa(list[0].AppServerId) + list[0].AppName + "数据库获取应用失败")
	}

	//db.Db.In("app_server_id", sid).Find(&srv)
	db.Db.Where("app_id=?", app[0].AppId).Find(&srv)
	ctl.Debug(len(sid))
	ctl.Debug(sid)
	ctl.Debug(srv)
	if srvNum && len(sid) != len(srv) {
		return false, ctl.Errorf(strconv.Itoa(list[0].AppServerId) + app[0].AppName + "数据库获取应用服数目不对")
	}
	if app[0].AppNginxId == 0 {
		//	return false, ctl.Errorf(strconv.Itoa(list[0].AppServerId) + app[0].AppName + "数据库获取nginxId失败")
		return true, nil
	}
	info := app.GetNginxInfo()
	ctl.Debug(info)

	for k, v := range info {
		ctl.Debug(k, v.AppName, v.AppServerId, v.AppServerLocal, v.RunState)
	}
	for _, v := range list {
		if _, ok := info[v.AppServerId]; !ok {
			return false, ctl.Errorf(strconv.Itoa(v.AppServerId) + v.AppName + "获取应用状态失败")
		}
		ctl.Debug(info[v.AppServerId].AppName, v.AppServerId, info[v.AppServerId].AppServerId, v.RunState, info[v.AppServerId].RunState)
		if info[v.AppServerId].RunState != 0 {
			return false, ctl.Errorf(strconv.Itoa(v.AppServerId) + v.AppName + "获取应用状态不一致")
		}

	}
	return true, nil
}
