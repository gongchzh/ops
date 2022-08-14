package agt

import (
	"ctl"
	"io/ioutil"
	"ops/pkg/al"
	"os"
	"strings"
)

type AppDisInfo struct {
	App       al.App
	AppServer al.AppServer
	AppFun    al.AppFun
	PS        string
	InsertApp string
}

type AppDisNew struct {
	DisInfo []AppDisInfo
	OldApps []AppInfo
	AppType string
}

func AppDiscovery() {
	var (
		an AppDisNew
	)
	if len(os.Args) != 3 {
		ctl.Panic(ctl.Errorf(""), "输入参数不正确")
	}
	an.OldApps, _ = GetApps()
	an.AppType = os.Args[2]
	an.FindApp(an.GetPS())
	an.CheckFile()
	an.CreateSQL()
}

func (an *AppDisNew) GetPS() string {
	ps, _ := ctl.Run("ps", "aux")
	return ps
}

func (an *AppDisNew) FindApp(ps string) {
	ctl.Debug(an.OldApps)
	for _, v := range strings.Split(ps, "\n") {
		if an.Exclude(v) {
			continue
		}
		switch an.AppType {
		case "JavaTomcat":
			if strings.Contains(v, "/data1/opt/tomcat") && strings.Contains(v, "java") {
				an.DisInfo = append(an.DisInfo, AppDisInfo{PS: v})
			}
		case "JavaJar":
			if strings.Contains(v, "-jar") {
				an.DisInfo = append(an.DisInfo, AppDisInfo{PS: v})
			}
		}
	}
	switch an.AppType {
	case "JavaTomcat":
		ctl.FatalErr(an.FindTomcatApp())
	case "JavaJar":
		ctl.FatalErr(an.FindTomcatApp())
	}
}

func (an *AppDisNew) FindTomcatApp() error {
	var (
		psd string
		err error
	)
	for k, v := range an.DisInfo {
		psd = ""
		ctl.Debug(v)
		for _, t := range regTomcat.FindAllString(v.PS, -1) {
			if psd == "" {
				psd = t
			}
			if psd != t {
				ctl.Debug(psd)
				ctl.Debug(t)
				return ctl.Errorf("获取tomcat应用异常，获取应用路径不一致")
			}
		}
		if !ctl.CheckDir(psd) {
			return ctl.Errorf("获取tomcat应用异常，tomcat路径不存在:%s", psd)
		}
		an.DisInfo[k].App.BasePath = strings.ReplaceAll(psd, "/opt/tomcat/", "/www/")
		if !ctl.CheckDir(an.DisInfo[k].App.BasePath) {
			return ctl.Errorf("获取tomcat应用异常，应用主路径不存在:%s", an.DisInfo[k].App.BasePath)
		}
		an.DisInfo[k].App.BackDir = strings.ReplaceAll(an.DisInfo[k].App.BasePath, "/www/", "/www/back/old/")
		an.DisInfo[k].App.ProgramDir = an.DisInfo[k].App.BasePath
		if ctl.CheckDir("/data1/www/back/test1") {
			an.DisInfo[k].App.UpdatePath = "/data1/www/back/test1"
		} else {
			backs, err := ctl.ListDir("/data1/www/back/")
			if err != nil {
				return err
			}
			for _, t := range backs {
				if strings.Contains(t, "test1") {
					an.DisInfo[k].App.UpdatePath = t
					break
				}
				if !ctl.CheckDir(t) {
					continue
				}
				if strings.Contains(t, "old") {
					continue
				}
				an.DisInfo[k].App.UpdatePath = t
			}
		}
		bases, err := ctl.ListDir(an.DisInfo[k].App.BasePath)
		if err != nil {
			return err
		}
		for _, t := range bases {
			if strings.Contains(t, ".war") {
				an.DisInfo[k].App.AppProgram = ctl.UnixSubDir(t)
				break
			}
		}
		if an.DisInfo[k].App.UpdatePath != "" {
			updates, err := ctl.ListDir(an.DisInfo[k].App.UpdatePath)
			if err != nil {
				return err
			}
			for _, t := range updates {
				if regWar.MatchString(t) && strings.Contains(t, ctl.UnixSubDir(an.DisInfo[k].App.BasePath)) {
					an.DisInfo[k].App.UpdateProgram = ctl.UnixSubDir(t)
					break
				}
			}
			if an.DisInfo[k].App.UpdateProgram == "" {
				for _, t := range updates {
					if regWar.MatchString(t) {
						if len(ctl.UnixSubDir(an.DisInfo[k].App.BasePath)) > 8 && strings.Contains(t, ctl.UnixSubDir(an.DisInfo[k].App.BasePath)[:8]) {
							an.DisInfo[k].App.UpdateProgram = ctl.UnixSubDir(t)
							break
						}
					}
				}
			}
		}
		an.DisInfo[k].App.UpdateScriptPath = "/home/app/shell/OpsAgent --app-update {AppName} {Md5} {Back}"
		an.DisInfo[k].App.RestartScriptPath = "/home/app/shell/OpsAgent --app-restart {AppName}"
		an.DisInfo[k].App.PortCheckInvl = 20
		an.DisInfo[k].App.AppTypeId = 1
		an.DisInfo[k].App.Status = 2
		if ctl.CheckFile(psd + "/conf/server.xml") {
			sf, err := ioutil.ReadFile(psd + "/conf/server.xml")
			if err == nil {
				for _, t := range strings.Split(string(sf), "\n") {
					if strings.Contains(t, "Connector") && strings.Contains(t, "port=") {
						portt := regTomcatPort.FindString(t)
						portt = regNum.FindString(portt)
						an.DisInfo[k].App.Port = ctl.AtoiNe(portt)
					}
				}
			}
		}
	}
	return err
}

