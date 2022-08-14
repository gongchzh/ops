package al

import (
	"ops/pkg/ctl"
	"encoding/json"
	"ops/pkg/cfg"
	"ops/pkg/db"
)

func ParseState(ss string) int {
	switch ss {
	case "未切服":
		return AppStatusNoNginx
	case "运行中":
		return AppStatusRun
	case "Ngx异常":
		return AppStatusNginxExpt
	case "端口异常":
		return AppStatusPortExpt
	case "进程异常":
		return AppStatusProcessExpt
	}
	return AppStatusUnknowExpt
}

func ParseStateStr(ss int) string {
	switch ss {
	case AppStatusNoNginx:
		return "未切服"
	case AppStatusRun:
		return "运行中"
	case AppStatusNginxExpt:
		return "Ngx异常"
	case AppStatusPortExpt:
		return "端口异常"
	case AppStatusProcessExpt:
		return "进程异常"
	}
	return "未知异常"
}

func GetAppNginxInfoCh(ch chan map[int]*NginxServerInfo, apps Apps, ids []int) {
	ch <- apps.GetNginxInfo()
}

func (apps Apps) GetNginxInfo() map[int]*NginxServerInfo {
	var (
		appt     Apps
		info     = make(map[int]*NginxServerInfo)
		err      error
		ngxs     []int
		ngxSrv   []AppNginxServer
		ngxInfos = make(map[int]map[int]NginxServerInfo)
		ch       = make(chan string)
		isOld    = make(map[int]int)
		isNew    = make(map[int]int)
		hostm    = make(map[int]int)
		srvs     Srvs
		srvm     = make(map[int]bool)
		res      string
		ids      []int
	)
	ids = apps.GetAppIds()
	err = db.Db.In("app_id", ids).GroupBy("app_nginx_id").Find(&appt)
	if err != nil {
		ctl.Log.Error(err)
		return info
	}
	ctl.Debug(appt)
	ctl.Debug(appt == nil)

	srvs = apps.GetSrvs()
	for _, v := range srvs {
		info[v.AppServerId] = &NginxServerInfo{}
		ngxInfos[v.AppServerId] = make(map[int]NginxServerInfo)
		srvm[v.AppServerId] = false
	}

	for _, v := range appt {
		ngxs = append(ngxs, v.AppNginxId)
	}

	err = db.Db.In("app_nginx_id", ngxs).Find(&ngxSrv)
	if err != nil {
		ctl.Log.Error(err)
		return info
	}
	for _, v := range ngxSrv {
		isOld[v.AppNginxId] += 1
		hostm[v.HostId] = v.AppNginxId
	}
	ngxSrv = nil
	err = db.Db.In("app_nginx_id", ngxs).GroupBy("host_id").Find(&ngxSrv)
	if err != nil {
		ctl.Log.Error(err)
		return info
	}
	for _, v := range ngxSrv {
		go cfg.GameHosts[v.HostId].SshChCmd(ch, cfg.App.NginxGetInfo)
	}
	for range ngxSrv {
		res = <-ch
		infot, ngxId, hostId := UnmarshalAppNginxInfo(res)
		if infot == nil {
			continue
		}
		for _, v := range ngxId {
			isNew[v] += 1
		}
		for k, v := range infot {
			if _, ok := ngxInfos[v.AppServerId]; !ok {
				continue
			}
			ngxInfos[k][hostId] = v
		}
	}

	for k, v := range ngxInfos {
		for s, t := range v {
			if _, ok := hostm[s]; !ok {

				srvm[k] = true
				continue
			}
			if isOld[hostm[s]] != isNew[hostm[s]] {
				srvm[k] = true
				continue
			}
			if info[k].AppName == "" {
				info[k].AppName = t.AppName
				info[k].AppServerId = t.AppServerId
				info[k].AppServerLocal = t.AppServerLocal
				info[k].RunState = t.RunState
			} else {
				if info[k].RunState != t.RunState {
					info[k].RunState = 2
					srvm[k] = true
				}
			}
		}
	}
	for k, _ := range ngxInfos {
		if srvm[k] {
			info[k].RunState = 2
		}
	}

	return info
}

func UnmarshalAppNginxInfo(res string) (map[int]NginxServerInfo, []int, int) {
	var (
		info  = make(map[int]NginxServerInfo)
		err   error
		infos []NginxInfo
		ngxId []int
	)
	err = json.Unmarshal([]byte(res), &infos)
	if err != nil {
		ctl.Log.Error(res)
		ctl.Log.Error(err)
		return nil, nil, 0
	}
	for _, v := range infos {
		for _, t := range v.ServerInfo {
			if _, ok := info[t.AppServerId]; ok {
				ctl.Log.Error("获取app_server状态重复", t.AppServerId, t.AppName, t.AppServerLocal, t.RunState)
				return nil, nil, 0
			}
			info[t.AppServerId] = t
		}
		ngxId = append(ngxId, v.AppNginx.AppNginxId)
	}
	return info, ngxId, infos[0].HostId
}
