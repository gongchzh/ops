package agt

import (
	"ctl"
	"encoding/json"
	"ops/pkg/al"
	"os"
	"strconv"
	"strings"
)

func AppGetInfo() {
	var (
		apps []AppInfo
		info []AppList
		ch   = make(chan AppList)
	)
	apps, err := GetApps()
	ctl.FatalErr(err)
	for k, _ := range apps {
		go apps[k].GetSrvInfo(ch)
	}
	for range apps {
		info = append(info, <-ch)
	}
	res, err := json.Marshal(&info)
	ctl.FatalErr(err)
	os.Stdout.Write(res)
}

func Contains(str string, arg ...string) string {
	if strings.Contains(str, arg[0]) {
		if len(arg) > 1 {
			return Contains(str, arg[1:]...)
		}
		return str
	}
	return ""
}

func FindString(ps string, arg ...string) string {
	var (
		res  string
		find string
	)
	for _, v := range strings.Split(ps, "\n") {
		find = Contains(v, arg...)
		if find == "" {
			continue
		}
		if res != "" {
			return "异常"
		}
		res = find
	}
	return res
}

func FindPid(pro string) (int, error) {
	p := strings.Split(pro, " ")
	if len(p) < 3 {
		return 0, ctl.Errorf("获取进程异常")
	}
	for k, v := range p {
		if v == "" {
			continue
		}
		if k == 0 {
			continue
		}
		if !regPro.MatchString(v) {
			return 0, ctl.Errorf("获取进程异常")
		}
		pro = v
		break
	}
	return strconv.Atoi(pro)
}

func (app *AppInfo) GetProcess() int {
	var (
		pro string
		err error
	)
	ps, err := ctl.Run("ps", "aux")
	if err != nil {
		return 0
	}
	switch app.AppTypeId {
	case 1:
		pro = FindString(ps, app.BasePath+"/"+app.AppProgram)
	case 2:
		pro = FindString(ps, TomcatDir+ctl.UnixSubDir(app.BasePath), "bin/java")
	case 3:
		pro = FindString(ps, app.AppProgram, "bin/java")
	}
	if pro == "异常" || pro == "" {
		ctl.Printf("获取进程ID异常")
		return 0
	}
	pid, _ := FindPid(pro)
	return pid
}

func (app *AppInfo) GetRunState(list AppList) AppList {
	var (
		pro string
	)
	list.RunState = 5
	defer func() {
		if list.RunState == 1 && app.CheckUsedCmd != "" && app.CheckUsedNot != "" {
			used, err := ctl.Run(app.CheckUsedCmd)
			if err == nil && used == app.CheckUsedNot {
				list.RunState = 0
			}
		}
	}()
	ps, err := ctl.Run("ps", "aux")
	switch app.AppTypeId {
	case 1:
		pro = FindString(ps, app.BasePath+"/"+app.AppProgram)
	case 2:
		pro = FindString(ps, TomcatDir+ctl.UnixSubDir(app.BasePath), "bin/java")
	case 3:
		pro = FindString(ps, app.AppProgram, "bin/java")
	case 4:
		list.RunState = 1
		return list
	default:
		list.RunState = al.AppStatusProcessExpt
		return list
	}
	if pro == "异常" || pro == "" {
		list.RunState = al.AppStatusProcessExpt
		return list
	}
	list.Pid, err = FindPid(pro)
	if err != nil || list.Pid == 0 {
		list.RunState = al.AppStatusProcessExpt
		return list
	}
	netInfo, err := ctl.Run("netstat", "-anpt")
	if err != nil {
		list.RunState = al.AppStatusPortExpt
		return list
	}
	netInfo = FindString(netInfo, strconv.Itoa(app.Port), strconv.Itoa(list.Pid), "LISTEN")
	if netInfo == "异常" || netInfo == "" {
		list.RunState = al.AppStatusPortExpt
		return list
	}
	list.RunState = 1
	return list
}

func (app *AppInfo) GetSrvInfo(ch chan AppList) {
	var (
		list AppList
	)
	list.CurMd5, _ = ctl.Md5sum(app.ProgramDir + "/" + app.AppProgram)
	list.NewMd5, _ = ctl.Md5sum(app.UpdatePath + "/" + app.UpdateProgram)
	infoc, err := os.Stat(app.ProgramDir + "/" + app.AppProgram)
	if err == nil {
		list.UpdateTime = infoc.ModTime().Format(ctl.TimeFormat)
	}
	infon, err := os.Stat(app.UpdatePath + "/" + app.UpdateProgram)
	if err == nil {
		list.NewTime = infon.ModTime().Format(ctl.TimeFormat)
	}
	list.AppId = app.AppId
	list.AppName = app.AppName
	list.AppServerId = app.AppServerId
	list.HostName = app.HostName
	list.HostId = app.HostId
	list.FunId = app.FunId
	list.AppProgram = app.AppProgram
	list.UpdateProgram = app.UpdateProgram
	list.Port = app.Port
	list.AppNginxId = app.AppNginxId
	list = app.GetRunState(list)
	ch <- list
}
