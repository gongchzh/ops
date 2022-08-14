package ah

import (
	"ops/pkg/ctl"
	"ops/pkg/al"
	"ops/pkg/auth"
	"ops/pkg/db"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AuditD struct {
	UpdateQueueId int
	FunName       string
	UpdateProgram string
	Md5           string
	UpdateTime    string
	Status        int
	FunId         int
	AppId         int
	UpEarlTime    string
	UpLateTime    string
}
type AuditD1 struct {
	UpdateQueueId string
	FunName       string
	UpdateProgram string
	Md5           string
	UpdateTime    string
	Status        string
	FunId         string
	AppId         string
	UpEarlTime    string
	UpLateTime    string
}

type AuditFileData struct {
	Md5Old        string
	Md5New        string
	AppId         int
	UpdateQueueId int
	UpdateTime    string
}

func AuditData(c *gin.Context) {
	var (
		aud []AuditD

		/*md5o    = make(map[int]string)
		md5n    = make(map[int]string)
		modTime = make(map[int]string)*/

	)
	defer func() {
		c.JSON(200, &aud)
	}()
	db.Db.SQL(`select a.update_queue_id,b.fun_id,b.fun_name,c.update_program,c.app_id,a.md5,a.update_time,a.status,a.up_earl_time,a.up_late_time from app_update_queue a,
app_fun b,app c
where a.fun_id=b.fun_id
and b.fun_id=c.fun_id
and a.status in (2,3)
and b.app_zone=?
group by a.update_queue_id
`, auth.GetState(c).State).Find(&aud)

	aud = GetAudInfo(aud, c)
	ctl.Debug(aud)
}

func GetAudInfo(aud []AuditD, c *gin.Context) []AuditD {
	var (
		ids   []int
		apps  al.Apps
		filem = make(map[int]map[int]*AuditFileData)
		list  []al.AppList
	)
	for _, v := range aud {
		if filem[v.AppId] == nil {
			filem[v.AppId] = make(map[int]*AuditFileData)
		}
		filem[v.AppId][v.UpdateQueueId] = &AuditFileData{}
		filem[v.AppId][v.UpdateQueueId].Md5Old = v.Md5
		filem[v.AppId][v.UpdateQueueId].UpdateQueueId = v.UpdateQueueId

		ids = append(ids, v.AppId)
	}

	db.Db.In("app_id", ids).Find(&apps)
	ctl.Debug(apps)
	if apps == nil {
		aud = nil
		return aud
	}
	list = apps.GetInfo()
	ctl.Debug(list)
	if list == nil {
		aud = nil
		return aud
	}
	for _, v := range list {

		if _, ok := filem[v.AppId]; !ok {
			continue
		}
		for k, _ := range filem[v.AppId] {
			filem[v.AppId][k].Md5New = v.NewMd5
			filem[v.AppId][k].UpdateTime = v.NewTime
		}
		ctl.Debug(v.AppName, v.RunState, v.NewMd5, v.NewTime)

	}
	for k, v := range aud {
		ctl.Debug(v.FunName, v.Status, v.UpdateTime)
		ctl.Debug(v.FunName, filem[v.AppId][v.UpdateQueueId].Md5New)
		ctl.Debug(v.FunName, filem[v.AppId][v.UpdateQueueId].Md5Old)
		if filem[v.AppId][v.UpdateQueueId].Md5New != filem[v.AppId][v.UpdateQueueId].Md5Old {
			aud[k].Status = 10
			aud[k].UpdateTime = ""
			continue
		}
		ctl.Debug(v.FunName, v.Status, v.UpdateTime)
		aud[k].UpdateTime = filem[v.AppId][v.UpdateQueueId].UpdateTime
	}
	return aud
}

func AuditPass(c *gin.Context) {
	var (
		res  string
		err  error
		aud  []AuditD1
		audn []AuditD
		audt AuditD
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
	err = c.BindJSON(&aud)
	if err != nil {
		return
	}
	if len(aud) == 0 {
		err = ctl.Errorf("未选择队列")
		return
	}
	ctl.Debug(aud)
	for _, v := range aud {
		if v.Status != "待审核" {
			err = ctl.Errorf(v.FunName + ":状态为非待审核")
			return
		}
		audt = AuditD{}
		db.Db.SQL(`select a.update_queue_id,b.fun_id,b.fun_name,c.update_program,c.app_id,a.md5,a.update_time,a.status,a.up_earl_time,a.up_late_time from app_update_queue a,
app_fun b,app c
where a.fun_id=b.fun_id
and b.fun_id=c.fun_id
and a.status=3
and a.update_queue_id=?
and a.md5=?
and b.fun_name=?
and c.update_program=?
and a.up_earl_time=?
and a.up_late_time=?
group by a.update_queue_id
`, v.UpdateQueueId, v.Md5, v.FunName, v.UpdateProgram, v.UpEarlTime, v.UpLateTime).Get(&audt)
		ctl.Debug(v.UpdateQueueId, v.Md5, v.FunName, v.UpdateProgram, v.UpEarlTime, v.UpLateTime)
		if audt.UpdateQueueId == 0 {
			err = ctl.Errorf(v.FunName + ":获取数据库队列失败")
			return
		}
		audn = append(audn, audt)
	}
	audn = GetAudInfo(audn, c)
	for _, v := range audn {
		if v.Status != 3 {
			err = ctl.Errorf(v.FunName + ":获取应用服文件状态异常")
		}
	}
	if len(aud) != len(audn) {
		err = ctl.Errorf(":校验通过队列数与原队列数不一致")
		return
	}
	for _, v := range audn {
		ln, err := db.Db.Where("update_queue_id=? and status=3", v.UpdateQueueId).Cols("status").Update(&al.AppUpdateQueue{Status: 2})
		if err != nil {
			err = ctl.Errorf(v.FunName + ":" + err.Error())
			return
		}
		if ln != 1 {
			err = ctl.Errorf(v.FunName+":修改的行数不为1,行数为:%d", ln)
			return
		}
		res += "审核应用功能名 " + v.FunName + " ,队列ID " + strconv.Itoa(v.UpdateQueueId) + "通过\n"
	}

}
func AuditDelQueue(c *gin.Context) {
	var (
		res string
		err error
		aud []AuditD1
		que al.AppUpdateQueue
	)
	defer func() {
		if err != nil {
			res = err.Error()
			ctl.Log.Error(err)
		} else {
			ctl.Log.Debug(err)
		}

		if err != nil {
			ctl.Log.Error(err)
		} else {
			ctl.Log.Debug("发送切服钉钉消息成功")
		}
		c.Data(200, "", []byte(res))

	}()
	err = c.BindJSON(&aud)
	if err != nil {
		return
	}
	ctl.Debug(len(aud))
	if len(aud) != 1 {
		err = ctl.Errorf("一次只能删除一条")
		return
	}

	db.Db.Where("update_queue_id=? and md5=? and up_earl_time=? and up_late_time=? and status=3", aud[0].UpdateQueueId, aud[0].Md5, aud[0].UpEarlTime, aud[0].UpLateTime).Get(&que)
	if que.UpdateQueueId == 0 {
		err = ctl.Errorf("获取数据库队列异常")
		return
	}
	ln, err := db.Db.Where("update_queue_id=? and md5=? and up_earl_time=? and up_late_time=? and status=3", aud[0].UpdateQueueId, aud[0].Md5, aud[0].UpEarlTime, aud[0].UpLateTime).Delete(&que)
	if err != nil {
		return
	}
	if ln != 1 {
		err = ctl.Errorf("删除队列的行数不对,行数为 %d", ln)
	}
	res = "删除应用功能名 " + aud[0].FunName + ",队列ID为" + strconv.Itoa(que.UpdateQueueId) + " 成功"
}
