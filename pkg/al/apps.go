package al

import (
	"ops/pkg/ctl"
	"encoding/json"
	"ops/pkg/cfg"
	"ops/pkg/db"
	"sort"
)

type Apps []App

func (apps Apps) GetAppIds() []int {
	var (
		ids []int
	)
	for _, v := range apps {
		ids = append(ids, v.AppId)
	}
	return ids
}

func (apps Apps) GetSrvs() Srvs {
	var (
		srvs Srvs
	)
	db.Db.In("app_id", apps.GetAppIds()).Find(&srvs)
	if srvs == nil || srvs[0].AppServerId == 0 {
		return nil
	}
	return srvs
}

func (apps Apps) GetFunIds() []int {
	var (
		ids []int
		idm = make(map[int]int)
	)
	for _, v := range apps {
		idm[v.FunId] = 0
	}
	for _, v := range apps {
		if idm[v.FunId] != 0 {
			continue
		}
		ids = append(ids, v.FunId)
		idm[v.FunId] = v.AppId
	}
	return ids
}

func (apps Apps) GetZone() (int, error) {
	var (
		fun []AppFun
	)
	if apps == nil || apps[0].AppId == 0 {
		return 0, ctl.Errorf("apps为空或有空值")
	}
	db.Db.In("fun_id", apps.GetFunIds()).GroupBy("app_zone").Find(&fun)
	if fun == nil || fun[0].FunId == 0 {
		return 0, ctl.Errorf("appfun为空或有空值")
	}
	if len(fun) != 1 {
		return 0, ctl.Errorf("appfun不为1")
	}
	return fun[0].AppZone, nil
}

func (apps Apps) GetAppM() map[int]App {
	var (
		appm = make(map[int]App)
	)
	for _, v := range apps {
		appm[v.AppId] = v
	}
	return appm
}

func (apps Apps) GetHosts() []int {
	var (
		ids []int
		srv []AppServer
	)
	db.Db.In("app_id", apps.GetAppIds()).GroupBy("host_id").Find(&srv)
	for _, v := range srv {
		ids = append(ids, v.HostId)
	}
	return ids
}

func (apps Apps) GetSrvm() map[int]AppServer {
	var (
		srvm = make(map[int]AppServer)
		srv  []AppServer
	)
	db.Db.In("app_id", apps.GetAppIds()).OrderBy("app_id,app_server_id").Find(&srv)
	for _, v := range srv {
		srvm[v.AppServerId] = v
	}
	return srvm
}

func (apps Apps) GetInfoM() map[int]AppList {
	var (
		infom = make(map[int]AppList)
	)
	for _, v := range apps.GetInfo() {
		infom[v.AppServerId] = v
	}
	return infom
}

func (apps Apps) GetInfo() []AppList {
	var (
		ids     []int
		ch      = make(chan string)
		list    []AppList
		listAll []AppList
		srv     []AppServer
		srvm    = make(map[int]AppServer)
		chn     = make(chan map[int]*NginxServerInfo)
		ngx     map[int]*NginxServerInfo
		err     error
		hid     []int
	)
	ids = apps.GetAppIds()
	srvm = apps.GetSrvm()
	go GetAppNginxInfoCh(chn, apps, ids)
	db.Db.In("app_id", ids).GroupBy("host_id").Find(&srv)
	hid = apps.GetHosts()
	for _, v := range hid {
		ctl.Debug(cfg.App.GetInfo)
		go cfg.GameHosts[v].SshChCmd(ch, cfg.App.GetInfo)
	}
	for range hid {
		list = nil
		err = json.Unmarshal([]byte(<-ch), &list)
		if err != nil {
			continue
		}
		ctl.Debug(list)
		for _, v := range list {
			if _, ok := srvm[v.AppServerId]; !ok {
				continue
			}
			listAll = append(listAll, v)
		}
	}
	ngx = <-chn
	ctl.Debug(ngx)
	if ngx != nil {
		for k, v := range listAll {
			if _, ok := ngx[v.AppServerId]; !ok {
				if listAll[k].RunState < 2 && listAll[k].AppNginxId != 0 {
					listAll[k].RunState = AppStatusNginxExpt
				}
			}
			if ngx[v.AppServerId] == nil {
				if listAll[k].RunState < 2 && listAll[k].AppNginxId != 0 {
					listAll[k].RunState = AppStatusNginxExpt
				}
			} else {
				if listAll[k].RunState < 2 && listAll[k].AppNginxId != 0 {
					listAll[k].RunState = ngx[v.AppServerId].RunState
				}
			}
		}
	} else {
		for k, _ := range listAll {
			if listAll[k].RunState < 2 && listAll[k].AppNginxId != 0 {
				listAll[k].RunState = AppStatusNginxExpt
			}
		}
	}
	sort.Sort(AppListS(listAll))
	ctl.Debug(listAll)
	/*	for _, v := range listAll {
		ctl.Debug(v.AppName, v.AppNginxId, v.RunState, v.HostName)
	}*/
	return listAll
}
