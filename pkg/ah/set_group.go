package ah

import (
	"ops/pkg/ctl"
	"ops/pkg/al"
	"ops/pkg/auth"
	"ops/pkg/db"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type GroupData struct {
	AppGroupId   int
	AppGroupName string
	MemberFun    string
	MemberGroup  string
	OpsEarlTime  string
	OpsLateTime  string
	Status       int
	State        int
	NewEarlTime  string
	NewLateTime  string
}

type GroupMemberFun struct {
	AppGroupId   int
	AppGroupName string
	FunName      string
	FunId        int
}

type GroupMemberSub struct {
	AppGroupId   int
	AppGroupName string
	SubGroupName string
	SubGroupId   int
}

func SetGroupData(c *gin.Context) {
	var (
		grp   []al.AppGroup
		info  []GroupData
		infot GroupData
		memf  []GroupMemberFun
		memfm = make(map[int]string)
		mems  []GroupMemberSub
		memsm = make(map[int]string)
	)
	ctl.Debug("abc")
	db.Db.Where("state=? or state=100", auth.GetState(c).State).Find(&grp)
	db.Db.SQL(`select a.app_group_id,a.app_group_name,c.fun_name from app_group a,app_group_member b,app_fun c
where a.app_group_id=b.app_group_id
and b.fun_id=c.fun_id
and (a.state=? or a.state=100)`, auth.GetState(c).State).Find(&memf)

	for _, v := range memf {
		memfm[v.AppGroupId] += v.FunName + " "
	}
	db.Db.SQL(`select a.app_group_id,a.app_group_name,c.app_group_name as sub_group_name from app_group a,app_group_sub_member b,app_group c
where a.app_group_id=b.app_group_id
and b.sub_group_id=c.app_group_id
and (a.state=? or a.state=100)`, auth.GetState(c).State).Find(&mems)
	for _, v := range mems {
		memsm[v.AppGroupId] += v.SubGroupName + " "
	}
	ctl.Debug(grp)
	for k, v := range grp {
		infot = GroupData{}
		infot.AppGroupId = v.AppGroupId
		infot.AppGroupName = v.AppGroupName
		infot.OpsEarlTime = v.OpsEarlTime
		infot.OpsLateTime = v.OpsLateTime
		infot.State = v.State
		te, err := time.ParseInLocation(ctl.TimeFormat, v.OpsEarlTime, time.Local)
		ctl.Debug(err)
		if err != nil {
			infot.NewEarlTime = time.Now().Format("2006-01-02 15:04:05")
		} else {
			if v.OpsEarlTime == "" {
				infot.NewEarlTime = time.Now().Format("2006-01-02 15:04:05")
			} else {
				ctl.Debug(time.Now())
				ctl.Debug(te)
				ctl.Debug(time.Now().Sub(te))
				if time.Now().Sub(te) > time.Hour*5 {
					infot.NewEarlTime = time.Now().Format("2006-01-02 15:04:05")
				}
			}
		}
		tt, err := time.ParseInLocation(ctl.TimeFormat, v.OpsLateTime, time.Local)
		if err != nil {
			infot.NewLateTime = time.Now().Add(time.Hour * 6).Format("2006-01-02 15:04:05")
		} else {
			if v.OpsEarlTime == "" {
				infot.NewLateTime = time.Now().Add(time.Hour * 6).Format("2006-01-02 15:04:05")
			} else {
				if tt.Sub(time.Now()) < time.Hour {
					infot.NewLateTime = time.Now().Add(time.Hour * 6).Format("2006-01-02 15:04:05")
				}
			}
		}
		ctl.Debug(grp[k].OpsEarlTime)
		ctl.Debug(grp[k].OpsLateTime)
		if time.Now().Sub(te) > time.Second/10 && tt.Sub(time.Now()) > time.Minute {
			infot.Status = 1
		}
		if infot.NewEarlTime == "" {
			infot.NewEarlTime = v.OpsEarlTime
		}
		if infot.NewLateTime == "" {
			infot.NewLateTime = v.OpsLateTime
		}
		if _, ok := memfm[v.AppGroupId]; ok {
			infot.MemberFun = memfm[v.AppGroupId]
		}
		if _, ok := memsm[v.AppGroupId]; ok {
			infot.MemberGroup = memsm[v.AppGroupId]
		}
		info = append(info, infot)
	}
	c.JSON(200, &info)

}

func SetGroupAdd(c *gin.Context) {
	var (
		grp al.AppGroup
		res string
		err error
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
	switch c.PostForm("State") {
	case "区域组":
		grp.State = auth.GetState(c).State
	case "跨区域组":
		grp.State = 100
	default:
		err = ctl.Errorf("类型异常")
		return
	}
	grp.AppGroupName = c.PostForm("AppGroupName")

	if grp.AppGroupName == "" {
		err = ctl.Errorf("应用组名为空")
		return
	}
	ln, err := db.Db.InsertOne(&grp)
	if err != nil {
		return
	}
	if ln != 1 {
		err = ctl.Errorf("插入的行数不对,插入的行 %d", ln)
		return
	}
	res = "插入" + grp.AppGroupName + "成功"
}

func SetGroupEdit(c *gin.Context) {
	ctl.Debug(c)
	var (
		info al.AppGroup
		res  string
		err  error
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
	info.AppGroupId, _ = strconv.Atoi(c.PostForm("AppGroupId"))
	info.AppGroupName = c.PostForm("AppGroupName")
	info.OpsEarlTime = c.PostForm("OpsEarlTime")
	info.OpsLateTime = c.PostForm("OpsLateTime")
	if info.AppGroupId == 0 || info.AppGroupName == "" || info.OpsEarlTime == "" || info.OpsLateTime == "" {
		err = ctl.Errorf("参数有空值")
		return
	}
	ln, err := db.Db.Where("app_group_id=? and app_group_name=?", info.AppGroupId, info.AppGroupName).Cols("ops_earl_time,ops_late_time").Update(&al.AppGroup{AppGroupId: info.AppGroupId, OpsEarlTime: info.OpsEarlTime, OpsLateTime: info.OpsLateTime})
	if err != nil {
		return
	}
	if ln != 1 {
		err = ctl.Errorf("修改维护时间的行数不对,受影响的行 %d", ln)
		return
	}
	res = "修改" + info.AppGroupName + "维护时间成功"
}

func SetGroupEditFun(c *gin.Context) {
	ctl.Debug(c)
	var (
		info     GroupData
		res      string
		err      error
		memf     []GroupMemberFun
		memfm    = make(map[string]GroupMemberFun)
		add      []string
		deln     []int
		fun      []al.AppFun
		memi     []al.AppGroupMember
		grpo     al.AppGroup
		memn     = make(map[string]int)
		ln1, ln2 int64
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
	info.AppGroupId, _ = strconv.Atoi(c.PostForm("AppGroupId"))
	info.AppGroupName = c.PostForm("AppGroupName")
	info.MemberFun = c.PostForm("MemberFun")
	if info.AppGroupId == 0 || info.AppGroupName == "" {
		err = ctl.Errorf("参数有空值")
		return
	}
	ctl.Debug(info.MemberFun)
	db.Db.Where("app_group_id=? and app_group_name=?", info.AppGroupId, info.AppGroupName).Get(&grpo)
	if grpo.AppGroupId == 0 {
		err = ctl.Errorf("获取数据库应用组异常")
		return
	}
	db.Db.SQL(`select a.app_group_id,a.app_group_name,c.fun_name,c.fun_id from app_group a,app_group_member b,app_fun c
where a.app_group_id=b.app_group_id
and b.fun_id=c.fun_id
and a.app_group_id=? and a.app_group_name=?`, info.AppGroupId, info.AppGroupName).Find(&memf)
	for _, v := range memf {
		memfm[v.FunName] = v
	}
	for _, v := range strings.Split(info.MemberFun, " ") {
		if v == "" {
			continue
		}
		memn[v] = 0
		if _, ok := memfm[v]; ok {
			continue
		}
		add = append(add, v)
	}
	for k, v := range memfm {
		if _, ok := memn[k]; ok {
			continue
		}
		deln = append(deln, v.FunId)
	}
	if deln != nil {
		ln1, err = db.Db.In("fun_id", deln).Where("app_group_id=?", info.AppGroupId).Delete(&al.AppGroupMember{})
		if err != nil {
			ctl.Debug(err)
			return
		}
	}
	if add != nil {
		db.Db.In("fun_name", add).Find(&fun)
		for _, v := range fun {
			memi = append(memi, al.AppGroupMember{FunId: v.FunId, AppGroupId: info.AppGroupId})
		}
		ln2, err = db.Db.Insert(&memi)
		if err != nil {
			ctl.Debug(err)
			return
		}
	}
	if deln == nil && add == nil {
		err = ctl.Errorf("未作任何修改")
		return
	}
	ctl.Debug(add)
	ctl.Debug(deln)
	if int(ln1) != len(deln) {
		err = ctl.Errorf("错误:删除的行不对,需删除的行为:%d,实际删除的行为:%d", len(deln), ln1)
	}
	if int(ln2) != len(add) {
		err = ctl.ErrorAf("%s\n错误:添加的行不对,需添加的行为:%d,实际添加的行为:%d", err, len(add), ln2)
	}
	if err == nil {
		res = "修改" + info.AppGroupName + "成员应用功能成功"
	}

}
func SetGroupEditGroup(c *gin.Context) {
	ctl.Debug(c)
	var (
		info     GroupData
		res      string
		err      error
		mems     []GroupMemberSub
		memsm    = make(map[string]GroupMemberSub)
		add      []string
		deln     []int
		grp      []al.AppGroup
		memi     []al.AppGroupSubMember
		memn     = make(map[string]int)
		ln1, ln2 int64
		grpo     al.AppGroup
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
	info.AppGroupId, _ = strconv.Atoi(c.PostForm("AppGroupId"))
	info.AppGroupName = c.PostForm("AppGroupName")
	info.MemberGroup = c.PostForm("MemberGroup")
	if info.AppGroupId == 0 || info.AppGroupName == "" {
		err = ctl.Errorf("参数有空值")
		return
	}
	ctl.Debug(info.MemberGroup)
	db.Db.Where("app_group_id=? and app_group_name=?", info.AppGroupId, info.AppGroupName).Get(&grpo)
	if grpo.AppGroupId == 0 {
		err = ctl.Errorf("获取数据库应用组异常")
		return
	}
	db.Db.SQL(`select a.app_group_id,a.app_group_name,c.app_group_name as sub_group_name,b.sub_group_id from app_group a,app_group_sub_member b,app_group c
where a.app_group_id=b.app_group_id
and b.sub_group_id=c.app_group_id
and a.app_group_id=? and a.app_group_name=?`, info.AppGroupId, info.AppGroupName).Find(&mems)
	for _, v := range mems {
		memsm[v.SubGroupName] = v
	}
	for _, v := range strings.Split(info.MemberGroup, " ") {
		if v == "" {
			continue
		}
		if v == info.AppGroupName {
			err = ctl.Errorf("成员不能包含自身")
			return
		}
		memn[v] = 0
		if _, ok := memsm[v]; ok {
			continue
		}
		add = append(add, v)
	}
	for k, v := range memsm {
		if _, ok := memn[k]; ok {
			continue
		}
		deln = append(deln, v.SubGroupId)
	}
	if add != nil {
		db.Db.In("app_group_name", add).Find(&grp)
		for _, v := range grp {
			if grpo.State != 100 && grpo.State != v.State {
				err = ctl.Errorf("普通组不能包含其它区组")
				return
			}
			memi = append(memi, al.AppGroupSubMember{SubGroupId: v.AppGroupId, AppGroupId: info.AppGroupId})
		}
		ln2, err = db.Db.Insert(&memi)
		if err != nil {
			ctl.Debug(err)
			return
		}
	}
	if deln != nil {
		ln1, err = db.Db.In("sub_group_id", deln).Where("app_group_id=?", info.AppGroupId).Delete(&al.AppGroupSubMember{})
		if err != nil {
			ctl.Debug(err)
			return
		}
	}

	if deln == nil && add == nil {
		err = ctl.Errorf("未作任何修改")
		return
	}
	ctl.Debug(add)
	ctl.Debug(deln)
	if int(ln1) != len(deln) {
		err = ctl.Errorf("错误:删除的行不对,需删除的行为:%d,实际删除的行为:%d", len(deln), ln1)
	}
	if int(ln2) != len(add) {
		err = ctl.ErrorAf("%s\n错误:添加的行不对,需添加的行为:%d,实际添加的行为:%d", err, len(add), ln2)
	}
	if err == nil {
		res = "修改" + info.AppGroupName + "成员应用组成功"
	}
}
