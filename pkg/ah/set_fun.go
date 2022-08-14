package ah

import (
	"ops/pkg/ctl"
	"ops/pkg/al"
	"ops/pkg/auth"
	"ops/pkg/db"

	"github.com/gin-gonic/gin"
)

func SetFunAdd(c *gin.Context) {
	var (
		res string
		err error
		fun al.AppFun
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
	fun.FunName = c.PostForm("FunName")
	if fun.FunName == "" {
		err = ctl.Errorf("功能名不能为空")
		return
	}
	fun.AppZone = auth.GetState(c).State
	ctl.Debug(fun)
	ln, err := db.Db.InsertOne(&fun)
	if err != nil {
		return
	}
	if ln != 1 {
		err = ctl.Errorf("插入的行数不对,插入的行 %d", ln)
		return
	}
	res = "插入" + fun.FunName + "成功"
}
