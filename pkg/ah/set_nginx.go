package ah

import (
	"ops/pkg/ctl"
	"ops/pkg/al"
	"ops/pkg/auth"
	"ops/pkg/db"
	"ops/pkg/op"
	"strings"

	"github.com/gin-gonic/gin"
)

type NginxSetData struct {
	AppNginxId        string
	AppNginxName      string
	NginxFirstScript  string
	NginxSecondScript string
	//	NginxViewScript   string
	NginxFile       string
	FunId           int
	FunName         string
	NginxFirstFile  string
	NginxSecondFile string
	Member          string
}

type InfoMember struct {
	AppNginxId       int
	AppNginxServerId int
	HostId           int
	HostName         string
}

func SetNginxData(c *gin.Context) {
	var (
		info  []NginxSetData
		infom = make(map[int]string)
		mem   []InfoMember
	)
	db.Db.SQL(`select b.fun_name,a.* from app_nginx a,app_fun b
where a.fun_id=b.fun_id
and b.app_zone=?`, auth.GetState(c).State).Find(&info)
	ctl.Debug(info)
	db.Db.SQL(`select b.host_name,a.* from app_nginx_server a,host b,app_nginx  c,app_fun d
where a.host_id=b.host_id
and a.app_nginx_id=c.app_nginx_id
and c.fun_id=d.fun_id
and b.state!=0
and d.app_zone=?`, auth.GetState(c).State).Find(&mem)
	for _, v := range mem {
		infom[v.AppNginxId] += v.HostName + " "
	}
	for k, v := range info {
		if _, ok := infom[ctl.AtoiNe(v.AppNginxId)]; ok {
			info[k].Member = infom[ctl.AtoiNe(v.AppNginxId)]
		}
	}
	c.JSON(200, &info)
	ctl.Debug("abc")
}