func (an *AppDisNew) CheckFile() {
	for k, v := range an.DisInfo {
		if !ctl.CheckDir(v.App.BasePath) {
			an.DisInfo[k].App.BasePath = ""
		}
		if !ctl.CheckDir(v.App.BackDir) {
			an.DisInfo[k].App.BackDir = ""
		}
		if !ctl.CheckDir(v.App.ProgramDir) {
			an.DisInfo[k].App.ProgramDir = ""
		}
		if !ctl.CheckDir(v.App.UpdatePath) {
			an.DisInfo[k].App.UpdatePath = ""
		}
		if !ctl.CheckFile(v.App.ProgramDir + "/" + v.App.AppProgram) {
			an.DisInfo[k].App.AppProgram = ""
		}
		if !ctl.CheckFile(v.App.UpdatePath + "/" + v.App.UpdateProgram) {
			an.DisInfo[k].App.BasePath = ""
		}
	}
}

func (an *AppDisNew) CreateSQL() {
	for k, v := range an.DisInfo {
		an.DisInfo[k].InsertApp = `insert into app(`
		an.DisInfo[k].InsertApp += "`app_name`,`port`,`port_check_invl`,`status`,`app_type_id`,`app_nginx_id`,`ngx_num`,"
		an.DisInfo[k].InsertApp += "`fun_id`,`base_path`,`back_dir`,`update_path`,`update_program`,`app_program`,"
		an.DisInfo[k].InsertApp += "`program_dir`,`log_path`,`log_format`,`update_script_path`,`restart_script_path`)"
		an.DisInfo[k].InsertApp += " values("
		an.DisInfo[k].InsertApp += `"",`
		an.DisInfo[k].InsertApp += ctl.Itoa(v.App.Port) + ","
		an.DisInfo[k].InsertApp += ctl.Itoa(v.App.PortCheckInvl) + ","
		an.DisInfo[k].InsertApp += ctl.Itoa(v.App.Status) + ","
		an.DisInfo[k].InsertApp += ctl.Itoa(v.App.AppTypeId) + ","
		an.DisInfo[k].InsertApp += ctl.Itoa(v.App.AppNginxId) + ","
		an.DisInfo[k].InsertApp += ctl.Itoa(v.App.NgxNum) + ","
		an.DisInfo[k].InsertApp += ctl.Itoa(v.App.FunId) + ","
		an.DisInfo[k].InsertApp += "\"" + v.App.BasePath + "\","
		an.DisInfo[k].InsertApp += "\"" + v.App.BackDir + "\","
		an.DisInfo[k].InsertApp += "\"" + v.App.UpdatePath + "\","
		an.DisInfo[k].InsertApp += "\"" + v.App.UpdateProgram + "\","
		an.DisInfo[k].InsertApp += "\"" + v.App.AppProgram + "\","
		an.DisInfo[k].InsertApp += "\"" + v.App.ProgramDir + "\","
		an.DisInfo[k].InsertApp += "\"" + v.App.UpdateScriptPath + "\","
		an.DisInfo[k].InsertApp += "\"" + v.App.RestartScriptPath + "\");"
		ctl.Printf(an.DisInfo[k].InsertApp)
	}
}

func (an *AppDisNew) Exclude(p string) bool {
	for _, v := range an.OldApps {
		switch {
		case v.AppTypeId == 1:
			if strings.Contains(p, v.BasePath) {
				return true
			}
			if strings.Contains(p, "/data1/opt") {
				ctl.Debug(v.BasePath)
				ctl.Debug(strings.Split(v.LogPath, "/logs/")[0])
			}
			if strings.Contains(p, strings.Split(v.LogPath, "/logs/")[0]) {
				return true
			}
		case v.AppTypeId == 2:
			if strings.Contains(p, v.AppProgram) {
				return true
			}
		}
	}
	return false
}
