package op

import (
	"bytes"
	"ops/pkg/ctl"
	"encoding/json"
	"io/ioutil"
	"ops/pkg/al"
	"ops/pkg/cfg"
	"ops/pkg/db"
	"ops/pkg/sh"
	"os"
	"strings"
	"time"

	"github.com/pkg/sftp"
)

func SyncInfo() error {
	err := SyncAppInfo()
	if err != nil {
		ctl.Log.Error(err)
		ctl.Debug(err)
		return err
	}
	err = SyncNginxInfo()
	if err != nil {
		ctl.Log.Error(err)
		ctl.Debug(err)
		return err
	}
	return err

}
func SyncNginxInfo() error {
	var (
		info       []al.NginxInfo
		infom      = make(map[int][]al.NginxInfo)
		srv        []al.AppNginxServer
		infot      al.NginxInfo
		appNginx   []al.AppNginx
		appNginxm  = make(map[int]al.AppNginx)
		err        error
		app        []al.App
		appm       = make(map[int][]al.App)
		appServer  []al.AppServer
		appServerm = make(map[int][]al.AppServer)
	)

	db.Db.Find(&appNginx)
	for _, v := range appNginx {
		appNginxm[v.AppNginxId] = v
	}
	db.Db.Where("app_nginx_id!=0").Find(&app)
	for _, v := range app {
		appm[v.AppNginxId] = append(appm[v.AppNginxId], v)
	}
	db.Db.Find(&appServer)
	ctl.Debug(appServer)
	for _, v := range appServer {
		appServerm[v.AppId] = append(appServerm[v.AppId], v)
	}
	ctl.Debug(appServerm)
	ctl.Debug(app)
	ctl.Debug(appm)
	//	ctl.Debug(appm[2])
	//	ctl.Debug(appServerm[1])
	//	ctl.Debug(appServerm[1][0])
	//	ctl.Debug(appServerm[1][0].AppServerLocal)
	for k, v := range appm {
		ctl.Debug(k, v, v[0].AppId, v[0].AppName)

	}
	db.Db.Find(&srv)
	for _, v := range srv {
		if _, ok := appNginxm[v.AppNginxId]; !ok {
			continue
		}
		if _, ok := appm[v.AppNginxId]; !ok {
			continue
		}
		infot = al.NginxInfo{}
		infot.AppNginxServerId = v.AppNginxServerId
		infot.AppNginx = appNginxm[v.AppNginxId]
		infot.HostId = v.HostId
		for _, t := range appm[v.AppNginxId] {
			ctl.Debug(t.AppId)
			ctl.Debug(t.AppName)
			ctl.Debug(appServerm[t.AppId])
			for _, y := range appServerm[t.AppId] {

				infot.ServerInfo = append(infot.ServerInfo, al.NginxServerInfo{AppName: t.AppName, AppServerId: y.AppServerId, AppServerLocal: y.AppServerLocal})
				ctl.Debug(infot.ServerInfo)
				ctl.Debug(y.AppServerId)
				ctl.Debug(y.AppServerLocal)
				ctl.Debug(t.AppId)
				ctl.Debug(y.AppId, y.AppServerId)
			}
		}
		info = append(info, infot)
	}
	for _, v := range info {
		infom[v.HostId] = append(infom[v.HostId], v)
	}
	for _, v := range infom {
		b, err := json.Marshal(v)
		if err != nil {
			ctl.Debug(err)
			ctl.Log.Error(err)
		}
		ctl.Debug(string(b))
		ctl.Debug(cfg.GameHosts)
		ctl.Debug(v)
		ctl.Debug(v[0])
		ctl.Debug(v[0].HostId)

		_, err = cfg.GameHosts[v[0].HostId].WriteFile(cfg.App.NginxInfo, b)
		ctl.Debug(cfg.App.NginxInfo)
		if err != nil {
			ctl.Log.Error(err)
			ctl.Debug(err)
			return err
		}
	}
	return err

}

func SyncAgent() error {
	var (
		err error
	)
	ctl.Debug("start sync2")
	fd, err := ioutil.ReadFile("d:/gotest/src/OpsAgent/OpsAgent")
	ctl.Debug(err)
	if err != nil {
		return err
	}
	if true {
		ctl.Debug(cfg.GameHosts)
		for _, v := range cfg.GameHosts {
			ctl.Debug(v)
			ln, err := v.WriteFile("/home/app/shell/OpsAgent", fd)

			ctl.Debug(ln, err)
		}
	} else {
		ln, err := cfg.GameHosts[93].WriteFile("/home/app/shell/OpsAgent", fd)
		ctl.Debug(ln, err)
	}
	ctl.Debug(fd)
	return err
}

