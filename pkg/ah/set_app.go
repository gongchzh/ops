package ah

import (
	"ops/pkg/ctl"
	"ops/pkg/al"
	"ops/pkg/auth"
	"ops/pkg/db"
	"ops/pkg/op"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AppSetData struct {
	AppId             int
	AppName           string
	FunName           string
	Port              int
	PortCheckInvl     int
	Status            int
	AppTypeName       string
	AppNginxName      string
	NgxNum            int
	BasePath          string
	BackDir           string
	UpdatePath        string
	UpdateProgram     string
	ProgramDir        string
	AppProgram        string
	LogPath           string
	LogFormat         string
	UpdateScriptPath  string
	RestartScriptPath string
}

func SetAppData(c *gin.Context) {
	var (
		info []AppSetData
	)

	/*db.Db.SQL(`select a.*,b.fun_name,c.app_nginx_name,d.app_type_name from app a,app_fun b,app_nginx c,app_type d
	where a.fun_id=b.fun_id
	and a.app_nginx_id=c.app_nginx_id
	and a.app_type_id=d.app_type_id
	and a.app_zone=?
	`, auth.GetState(c).State).Find(&info)*/

	/*	db.Db.SQL(`select e.*,c.app_nginx_name from  (select a.*,b.fun_name,d.app_type_name from app a,app_fun b,app_type d
		where a.fun_id=b.fun_id

		and a.app_type_id=d.app_type_id
		and b.app_zone=?) e left join app_nginx c
		on e.app_nginx_id=c.app_nginx_id
		`, auth.GetState(c).State).Find(&info)*/
	db.Db.SQL(`select e.*,c.app_nginx_name from  (select a.*,b.fun_name,d.app_type_name from app_type d,app a left join app_fun b 
on a.fun_id=b.fun_id
where 
 a.app_type_id=d.app_type_id
and (b.app_zone=? or a.fun_id=0)) e left join app_nginx c
on e.app_nginx_id=c.app_nginx_id`, auth.GetState(c).State).Find(&info)
	ctl.Debug(info)
	c.JSON(200, &info)
	ctl.Debug("abc")
}

func SetAppAdd(c *gin.Context) {
	ctl.Debug(c)
	var (
		info    AppSetData
		res     string
		err     error
		fun     al.AppFun
		app     al.App
		ngx     al.AppNginx
		appType al.AppType
		appn    al.App
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
	info.AppNginxName = c.PostForm("AppNginxName")
	info.AppProgram = c.PostForm("AppProgram")
	info.AppTypeName = c.PostForm("AppTypeName")
	info.BackDir = c.PostForm("BackDir")
	info.BasePath = c.PostForm("BasePath")
	info.FunName = c.PostForm("FunName")
	info.LogFormat = c.PostForm("LogFormat")
	info.LogPath = c.PostForm("LogPath")
	info.ProgramDir = c.PostForm("ProgramDir")
	info.RestartScriptPath = c.PostForm("RestartScriptPath")
	info.UpdatePath = c.PostForm("UpdatePath")
	info.UpdateProgram = c.PostForm("UpdateProgram")
	info.UpdateScriptPath = c.PostForm("UpdateScriptPath")
	info.Port, _ = strconv.Atoi(c.PostForm("Port"))
	info.PortCheckInvl, _ = strconv.Atoi(c.PostForm("PortCheckInvl"))

	if info.AppName == "" {
		err = ctl.Errorf("应用名为空")
		return
	} else {
		db.Db.Where("app_name=?", info.AppName).Get(&appn)
		if appn.AppId != 0 {
			err = ctl.Errorf("应用名己存在")
			return
		}
	}
	app.AppName = info.AppName
	if info.AppProgram == "" {
		err = ctl.Errorf("应用文件名为空")
		return
	}
	app.AppProgram = info.AppProgram
	if info.AppNginxName != "" {
		db.Db.Where("app_nginx_name=?", info.AppNginxName).Get(&ngx)
		if ngx.AppNginxId == 0 {
			err = ctl.Errorf("nginx不存在")
			return
		}
		app.AppNginxId = ngx.AppNginxId
	}
	if info.AppTypeName == "" {
		err = ctl.Errorf("应用类型为空")
		return
	}
	db.Db.Where("app_type_name=?", info.AppTypeName).Get(&appType)
	if appType.AppTypeId == 0 {
		err = ctl.Errorf("应用类型不存在")
		return
	}
	app.AppTypeId = appType.AppTypeId
	if info.BackDir == "" {
		err = ctl.Errorf("备份目录为空")
		return
	}
	app.BackDir = info.BackDir
	if info.BasePath == "" {
		err = ctl.Errorf("主目录为空")
		return
	}
	app.BasePath = info.BasePath
	if info.FunName == "" {
		err = ctl.Errorf("应用功能名为空")
		return
	}
	db.Db.Where("fun_name=?", info.FunName).Get(&fun)
	if fun.FunId == 0 {
		err = ctl.Errorf("应用功能名不存在")
		return
	}
	app.FunId = fun.FunId

	if info.Port == 0 {
		err = ctl.Errorf("端口为空")
		return
	}
	app.Port = info.Port

	if info.PortCheckInvl == 0 && app.AppTypeId != 4 {
		err = ctl.Errorf("端口检查间隔为空")
		return
	}
	app.PortCheckInvl = info.PortCheckInvl

	if info.ProgramDir == "" {
		err = ctl.Errorf("应用文件目录为空")
		return
	}
	app.ProgramDir = info.ProgramDir

	if info.RestartScriptPath == "" && app.AppTypeId != 4 {
		err = ctl.Errorf("重启脚本为空")
		return
	}
	app.RestartScriptPath = info.RestartScriptPath

	if info.UpdatePath == "" {
		err = ctl.Errorf("更新目录为空")
		return
	}
	app.UpdatePath = info.UpdatePath

	if info.UpdateProgram == "" {
		err = ctl.Errorf("更新文件为空")
		return
	}
	app.UpdateProgram = info.UpdateProgram

	if info.UpdateScriptPath == "" {
		err = ctl.Errorf("更新脚本为空")
		return
	}
	app.UpdateScriptPath = info.UpdateScriptPath

	switch c.PostForm("NgxNum") {
	case "第一":
		info.NgxNum = 0

	case "第二":
		info.NgxNum = 1
	default:
		if app.AppNginxId != 0 {
			err = ctl.Errorf("nginx num 参数错误")
		}

	}
	app.NgxNum = info.NgxNum
	if app.AppNginxId != 0 {
		db.Db.Where("app_nginx_id=? and ngx_num=?", app.AppNginxId, app.NgxNum).Get(&appn)
		if appn.AppId != 0 {
			err = ctl.Errorf("nginx序号己存在")
			return
		}
	}
	app.Status = 2
	ln, err := db.Db.InsertOne(&app)
	if err != nil {
		return
	}
	if ln != 1 {
		err = ctl.Errorf("错误:插入的行不为1,插入的行为:%d", ln)
		return
	}
	res = "插入" + info.AppName + "成功"
	go func() { op.SyncInfo() }()
}

func SetAppEdit(c *gin.Context) {
	ctl.Debug(c)
	var (
		info    AppSetData
		res     string
		err     error
		fun     al.AppFun
		app     al.App
		ngx     al.AppNginx
		appType al.AppType
		appn    []al.App
		appo    al.App
		diff    string
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
	info.AppId, _ = strconv.Atoi(c.PostForm("AppId"))
	info.AppName = c.PostForm("AppName")
	info.AppNginxName = c.PostForm("AppNginxName")
	info.AppProgram = c.PostForm("AppProgram")
	info.AppTypeName = c.PostForm("AppTypeName")
	info.BackDir = c.PostForm("BackDir")
	info.BasePath = c.PostForm("BasePath")
	info.FunName = c.PostForm("FunName")
	info.LogFormat = c.PostForm("LogFormat")
	info.LogPath = c.PostForm("LogPath")
	info.ProgramDir = c.PostForm("ProgramDir")
	info.RestartScriptPath = c.PostForm("RestartScriptPath")
	info.UpdatePath = c.PostForm("UpdatePath")
	info.UpdateProgram = c.PostForm("UpdateProgram")
	info.UpdateScriptPath = c.PostForm("UpdateScriptPath")
	info.Port, _ = strconv.Atoi(c.PostForm("Port"))
	info.PortCheckInvl, _ = strconv.Atoi(c.PostForm("PortCheckInvl"))
	ctl.Debug(info)
	if info.AppId == 0 {
		err = ctl.Errorf("应用ID为空")
		return
	}
	app.AppId = info.AppId
	if info.AppName == "" {
		err = ctl.Errorf("应用名为空")
		return
	}
	app.AppName = info.AppName
	db.Db.Where("app_id=? and app_name=?", app.AppId, app.AppName).Get(&appo)
	if appo.AppId == 0 {
		err = ctl.Errorf("应用不存在")
		return
	}
	app.AppProgram = info.AppProgram
	if app.AppProgram == "" {
		err = ctl.Errorf("应用文件为空")
		return
	}
	if app.AppProgram != appo.AppProgram {
		diff += "应用文件名,"
	}
	if info.AppNginxName != "" {
		db.Db.Where("app_nginx_name=?", info.AppNginxName).Get(&ngx)
		if ngx.AppNginxId == 0 {
			err = ctl.Errorf("nginx不存在")
			return
		}
		app.AppNginxId = ngx.AppNginxId

	}
	if app.AppNginxId != appo.AppNginxId {
		diff += "Nginx名,"
	}
	if info.AppTypeName == "" {
		err = ctl.Errorf("应用类型为空")
		return
	}
	db.Db.Where("app_type_name=?", info.AppTypeName).Get(&appType)
	if appType.AppTypeId == 0 {
		err = ctl.Errorf("应用类型不存在")
		return
	}
	app.AppTypeId = appType.AppTypeId
	if app.AppTypeId != appo.AppTypeId {
		diff += "应用类型,"
	}
	if info.BackDir == "" {
		err = ctl.Errorf("备份目录为空")
		return
	}
	app.BackDir = info.BackDir
	if app.BackDir != appo.BackDir {
		diff += "备份目录,"
	}
	if info.BasePath == "" {
		err = ctl.Errorf("主目录为空")
		return
	}
	app.BasePath = info.BasePath
	if app.BasePath != appo.BasePath {
		diff += "主目录,"
	}
	if info.FunName == "" {
		err = ctl.Errorf("应用功能名为空")
		return
	}
	db.Db.Where("fun_name=?", info.FunName).Get(&fun)
	if fun.FunId == 0 {
		err = ctl.Errorf("应用功能名不存在")
		return
	}
	app.FunId = fun.FunId
	if app.FunId != appo.FunId {
		diff += "应用功能,"
	}
	if info.Port == 0 {
		err = ctl.Errorf("端口为空")
		return
	}
	app.Port = info.Port
	if app.Port != appo.Port {
		diff += "端口,"
	}
	if info.PortCheckInvl == 0 && app.AppTypeId != 4 {
		err = ctl.Errorf("端口检查间隔为空")
		return
	}
	app.PortCheckInvl = info.PortCheckInvl
	if app.PortCheckInvl != appo.PortCheckInvl {
		diff += "端口检查间隔,"
	}
	if info.ProgramDir == "" {
		err = ctl.Errorf("应用文件目录为空")
		return
	}
	app.ProgramDir = info.ProgramDir
	if app.ProgramDir != appo.ProgramDir {
		diff += "应用文件目录,"
	}
	if info.RestartScriptPath == "" && app.AppTypeId != 4 {
		err = ctl.Errorf("重启脚本为空")
		return
	}
	app.RestartScriptPath = info.RestartScriptPath
	if app.RestartScriptPath != appo.RestartScriptPath {
		diff += "重启脚本,"
	}
	if info.UpdatePath == "" {
		err = ctl.Errorf("更新目录为空")
		return
	}
	app.UpdatePath = info.UpdatePath
	if app.UpdatePath != appo.UpdatePath {
		diff += "更新目录,"
	}
	if info.UpdateProgram == "" {
		err = ctl.Errorf("更新文件为空")
		return
	}
	app.UpdateProgram = info.UpdateProgram
	if app.UpdateProgram != appo.UpdateProgram {
		diff += "更新文件,"
	}
	if info.UpdateScriptPath == "" {
		err = ctl.Errorf("更新脚本为空")
		return
	}
	app.UpdateScriptPath = info.UpdateScriptPath
	if app.UpdateScriptPath != appo.UpdateScriptPath {
		diff += "更新脚本,"
	}
	switch c.PostForm("NgxNum") {
	case "第一":
		info.NgxNum = 0

	case "第二":
		info.NgxNum = 1
	default:
		if app.AppNginxId != 0 {
			err = ctl.Errorf("nginx num 参数错误")
		}

	}
	app.NgxNum = info.NgxNum
	if app.NgxNum != appo.NgxNum {
		diff += "Nginx序号,"
	}
	if app.AppNginxId != 0 {
		db.Db.Where("app_nginx_id=?", app.AppNginxId).Get(&appn)
		if len(appn) >= 2 {
			if (appn[0].NgxNum == 0 && appn[1].NgxNum == 1) || (appn[0].NgxNum == 1 && appn[1].NgxNum == 0) {
				err = ctl.Errorf("nginx序号异常")
				return
			}
			if appn[0].AppId == app.AppId && appn[1].NgxNum == app.NgxNum {
				err = ctl.Errorf("nginx序号异常")
				return
			}
			if appn[1].AppId == app.AppId && appn[0].NgxNum == app.NgxNum {
				err = ctl.Errorf("nginx序号异常")
				return
			}

		}
	}
	ctl.Debug(auth.GetState(c))
	switch c.PostForm("Status") {
	case "启用":
		app.Status = 2
	case "禁用":
		app.Status = 0
	default:
		err = ctl.Errorf("状态有异常")
		return
	}
	if app.Status != appo.Status {
		diff += "状态,"
	}
	if diff == "" {
		res = "未作任何修改"
		return
	}
	//	ln, err := db.Db.Where("app_id=? and app_name=?", app.AppId, app.AppName).Update(&app)
	ln, err := db.Db.Where("app_id=? and app_name=?", app.AppId, app.AppName).Cols("port,port_check_invl,status,app_type_id,app_nginx_id,ngx_num,fun_id,base_path,back_dir,update_path,update_program,app_program,program_dir,log_path,log_format,update_script_path,restart_script_path").Update(&app)
	if err != nil {
		return
	}
	if ln != 1 {
		err = ctl.Errorf("错误:修改的行不为1,插入的行为:%d", ln)
		return
	}
	res = "修改" + info.AppName + "成功\n修改项:" + diff[:len(diff)-1]
	go func() { op.SyncInfo() }()
}