func SetNgxMember(ms string, nid int) (string, error) {
	var (
		memo     []InfoMember
		host     []db.Host
		memns    []string
		memom    = make(map[int]string)
		memnm    = make(map[int]string)
		add, del []int
		err      error
		srv      []al.AppNginxServer
		res      string
	)
	for _, v := range strings.Split(ms, " ") {
		if v == "" {
			continue
		}
		memns = append(memns, v)
	}
	db.Db.SQL(`select b.host_name,a.* from app_nginx_server a,host b
where a.host_id=b.host_id
and b.state!=0
and a.app_nginx_id=?`, nid).Find(&memo)
	db.Db.In("host_name", memns).Where("state!=0").Find(&host)
	if len(memns) != len(host) {
		return res, ctl.Errorf("核对主机数不一致")
	}
	for _, v := range host {
		memnm[v.HostId] = v.HostName
	}
	ctl.Debug(memns)
	ctl.Debug(host)
	ctl.Debug(memo)
	ctl.Debug(memnm)
	for _, v := range memo {
		memom[v.HostId] = v.HostName
		if _, ok := memnm[v.HostId]; !ok {
			del = append(del, v.HostId)
			res += v.HostName + ","
		}
	}
	if del != nil {
		res = "删除的主机为:" + res
	}
	for k, v := range memnm {
		if _, ok := memom[k]; !ok {
			if add == nil {
				if res != "" {
					res += "\n"
				}
				res += "添加的主机为:"
			}
			add = append(add, k)
			res += v + ","
			ctl.Debug(v)
			ctl.Debug(res)
			srv = append(srv, al.AppNginxServer{AppNginxId: nid, HostId: k})
		}
	}
	ctl.Debug(res)

	if add == nil && del == nil {
		return res, ctl.Errorf("nginx 成员主机未作任何修改")
	}
	ln1, err := db.Db.In("host_id", del).Where("app_nginx_id=?", nid).Delete(&al.AppNginxServer{AppNginxId: nid})
	if err != nil {
		return res, err
	}
	if int(ln1) != len(del) {
		return res, ctl.Errorf("错误:删除的行与需删除的行不一致")
	}
	ln2, err := db.Db.Insert(&srv)
	if err != nil {
		return res, err
	}
	if int(ln2) != len(add) {
		return res, ctl.Errorf("错误:添加的行与需添加的行不一致")
	}
	return res, err
}
func SetNginxAdd(c *gin.Context) {
	ctl.Debug(c)
	var (
		info      NginxSetData
		res, res1 string
		err       error
		fun       al.AppFun
		ngx       al.AppNginx
	)
	defer func() {
		if err != nil {
			res += err.Error()
			ctl.Log.Error(err)
		} else {
			ctl.Log.Debug(err)
		}
		ctl.Log.Debug(res)
		c.Data(200, "", []byte(res))
	}()

	info.AppNginxName = c.PostForm("AppNginxName")
	info.FunName = c.PostForm("FunName")
	info.NginxFile = c.PostForm("NginxFile")
	info.NginxFirstFile = c.PostForm("NginxFirstFile")
	info.NginxFirstScript = c.PostForm("NginxFirstScript")
	info.NginxSecondFile = c.PostForm("NginxSecondFile")
	info.NginxSecondScript = c.PostForm("NginxSecondScript")
	info.Member = c.PostForm("Member")
	if info.AppNginxName == "" {
		err = ctl.Errorf("nginx名为空")
		return
	} else {
		db.Db.Where("app_nginx_name=?", info.AppNginxName).Get(&ngx)
		ctl.Log.Debug(ngx)
		ctl.Log.Debug(info.AppNginxName)
		ctl.Debug(ngx)
		ctl.Debug(info.AppNginxName)
		if ngx.AppNginxId != 0 {
			err = ctl.Errorf("nginx名己存在")
			return
		}
	}
	if info.FunName == "" {
		err = ctl.Errorf("应用功能名为空")
		return
	} else {
		db.Db.Where("fun_name=?", info.FunName).Get(&fun)
		ctl.Debug(fun)
		ctl.Debug(info.FunName)
		if fun.FunId == 0 {
			err = ctl.Errorf("应用功能名不存在")
			return
		}
	}
	if info.NginxFile == "" {
		err = ctl.Errorf("nginx文件为空")
		return
	}
	if info.NginxFirstFile == "" {
		err = ctl.Errorf("nginx一服文件为空")
		return
	}
	if info.NginxFirstScript == "" {
		err = ctl.Errorf("nginx一服脚本为空")
		return
	}
	if info.NginxSecondFile == "" {
		err = ctl.Errorf("nginx二服文件为空")
		return
	}
	if info.NginxSecondScript == "" {
		err = ctl.Errorf("nginx二服脚本为空")
		return
	}

	ngx.AppNginxName = info.AppNginxName
	ngx.FunId = fun.FunId
	ngx.NginxFile = info.NginxFile
	ngx.NginxFirstFile = info.NginxFirstFile
	ngx.NginxFirstScript = info.NginxFirstScript
	ngx.NginxSecondFile = info.NginxSecondFile
	ngx.NginxSecondScript = info.NginxSecondScript
	ln, err := db.Db.InsertOne(&ngx)
	if err != nil {
		return
	}
	if ln != 1 {
		err = ctl.Errorf("错误:插入的行不为1,插入的行为:%d", ln)
		return
	}
	res = "插入" + ngx.AppNginxName + "成功"
	db.Db.Where("app_nginx_name=? and fun_id=?", info.AppNginxName, fun.FunId).Get(&ngx)
	res1, err = SetNgxMember(info.Member, ngx.AppNginxId)
	res += "\n" + res1
}

