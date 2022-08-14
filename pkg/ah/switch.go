package ah

import (
	"ops/pkg/ctl"
	"ops/pkg/al"
	"ops/pkg/auth"
	"ops/pkg/cfg"
	"ops/pkg/db"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func Switch(c *gin.Context) {
	var (
		infoo  []al.AppList1
		infon  []al.AppList
		res    string
		rest   string
		err    error
		app    al.App
		script string
		ngx    al.AppNginx
		srv    []al.AppNginxServer
	)
	ctl.Debug(c.Request)
	defer func() {
		if err != nil {
			res = err.Error()
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
	notUsed, err := CheckNotUsed(infon, true)
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
	db.Db.Where("app_nginx_id=?", app.AppNginxId).Get(&ngx)
	if ngx.AppNginxId == 0 {
		err = ctl.Errorf("获取数据库nginx异常")
		return
	}
	db.Db.Where("app_nginx_id=?", ngx.AppNginxId).Find(&srv)
	if len(srv) == 0 || srv[0].AppNginxId == 0 {
		err = ctl.Errorf("获取数据库nginx服异常")
		return
	}
	url := c.Request.URL.String()
	if (app.NgxNum == 0 && strings.Contains(url, "switch_sed")) || (app.NgxNum == 1 && strings.Contains(url, "switch_frt")) {
		err = ctl.Errorf("选中的服务器与切服按钮不一致")
		return
	}
	if app.NgxNum == 0 {
		script = ngx.NginxFirstScript
	} else {
		script = ngx.NginxSecondScript
	}
	ctl.Debug(script)
	for _, v := range srv {
		rest, err = cfg.GameHosts[v.HostId].RunCmd(script)
		if err != nil {
			if cfg.GameHosts[v.HostId] != nil {
				rest = "[" + cfg.GameHosts[v.HostId].HostName + "]" + rest
			}
			return
		}
		res += "[" + cfg.GameHosts[v.HostId].HostName + "]" + rest + "\n"

	}
	ctl.Debug(res)
	res += "切到" + app.AppName + "成功"
	//res = "状态正常"
	ctl.Debug(infoo)

}

func ParseAppList(infoo []al.AppList1) []al.AppList {
	var (
		infon []al.AppList
		infot al.AppList
		srv   al.AppServer
		app   al.App
	)
	for _, v := range infoo {
		ctl.Debug(v.AppId)
		ctl.Debug(v.AppName)
		ctl.Debug(v.AppProgram)
		ctl.Debug(v.AppServerId)
		ctl.Debug(v.HostName)
		ctl.Debug(v.Port)
		ctl.Debug(v.CurMd5)
		ctl.Debug(v.HostId)
		ctl.Debug(v.NewTime)
		ctl.Debug(v.UpdateTime)
		ctl.Debug(v.RunState)

		infot = al.AppList{}
		app = al.App{}
		infot.AppServerId, _ = strconv.Atoi(v.AppServerId)
		infot.Port, _ = strconv.Atoi(v.Port)
		infot.AppName = v.AppName
		infot.AppProgram = v.AppProgram
		infot.CurMd5 = v.CurMd5
		infot.HostName = v.HostName
		infot.NewMd5 = v.NewMd5
		infot.NewTime = v.NewTime
		infot.UpdateTime = v.UpdateTime
		infot.UpdateProgram = v.UpdateProgram
		srv = al.AppServer{}
		db.Db.SQL(`select a.* from app_server a,host b
where a.host_id=b.host_id
and a.app_server_id=?
and b.host_name=?`, v.AppServerId, v.HostName).Get(&srv)
		if srv.AppServerId == 0 {
			return nil
		}
		infot.AppId = srv.AppId
		infot.HostId = srv.HostId
		db.Db.Where("app_id=?", infot.AppId).Get(&app)
		infot.FunId = app.FunId
		switch v.RunState {
		case "未切服":
			infot.RunState = 0
		case "运行中":
			infot.RunState = 1
		case "Ngx异常":
			infot.RunState = 2
		case "端口异常":
			infot.RunState = 3
		case "进程异常":
			infot.RunState = 4
		default:
			infot.RunState = 2
		}
		infon = append(infon, infot)
	}

	return infon
}
