package ah

import (
	"ops/pkg/ctl"
	"ops/pkg/al"
	"ops/pkg/auth"
	"ops/pkg/cfg"
	"ops/pkg/db"
	"time"

	"github.com/gin-gonic/gin"
)

func Restart(c *gin.Context) {
	var (
		infoo  []al.AppList1
		infon  []al.AppList
		res    string
		rest   db.AppResult
		err    error
		app    al.App
		script string
		srv    al.Srvs
		ch     = make(chan db.AppResult)
	)
	ctl.Debug(c.Request)
	defer func() {
		if err != nil {

			res += err.Error()
			ctl.Log.Error(err)
		} else {
			ctl.Log.Debug(err)
		}
		c.Data(200, "", []byte(res))

	}()
	err = c.BindJSON(&infoo)
	if err != nil {
		return
	}
	infon = ParseAppList(infoo)
	notUsed, err := CheckNotUsed(infon, false)
	if err != nil {
		return
	}

	if !notUsed {
		err = ctl.Errorf("获取选中应用服状态异常")
		return
	}
	_, err = db.Db.SQL(`select * from app a,app_fun b where a.fun_id=b.fun_id and a.app_id=? and b.app_zone=?`, infon[0].AppId, auth.GetState(c).State).Get(&app)
	if err != nil {
		return
	}

	if app.AppId == 0 {
		err = ctl.Errorf("获取数据库应用异常")
		return
	}
	if app.RestartScriptPath == "" {
		err = ctl.Errorf("数据库中应用重启脚本配置为空")
		return
	}
	for _, v := range infon {
		db.Db.Where("app_server_id=?", v.AppServerId).Find(&srv)
	}
	ctl.Debug(srv)
	if len(srv) == 0 || srv[0].AppServerId == 0 {
		err = ctl.Errorf("获取数据库应用服异常")
		return
	}

	ctl.Debug(script)
	for _, v := range srv {
		//	rest, err = cfg.GameHosts[v.HostId].RunCmd("source /etc/profile && " + app.ParseScript(app.RestartScriptPath))

		go cfg.GameHosts[v.HostId].SshAppChCmd(ch, v.ShAppServer(), "source /etc/profile && "+app.ParseScript(app.RestartScriptPath))
		/*	if err != nil {
				if cfg.GameHosts[v.HostId] != nil {
					rest = "[" + cfg.GameHosts[v.HostId].HostName + "]" + rest
				}
				return
			}
			res += "[" + cfg.GameHosts[v.HostId].HostName + "]" + rest + "\n"
			res += "重启" + app.AppName + "成功\n"*/
	}
	time.Sleep(time.Second * time.Duration(app.PortCheckInvl))
	stem := srv.CheckState()
	ctl.Debug(stem)
	for range srv {
		rest = <-ch
		res += rest.Result
		if rest.Err != nil {
			res += rest.Err.Error()
		}
		if _, ok := stem[rest.AppServerId]; !ok {
			ctl.Debug(rest.AppServerId)
			ctl.Debug(rest.AppId)
			res += "错误:重启应用失败,检查应用状态异常\n"
		} else {
			if stem[rest.AppServerId].Error != "" {
				res += "错误:重启应用失败," + stem[rest.AppServerId].Error + "\n"
			} else {
				res += "重启" + app.AppName + "成功 PID:" + ctl.Itoa(stem[rest.AppServerId].Pid) + "\n"
			}
		}
	}

	ctl.Debug(res)

	//res = "状态正常"
	ctl.Debug(infoo)

}
