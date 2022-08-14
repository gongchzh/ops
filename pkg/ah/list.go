package ah

import (
	"ops/pkg/ctl"
	"ops/pkg/al"
	"ops/pkg/auth"
	"ops/pkg/db"

	"github.com/gin-gonic/gin"
)

func ListData(c *gin.Context) {
	var (
		apps al.Apps
	)
	ctl.Debug(auth.GetState(c))
	ctl.Debug(auth.GetState(c).State)
	db.Db.SQL(`select a.* from app a,app_fun b
	where a.fun_id=b.fun_id
	and b.app_zone=? order by a.app_name`, auth.GetState(c).State).Find(&apps)
	db.Db.SQL(`select a.* from app a left join app_fun b on a.fun_id=b.fun_id
where (b.app_zone=? or a.fun_id=0)`, auth.GetState(c).State).Find(&apps)
	//	c.JSON(200, GetAppInfo(app, auth.GetState(c).State))
	//	c.JSON(200, GetAppInfo(app, auth.GetState(c).State))
	ctl.Debug(apps)
	ctl.Debug(apps.GetInfo())
	c.JSON(200, apps.GetInfo())
}
