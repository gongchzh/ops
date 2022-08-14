package ah

import (
	"ops/pkg/ctl"
	"ops/pkg/al"
	"ops/pkg/auth"
	"ops/pkg/db"
	"ops/pkg/op"

	"github.com/gin-gonic/gin"
)

type ServerSetData struct {
	AppServerId    string
	AppName        string
	HostName       string
	AppServerLocal string
}

func SetServerData(c *gin.Context) {
	var (
		info []ServerSetData
	)
	db.Db.SQL(`select a.app_server_id,b.app_name,c.host_name,a.app_server_local from app_server a,app b,host c,app_fun d
where a.app_id=b.app_id
and a.host_id=c.host_id
and b.fun_id=d.fun_id
and b.status=2
and d.app_zone=?
and c.state!=0`, auth.GetState(c).State).Find(&info)
	ctl.Debug(info)
	c.JSON(200, &info)
	ctl.Debug("abc")
}

func SetServerAdd(c *gin.Context) {
	ctl.Debug(c)
	var (
		info ServerSetData
		res  string
		err  error
		app  al.App
		host db.Host
		srv  al.AppServer
	)
	defer func() {
		if err != nil {
			res = err.Error()
			ctl.Log.Error(err)
		} else {
			ctl.Log.Debug(err)
		}
		c.Data(200, "", []byte(res))
	}()
	info.AppName = c.PostForm("AppName")
	info.HostName = c.PostForm("HostName")
	info.AppServerLocal = c.PostForm("AppServerLocal")
	if info.AppName == "" {
		err = ctl.Errorf("应用名为空")
		return
	} else {
		db.Db.Where("app_name=? and status=2", info.AppName).Get(&app)
		if app.AppId == 0 {
			err = ctl.Errorf("应用名不存在")
			return
		}
	}
	if info.HostName == "" {
		err = ctl.Errorf("主机名为空")
		return
	} else {
		db.Db.Where("host_name=? and state!=0", info.HostName).Get(&host)
		if host.HostId == 0 {
			err = ctl.Errorf("主机名不存在")
			return
		}
	}
	srv.AppId = app.AppId
	srv.AppServerLocal = info.AppServerLocal
	srv.HostId = host.HostId
	ln, err := db.Db.InsertOne(&srv)
	if err != nil {
		return
	}
	if ln != 1 {
		err = ctl.Errorf("错误:插入的行不为1,插入的行为:%d", ln)
		return
	}
	res = "插入" + info.AppName + " " + info.HostName + "成功"
	go func() { op.SyncInfo() }()
}

func SetServerEdit(c *gin.Context) {
	ctl.Debug(c)
	var (
		info ServerSetData
		res  string
		srv  al.AppServer
		diff string
		app  al.App
		host db.Host
		err  error
	)
	defer func() {
		if err != nil {
			ctl.Debug(info)
			res = err.Error()
			ctl.Log.Error(err)
		} else {
			ctl.Log.Debug(err)
		}
		c.Data(200, "", []byte(res))
	}()
	info.AppServerId = c.PostForm("AppServerId")
	info.AppName = c.PostForm("AppName")
	info.HostName = c.PostForm("HostName")
	info.AppServerLocal = c.PostForm("AppServerLocal")
	srv.AppServerId = ctl.AtoiNe(info.AppServerId)
	if srv.AppServerId == 0 {
		err = ctl.Errorf("应用服ID为空")
		return
	} else {
		srv = al.AppServer{}
		ctl.Debug(info.AppServerId)
		db.Db.Where("app_server_id=?", info.AppServerId).Get(&srv)
		if srv.AppServerId == 0 {
			err = ctl.Errorf("应用服不存在")
			return
		}
	}
	if info.AppName == "" {
		err = ctl.Errorf("应用名为空")
		return
	} else {
		db.Db.Where("app_name=? and status=2", info.AppName).Get(&app)
		if app.AppId == 0 {
			err = ctl.Errorf("应用名不存在")
			return
		}
	}

	if app.AppId != srv.AppId {
		diff += "应用名,"
	}
	if info.HostName == "" {
		err = ctl.Errorf("主机名为空")
		return
	} else {
		db.Db.Where("host_name=? and state!=0", info.HostName).Get(&host)
		if host.HostId == 0 {
			err = ctl.Errorf("主机名不存在")
			return
		}
	}
	if host.HostId != srv.HostId {
		diff += "主机名,"
	}
	if info.AppServerLocal != srv.AppServerLocal {
		diff += "内网地址,"
	}
	srv.AppId = app.AppId
	srv.AppServerLocal = info.AppServerLocal
	srv.HostId = host.HostId
	if diff == "" {
		res = "未作任何修改"
		return
	}
	ln, err := db.Db.Where("app_server_id=?", srv.AppServerId).Update(&srv)
	if err != nil {
		return
	}
	if ln != 1 {
		err = ctl.Errorf("错误:修改的行不为1,插入的行为:%d", ln)
		return
	}
	res = "修改" + info.AppName + " " + info.AppServerId + "成功\n修改项:" + diff[:len(diff)-1]
	go func() { op.SyncInfo() }()
}