func SyncAgentFile(fs bytes.Buffer, v *sh.SSHHost, sinfo os.FileInfo, dst string) {
	ctl.Debug(v)
	sess, err := sftp.NewClient(v.SshClient)
	if err != nil {
		ctl.Debug(err)
		ctl.Log.Error(err)
		return
	}
	defer sess.Close()
	dinfo, err := sess.Stat(dst)
	if err != nil {
		if _, err = sess.Stat(ctl.UnixDir(dst)); err != nil {
			err = sess.Mkdir(ctl.UnixDir(dst))
			if err != nil {
				ctl.Debug(err)
				ctl.Log.Error(err)
				return
			}
		}
	}
	if dinfo != nil {
		invl := sinfo.ModTime().Sub(dinfo.ModTime())
		ctl.Debug(invl, invl < time.Second)
		if sinfo.Size() == dinfo.Size() && ((invl > 0 && invl < time.Second) || (invl < 0 && invl > time.Second*-1)) {
			//return
			cmd := v.SshCmd("/home/app/shell/OpsAgent")
			ctl.Debug(cmd)
			if strings.Contains(cmd, "Text file busy") {
				sess.Remove("/home/app/shell/OpsAgent")
				//time.Sleep(time.Second * 10)
				//	ft, err := sess.Open("/home/app/shell/OpsAgent")
				ctl.Debug(err)

				//	ft.Close()
				//	return
				//return
			} else {

				return
			}
		}
	}

	fd, err := sess.Create("/home/app/shell/OpsAgent")
	if err != nil {
		ctl.Debug(err)
		ctl.Log.Error(err)
		return
	}
	defer fd.Close()
	ln, err := fd.Write(fs.Bytes())
	if err != nil {
		ctl.Debug(err)
		ctl.Log.Error(err)
		return
	}
	if int64(ln) != sinfo.Size() {
		err = ctl.Errorf(v.HostName+":复制文件大小不一致,源大小%s,目标:%s", sinfo.Size(), ln)
		ctl.Debug(err)
		ctl.Log.Error(err)
		return
	}
	ctl.Log.Debug(v.HostName + "上传OpsAgent成功")
	ctl.Debug(v.HostName + "上传OpsAgent成功")
	err = sess.Chtimes("/home/app/shell/OpsAgent", sinfo.ModTime(), sinfo.ModTime())
	if err != nil {
		ctl.Debug(err)
		ctl.Log.Error(err)
		return
	}
	ctl.Log.Debug(v.HostName + "修改OpsAgent时间成功")
	ctl.Debug(v.HostName + "修改OpsAgent时间成功")
	err = sess.Chmod("/home/app/shell/OpsAgent", 0744)
	if err != nil {
		ctl.Debug(err)
		ctl.Log.Error(err)
		return
	}
	ctl.Log.Debug(v.HostName + "修改OpsAgent权限成功")
	ctl.Debug(v.HostName + "修改OpsAgent权限成功")
	fd.Close()
}
func SyncAgent1() error {
	var (
		err error
		fs1 bytes.Buffer
	)
	ctl.Debug("start sync agent1")
	fs, err := os.Open("d:/gotest/src/OpsAgent/OpsAgent")
	if err != nil {
		return err
	}
	fs1.ReadFrom(fs)
	defer fs.Close()
	sinfo, err := os.Stat("d:/gotest/src/OpsAgent/OpsAgent")
	if err != nil {
		return err
	}
	dst := "/home/app/shell/OpsAgent"
	for _, v := range cfg.GameHosts {
		go SyncAgentFile(fs1, v, sinfo, dst)
		time.Sleep(time.Second * 5)
	}
	return err
}
func SyncAppInfo() error {
	ctl.Debug(cfg.GameHosts)

	var (
		app  []al.AppInfo
		apps = make(map[int][]al.AppInfo)
		err  error
	)
	db.Db.SQL(`select c.*,a.app_server_id,b.host_id,b.host_name,d.check_used_cmd,d.check_used_not from app_server a,host b,app c,app_fun d
where a.host_id=b.host_id
and a.app_id=c.app_id
and c.fun_id=d.fun_id
and b.state!=0
and c.status!=0`).Find(&app)
	for _, v := range app {
		apps[v.HostId] = append(apps[v.HostId], v)
	}
	for k, v := range apps {
		ctl.Debug(k, v)
	}
	ctl.Debug(len(apps))
	for k, v := range apps {
		//	ctl.Debug(k)

		b, err := json.Marshal(v)
		if err != nil {
			ctl.Debug(err)
			ctl.Log.Error(err)
			return err
		}
		//	ctl.Debug(k, string(b))
		if _, ok := cfg.GameHosts[k]; ok {
			_, err = cfg.GameHosts[k].WriteFile(cfg.App.Info, b)
			ctl.Debug(err)
			if err != nil {
				ctl.Log.Error(err)
				ctl.Debug(err)
				return err
			}
			//	ctl.Debug(cfg.GameHosts[k].HostName, err)
		}
	}
	ctl.Debug("stop sync info")
	return err
}