func SetNginxEdit(c *gin.Context) {
	ctl.Debug(c)
	var (
		info NginxSetData
		res  string
		ngx  al.AppNginx
		diff string
		fun  al.AppFun
		res1 string
		err  error
	)
	defer func() {
		if err != nil {
			ctl.Debug(info)
			res += err.Error()
			ctl.Log.Error(err)
		} else {
			ctl.Log.Debug(err)
		}
		c.Data(200, "", []byte(res))
	}()

	info.AppNginxId = c.PostForm("AppNginxId")
	info.AppNginxName = c.PostForm("AppNginxName")
	info.FunName = c.PostForm("FunName")
	info.NginxFile = c.PostForm("NginxFile")
	info.NginxFirstFile = c.PostForm("NginxFirstFile")
	info.NginxFirstScript = c.PostForm("NginxFirstScript")
	info.NginxSecondFile = c.PostForm("NginxSecondFile")
	info.NginxSecondScript = c.PostForm("NginxSecondScript")
	ngx.AppNginxId = ctl.AtoiNe(info.AppNginxId)
	info.Member = c.PostForm("Member")
	if ngx.AppNginxId == 0 {
		err = ctl.Errorf("nginx ID为空")
		return
	} else {
		ngx = al.AppNginx{}
		ctl.Debug(info.AppNginxId)
		db.Db.Where("app_nginx_id=?", info.AppNginxId).Get(&ngx)
		if ngx.AppNginxId == 0 {
			err = ctl.Errorf("应用nginx不存在")
			return
		}
	}
	if info.AppNginxName == "" {
		err = ctl.Errorf("应用nginx名为空")
		return
	} else {
		db.Db.Where("app_nginx_name=?", info.AppNginxName).Get(&ngx)
		if ngx.AppNginxId == 0 {
			err = ctl.Errorf("应用nginx不存在")
			return
		}
		if info.AppNginxName != info.AppNginxName {
			diff += "应用nginx名,"
		}
	}

	if info.FunName == "" {
		err = ctl.Errorf("应用功能名为空")
		return
	} else {
		db.Db.Where("fun_name=? and app_zone=?", info.FunName, auth.GetState(c).State).Get(&fun)
		if fun.FunId == 0 {
			err = ctl.Errorf("应用功能名不存在")
			return
		}
		if fun.FunId != ngx.FunId {
			diff += "应用功能名,"
		}
	}

	if info.NginxFile == "" {
		err = ctl.Errorf("nginx文件为空")
		return
	} else if info.NginxFile != ngx.NginxFile {
		diff += "nginx文件,"
	}
	if info.NginxFirstFile == "" {
		err = ctl.Errorf("nginx一服文件为空")
		return
	} else if info.NginxFirstFile != ngx.NginxFirstFile {
		diff += "nginx一服文件,"
	}
	if info.NginxFirstScript == "" {
		err = ctl.Errorf("nginx一服脚本为空")
		return
	} else if info.NginxFirstScript != ngx.NginxFirstScript {
		diff += "nginx一服脚本,"
	}
	if info.NginxSecondFile == "" {
		err = ctl.Errorf("nginx二服文件为空")
		return
	} else if info.NginxSecondFile != ngx.NginxSecondFile {
		diff += "nginx二服文件,"
	}
	if info.NginxSecondScript == "" {
		err = ctl.Errorf("nginx二服脚本为空")
		return
	} else if info.NginxSecondScript != ngx.NginxSecondScript {
		diff += "nginx二服脚本,"
	}
	if diff == "" {
		res = "nginx配置项未作任何修改"
		res1, err = SetNgxMember(info.Member, ngx.AppNginxId)
		res += "\n" + res1
		return
	}
	ngx.AppNginxName = info.AppNginxName
	ngx.FunId = fun.FunId
	ngx.NginxFile = info.NginxFile
	ngx.NginxFirstFile = info.NginxFirstFile
	ngx.NginxFirstScript = info.NginxFirstScript
	ngx.NginxSecondFile = info.NginxSecondFile
	ngx.NginxSecondScript = info.NginxSecondScript

	ln, err := db.Db.Where("app_nginx_id=?", ngx.AppNginxId).Update(&ngx)
	if err != nil {
		return
	}
	if ln != 1 {
		err = ctl.Errorf("错误:修改的行不为1,插入的行为:%d", ln)
		return
	}
	res = "修改" + info.AppNginxName + "成功\n修改项:" + diff[:len(diff)-1]
	res1, err = SetNgxMember(info.Member, ngx.AppNginxId)
	res += "\n" + res1
	go func() { op.SyncInfo() }()
}

func SetNginxDel(c *gin.Context) {
	ctl.Debug(c)
	var (
		info NginxSetData
		res  string
		ngx  al.AppNginx
		err  error
	)
	defer func() {
		if err != nil {
			ctl.Debug(info)
			res += err.Error()
			ctl.Log.Error(err)
		} else {
			ctl.Log.Debug(err)
		}
		c.Data(200, "", []byte(res))
	}()
	err = c.BindJSON(&info)
	ctl.Debug(info)
	ctl.Debug(err)
	if err != nil {
		return
	}
	ngx.AppNginxId = ctl.AtoiNe(info.AppNginxId)
	if ngx.AppNginxId == 0 {
		err = ctl.Errorf("nginx ID为空")
		return
	}
	_, err = db.Db.Where("app_nginx_id=?", ngx.AppNginxId).Delete(&ngx)
	if err != nil {
		return
	}
	_, err = db.Db.Where("app_nginx_id=?", ngx.AppNginxId).Delete(&al.AppNginxServer{})
	if err != nil {
		return
	}
	res = "删除" + info.AppNginxName + "成功"
}
