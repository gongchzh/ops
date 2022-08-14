package op

import (
	"ops/pkg/ctl"
	"ops/pkg/al"
	"ops/pkg/cfg"
	"ops/pkg/db"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type QueueInfo struct {
	AppServerId   int
	AppName       string
	UpdateProgram string
	HostName      string
	CurMd5        string
	Port          int
	Pid           int
	RunState      int
	FunId         int
	UpdateTime    string
	UpdateStatus  int
	NewMd5        string
	NewTime       string
	OpsStatus     int
	UpEarlTime    string
	UpLateTime    string
}
type QueueInfo1 struct {
	AppServerId   string
	AppName       string
	UpdateProgram string
	HostName      string
	CurMd5        string
	Port          string
	Pid           string
	RunState      string
	UpdateTime    string
	UpdateStatus  string
	NewMd5        string
	NewTime       string
	OpsStatus     string
	UpEarlTime    string
	UpLateTime    string
}

type UpdateQueueInfo struct {
	QueueInfo  []QueueInfo1
	AutoSwitch bool
}
type UpdateOpt struct {
	Fun          al.AppFun
	App          al.App
	Server       []al.AppServer
	Queue        al.AppUpdateQueue
	ListM        map[int]al.AppList
	Info         UpdateQueueInfo
	IsSingle     bool
	SrvHost      map[int]int
	Res          string
	UpdateScript string
	//	Result   map[int]*db.Result
}

func (opt *UpdateOpt) ParseScript() {
	opt.UpdateScript = strings.ReplaceAll(opt.App.UpdateScriptPath, "{Md5}", opt.Queue.Md5)
	opt.UpdateScript = strings.ReplaceAll(opt.UpdateScript, "{Port}", ctl.Itoa(opt.App.Port))
	opt.UpdateScript = strings.ReplaceAll(opt.UpdateScript, "{AppName}", opt.App.AppName)
	opt.UpdateScript = strings.ReplaceAll(opt.UpdateScript, "{md5}", opt.Queue.Md5)
	opt.UpdateScript = strings.ReplaceAll(opt.UpdateScript, "{port}", ctl.Itoa(opt.App.Port))
	if opt.Queue.Back == "" {
		opt.Queue.Back = ctl.UnixSubDir(opt.App.BasePath) + ".bak." + ctl.NowMinute()
	}
	opt.UpdateScript = strings.ReplaceAll(opt.UpdateScript, "{Back}", opt.Queue.Back)
}
func (opt *UpdateOpt) Check(c *gin.Context) error {
	var (
		err  error
		sid  []int
		aid  []int
		app  al.Apps
		appt al.App
		srvt al.AppServer
	)

	for _, v := range opt.Info.QueueInfo {
		sidt, err := ctl.Atoi(v.AppServerId)
		if err != nil {
			return err
		}
		if v.UpdateStatus != "未更新" {
			return ctl.Errorf("更新状态非未更新")
		}
		srvt = al.AppServer{}
		appt = al.App{}
		db.Db.Where("app_server_id=?", sidt).Get(&srvt)
		if srvt.AppServerId == 0 {
			return ctl.Errorf("获取应用服异常")
		}
		db.Db.Where("app_id=?", srvt.AppId).Get(&appt)
		if appt.AppId == 0 {
			return ctl.Errorf("获取数据库应用异常")
		}
		opt.App = appt
		db.Db.Where("fun_id=?", opt.App.FunId).Get(&opt.Fun)
		ctl.Debug(al.ParseState(v.RunState))
		ctl.Debug(opt.Fun.UpdateState)
		ctl.Debug(opt.Fun)
		if al.ParseState(v.RunState) > opt.Fun.UpdateState {
			return ctl.Errorf("当前运行状态禁止操作")
		}

		ctl.Debug(opt.App.AppTypeId, opt.App.AppName)
		if al.ParseState(v.RunState) <= 1 && (!cfg.RegInt.MatchString(v.Pid) || v.Pid == "0") {
			return ctl.Errorf("进程异常")
		}
		if len(v.NewMd5) != 32 {
			return ctl.Errorf("新MD5长度不对")
		}
		if !cfg.RegMd5.MatchString(v.NewMd5) {
			return ctl.Errorf("新MD5格式不对")
		}
		if opt.Info.QueueInfo[0].NewMd5 != v.NewMd5 {
			return ctl.Errorf("选择更新的MD5有多个")
		}
		if v.UpdateTime == "" {
			return ctl.Errorf("新包MD5时间异常")
		}
		if v.UpEarlTime == "" {
			return ctl.Errorf("新包最早更新时间异常")
		}
		if v.UpLateTime == "" {
			return ctl.Errorf("新包最迟更新时间异常")
		}
		if v.HostName == "" {
			return ctl.Errorf("主机名异常")
		}
		if v.UpdateProgram == "" {
			return ctl.Errorf("应用源文件异常")
		}
		if len(v.CurMd5) != 32 {
			return ctl.Errorf("当前MD5长度不对")
		}
		if !cfg.RegMd5.MatchString(v.CurMd5) {
			return ctl.Errorf("当前MD5格式不对")
		}
		te, err := time.ParseInLocation(ctl.TimeFormat, v.UpEarlTime, time.Local)
		if err != nil {
			return ctl.Errorf("转换最早时间异常:%s", err)
		}
		tl, err := time.ParseInLocation(ctl.TimeFormat, v.UpLateTime, time.Local)
		if err != nil {
			return ctl.Errorf("转换最迟时间异常:%s", err)
		}
		if !(v.OpsStatus == "维护中" || (time.Now().Sub(te) > time.Second && tl.Sub(time.Now()) > time.Second)) {
			return ctl.Errorf("不在维护或更新时间段")
		}
		sid = append(sid, sidt)
	}
	db.Db.In("app_server_id", sid).Find(&opt.Server)
	if opt.IsSingle && len(opt.Server) != 1 {
		return ctl.Errorf("获取数据库应用服数目不对")
	}
	if len(opt.Info.QueueInfo) != len(opt.Server) {
		return ctl.Errorf("获取数据库应用服数目不对")
	}
	opt.SrvHost = make(map[int]int)
	for _, v := range opt.Server {
		aid = append(aid, v.AppId)
		opt.SrvHost[v.HostId] = v.AppServerId
	}
	db.Db.In("app_id", aid).Find(&app)
	if len(app) != 1 {
		return ctl.Errorf("获取数据库应用数目不对")
	}
	if app[0].AppId == 0 {
		return ctl.Errorf("获取数据库应用异常")
	}

	opt.App = app[0]
	if opt.App.UpdateScriptPath == "" {
		return ctl.Errorf("应用更新脚本配置为空")
	}
	db.Db.Where("md5=? and fun_id=? and status in(1,2,3)", opt.Info.QueueInfo[0].NewMd5, opt.App.FunId).Get(&opt.Queue)
	ctl.Debug(opt.Info.QueueInfo[0].NewMd5, opt.App.FunId)
	if opt.Queue.UpdateQueueId == 0 {
		return ctl.Errorf("获取数据库队列异常")
	}
	opt.ListM = app.GetInfoM()
	for _, v := range opt.Info.QueueInfo {
		server := ctl.AtoiNe(v.AppServerId)
		if _, ok := opt.ListM[server]; !ok {
			return ctl.Errorf("获取应用服信息异常")
		}
		if opt.ListM[server].NewMd5 != v.NewMd5 {
			return ctl.Errorf("应用服新MD5与队列新MD5不一致")
		}
		if v.CurMd5 != opt.ListM[server].CurMd5 {
			return ctl.Errorf("应用服当前MD5与前端当前MD5不一致")
		}
		if v.AppName != opt.ListM[server].AppName {
			return ctl.Errorf("应用服应用名与前端应用名不一致")
		}
		if v.HostName != opt.ListM[server].HostName {
			return ctl.Errorf("应用服主机名与前端主机名不一致")
		}

		if opt.ListM[server].RunState > opt.Fun.UpdateState {
			return ctl.Errorf("当前运行状态禁止操作")
		}
		if opt.ListM[server].NewMd5 == opt.ListM[server].CurMd5 {
			return ctl.Errorf("检验应用服新MD5与当前MD5相同")
		}
		if v.UpdateProgram != opt.ListM[server].UpdateProgram {
			return ctl.Errorf("应用服源文件名与前端源文名不一致")
		}
		if v.UpEarlTime != opt.Queue.UpEarlTime {
			return ctl.Errorf("队列与前端最早更新时间不一致")
		}
		if v.UpLateTime != opt.Queue.UpLateTime {
			return ctl.Errorf("队列与前端最迟更新时间不一致")
		}
		if v.UpdateTime != opt.ListM[server].UpdateTime {
			return ctl.Errorf("应用服与队列当前文件更新时间不一致")
		}
		if CheckOpsStatus(opt.Queue.FunId) && v.OpsStatus != "维护中" {
			return ctl.Errorf("队列与前端维护状态不一致")
		}
	}
	opt.ParseScript()
	return err
}

func CheckOpsStatus(fid int) bool {
	var (
		fun al.AppFun
		ops bool
	)
	db.Db.Where("fun_id=?", fid).Get(&fun)
	if fun.FunId == 0 {
		return false
	}
	if fun.ParentFunId != 0 {
		ops = CheckOpsStatus(fun.ParentFunId)
	}
	if ops {
		return ops
	}
	if fun.OpsEarlTime == "" || fun.OpsLateTime == "" {
		return false
	}
	te, err := time.ParseInLocation(ctl.TimeFormat, fun.OpsEarlTime, time.Local)
	if err != nil {
		return false
	}
	tl, err := time.ParseInLocation(ctl.TimeFormat, fun.OpsLateTime, time.Local)
	if err != nil {
		return false
	}
	if time.Now().Sub(te) > time.Second && tl.Sub(time.Now()) > time.Second {
		return true
	}
	return false

}

func ParseResHtml(res string) string {
	var (
		res1 string
	)
	for _, v := range strings.Split(res, "\n") {
		if strings.Contains(v, "ERROR") {
			res1 += "<p style='color:red;font-size:20px' >"
		} else {
			res1 += "<p style='font-size:20px'>"
		}
		res1 += v + "</p>"
	}
	return res1
}
func (opt *UpdateOpt) Update() error {
	var (
		ch     = make(chan db.Result)
		res    db.Result
		err    error
		sid    []int
		iid    []int
		srv    []al.AppServer
		isFail = make(map[int]bool)
	)

	for _, v := range opt.Server {
		//opt.UpdateScript = strings.ReplaceAll(opt.UpdateScript, "{HostName}", cfg.GameHosts[v.HostId].HostName)
		opt.ParseScript()
		ctl.Debug(opt.App.AppName, v.AppServerId, v.HostId, opt.UpdateScript)
		ctl.Debug(opt.App.AppName, v.AppServerId, v)
		ctl.Log.Debug(opt.App.AppName, v.AppServerId, v.HostId, opt.UpdateScript)
		ctl.Log.Debug(opt.App.AppName, v.AppServerId, v)
		ctl.Debug("source /etc/profile && " + opt.UpdateScript)
		go cfg.GameHosts[v.HostId].SshHostChCmdErr(ch, "source /etc/profile && "+opt.UpdateScript)
		isFail[v.AppServerId] = false
	}
	for range opt.Server {
		res = <-ch
		ctl.Debug(res.HostName, res)
		ctl.Log.Debug(res.HostName, res)
		opt.Res += "[INFO]应用服ID:" + ctl.Itoa(opt.SrvHost[res.HostId]) + " 服务器:" + res.HostName + " 应用名:" + opt.App.AppName + "开始执行升级......\n"
		opt.Res += res.Result
		if res.Err == nil {
			if cfg.RegAppSuc.MatchString(res.Result) {
				opt.Res += "[INFO]应用服ID:" + ctl.Itoa(opt.SrvHost[res.HostId]) + " 服务器:" + res.HostName + " 应用名:" + opt.App.AppName + "升级成功\n"
			} else {
				if opt.Fun.SucForce == 1 {
					isFail[opt.SrvHost[res.HostId]] = true
					err = ctl.ErrorAf("%s[ERROR]%s\n", err, "未检测到升级成功")
					opt.Res += "[ERROR]" + err.Error()
					opt.Res += "[ERROR]应用服ID:" + ctl.Itoa(opt.SrvHost[res.HostId]) + " 服务器:" + res.HostName + " 应用名:" + opt.App.AppName + "升级失败\n"
				} else {
					opt.Res += "[INFO]应用服ID:" + ctl.Itoa(opt.SrvHost[res.HostId]) + " 服务器:" + res.HostName + " 应用名:" + opt.App.AppName + "获取升级状态异常，可能升级成功\n"
				}
			}
		} else {
			isFail[opt.SrvHost[res.HostId]] = true
			err = ctl.ErrorAf("%s[ERROR]%s\n", err, res.Err)
			opt.Res += "[ERROR]" + err.Error()
			opt.Res += "[ERROR]应用服ID:" + ctl.Itoa(opt.SrvHost[res.HostId]) + " 服务器:" + res.HostName + " 应用名:" + opt.App.AppName + "升级失败\n"
		}
	}
	time.Sleep(time.Second * time.Duration(opt.App.PortCheckInvl))
	list := al.Apps{opt.App}.GetInfoM()
	for _, v := range opt.Server {
		if _, ok := list[v.AppServerId]; !ok {
			err = ctl.ErrorAf("%s[ERROR]获取应用服信息异常,应用服ID:%d\n", err, v.AppServerId)
			opt.Res = ctl.Sprintf("%s[ERROR]获取应用服信息异常,应用服ID:%d\n", opt.Res, v.AppServerId)
			continue
		}
		if list[v.AppServerId].CurMd5 != opt.Queue.Md5 {
			err = ctl.ErrorAf("%s[ERROR][检查MD5]检验应用服当前MD5不一致,应用服ID:%d,服务器:%s,MD5:%s\n", err, v.AppServerId, list[v.AppServerId].HostName, opt.Queue.Md5)
			opt.Res = ctl.Sprintf("%s[ERROR][检查MD5]检验应用服当前MD5不一致,应用服ID:%d,服务器:%s,MD5:%s\n", opt.Res, v.AppServerId, list[v.AppServerId].HostName, opt.Queue.Md5)
			isFail[v.HostId] = true
		} else {
			opt.Res += "[INFO][检查MD5]应用服ID:" + ctl.Itoa(v.AppServerId) + " 服务器:" + list[v.AppServerId].HostName + " 应用名:" + opt.App.AppName + " 检查正常\n"
		}
		if list[v.AppServerId].Pid == 0 {
			err = ctl.ErrorAf("%s[ERROR][检查进程ID]检查进程ID异常,应用服ID:%d,服务器:%s,应用名:%s\n", err, v.AppServerId, list[v.AppServerId].HostName, opt.App.AppName)
			opt.Res = ctl.Sprintf("%s[ERROR][检查进程ID]检查进程ID异常,应用服ID:%d,服务器:%s,应用名:%s\n", opt.Res, v.AppServerId, list[v.AppServerId].HostName, opt.App.AppName)
			isFail[v.HostId] = true
		} else {
			opt.Res += "[INFO][检查进程ID]应用服ID:" + ctl.Itoa(v.AppServerId) + " 服务器:" + list[v.AppServerId].HostName + " 应用名:" + opt.App.AppName + " 检查正常,进程ID:" + ctl.Itoa(list[v.AppServerId].Pid) + "\n"
		}
		if list[v.AppServerId].RunState > 1 {
			isFail[v.HostId] = true
			err = ctl.ErrorAf("%s[ERROR][运行状态]检查运行状态异常,应用服ID:%d,服务器:%s,应用名:%s\n", err, v.AppServerId, list[v.AppServerId].HostName, opt.App.AppName)
			opt.Res = ctl.Sprintf("%s[ERROR][运行状态]检查运行状态异常,应用服ID:%d,服务器:%s,应用名:%s\n", opt.Res, v.AppServerId, list[v.AppServerId].HostName, opt.App.AppName)
		}
		if !isFail[v.HostId] {
			opt.Res += "[INFO]应用服ID:" + ctl.Itoa(v.AppServerId) + " 服务器:" + list[v.AppServerId].HostName + " 应用名:" + opt.App.AppName + "升级成功\n"
		}
	}
	for _, v := range list {
		if v.NewMd5 == v.CurMd5 {
			iid = append(iid, v.AppServerId)
		}
	}
	for k, v := range isFail {
		if !v {
			sid = append(sid, k)
		}
	}
	db.Db.Where("app_id=?", opt.App.AppId).Find(&srv)
	//opt.AutoSwitch = true
	ctl.Debug(len(iid) == len(list) && len(list) == len(srv))
	ctl.Debug(len(sid) == len(opt.Server) && err == nil)
	ctl.Debug(len(sid))
	ctl.Debug(len(opt.Server))
	ctl.Debug(isFail)
	if len(iid) == len(list) && len(list) == len(srv) && len(sid) == len(opt.Server) && err == nil {
		if len(iid) == len(sid) {
			opt.Queue.StartTime = ctl.Now()
		}
		opt.Queue.EndTime = ctl.Now()
		opt.Queue.UpdateTime = opt.Queue.EndTime
		opt.Queue.Status = 0
		opt.Res += "[INFO]应用名:" + opt.App.AppName + " md5:" + opt.Queue.Md5 + " 全部升级完毕\n"
		if len(iid) == len(sid) && opt.Info.AutoSwitch {
			res1, err1 := NginxSwitch(opt.App)
			if err1 != nil {
				opt.Res += "[ERROR][自动切服]" + err1.Error() + "\n"
				err = ctl.ErrorAf("%s[ERROR][自动切服]%s\n", err, err1)
			} else {
				opt.Res += "[INFO][自动切服]" + res1
			}
		} else {
			opt.Res += "[ERROR]由于不是更新全部服,并且全部成功,不自动切服\n"
			err = ctl.ErrorAf("%s[ERROR]由于不是更新全部服,并且全部成功,不自动切服\n", err)
		}
	} else if sid != nil {
		opt.Queue.UpdateTime = ctl.Now()
		opt.Queue.Status = 1
		opt.Res += "[ERROR]由于不是更新全部服,并且全部成功,不自动切服\n"
		err = ctl.ErrorAf("%s[ERROR]由于不是更新全部服,并且全部成功,不自动切服\n", err)
	}
	ctl.Debug(iid)
	ctl.Debug(sid)
	ctl.Debug(len(list))
	ctl.Debug(len(opt.Server))
	ctl.Debug(err)
	if sid != nil {
		ln, err1 := db.Db.Where("update_queue_id=? and md5=? and fun_id=?", opt.Queue.UpdateQueueId, opt.Queue.Md5, opt.Queue.FunId).Cols("start_time,update_time,end_time,status,back").Update(&opt.Queue)
		if err1 != nil {
			err = ctl.ErrorAf("%s[ERROR]修改队列信息:%s\n", err, err1)
			opt.Res = ctl.Sprintf("%s[ERROR]修改队列信息:%s\n", opt.Res, err1)
		}
		if ln != 1 {
			err = ctl.ErrorAf("%s[ERROR]:修改队列信息:行数不为1", err)
		} else {
			opt.Res += "[INFO]应用名:" + opt.App.AppName + " md5:" + opt.Queue.Md5 + " 修改数据库更新队列成功"
		}
	}
	ctl.Debug(opt.Res)
	ctl.Log.Debug(opt.Res)
	opt.Res = ParseResHtml(opt.Res)
	return err
}

func NginxSwitch(app al.App) (res string, err error) {
	var (
		ngx          al.AppNginx
		srv          []al.AppNginxServer
		rest, script string
	)
	if app.AppNginxId == 0 {
		err = ctl.Errorf("应用nginxID配置为空")
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
				err = ctl.Errorf("[" + cfg.GameHosts[v.HostId].HostName + "]" + rest + err.Error())
			}
			return
		}
		res += "[" + cfg.GameHosts[v.HostId].HostName + "]" + rest + "\n"
		res += "[" + cfg.GameHosts[v.HostId].HostName + "]切到" + app.AppName + "成功\n"
	}
	return
}
