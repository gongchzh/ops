package agt

import (
	"ctl"
	"encoding/json"
	"io/ioutil"
	"ops/pkg/al"
	"os"
	"strings"
)

func AppNginxGetInfo() {
	var (
		info []al.NginxInfo
	)
	b, err := ioutil.ReadFile(NginxInfoFile)
	ctl.FatalErr(err)
	err = json.Unmarshal(b, &info)
	ctl.FatalErr(err)
	for k, v := range info {
		info[k] = CheckNginxInfo(v)
	}
	bs, err := json.Marshal(&info)
	ctl.FatalErr(err)
	os.Stdout.Write(bs)
}

func CheckNginxInfo(info al.NginxInfo) al.NginxInfo {
	b, err := ioutil.ReadFile(NginxUpDir + info.AppNginx.NginxFile)
	ctl.FatalErr(err)
	for k, v := range info.ServerInfo {
		if CheckServerInfo(v, string(b)) {
			info.ServerInfo[k].RunState = 1
		}
	}
	return info
}
func CheckServerInfo(srvInfo al.NginxServerInfo, fd string) bool {
	for _, v := range strings.Split(fd, "\n") {
		if regNgxSrv.MatchString(v) {
			if strings.Contains(v, srvInfo.AppServerLocal) {
				return true
			}
		}
	}
	return false
}
