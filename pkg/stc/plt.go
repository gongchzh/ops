package stc

import (
	"ops/pkg/ctl"
	"ops/pkg/auth"
	"ops/pkg/cfg"
	"ops/pkg/db"

	"github.com/gin-gonic/gin"
)

func PltList(c *gin.Context) {
	c.HTML(200, "plt_list.html", nil)
}
func PltAudit(c *gin.Context) {
	c.HTML(200, "plt_audit.html", nil)
}
func PltHeader2(c *gin.Context) {
	ctl.Debug("plt.html")
	ctl.Debug(cfg.GameHeaderHtml.NormState)
	c.HTML(200, "plt.html", &cfg.GameHeaderHtml)
}
func Test1JiekouLeft(c *gin.Context) {
	ctl.Debug("test1.html")
	c.HTML(200, "test1_jiekou_left.html", nil)
}
func Test1JiekouList(c *gin.Context) {
	ctl.Debug("test1.html")
	c.HTML(200, "test1_jiekou_list.html", nil)
}

func PltLeft(c *gin.Context) {
	var (
		user db.Auth
	)

	if c.GetString("user") == "" {
		return
	}

	db.Db.Where("user=?", c.GetString("user")).Get(&user)
	if user.AuthId == 0 {
		return
	}
	c.HTML(200, "plt_left.html", auth.GetAces(user.AuthStatus))
}

/*func PltOpsTime(c *gin.Context) {
	c.HTML(200, "plt_ops_time.html", nil)
}*/

func PltSetFun(c *gin.Context) {
	c.HTML(200, "plt_set_fun.html", nil)
}
func PltSetGroup(c *gin.Context) {
	c.HTML(200, "plt_set_group.html", nil)
}
func PltQueue(c *gin.Context) {
	c.HTML(200, "plt_queue.html", nil)
}
func PltSetApp(c *gin.Context) {
	c.HTML(200, "plt_set_app.html", nil)
}
func PltSetServer(c *gin.Context) {
	c.HTML(200, "plt_set_server.html", nil)
}

func PltSetNginx(c *gin.Context) {
	c.HTML(200, "plt_set_nginx.html", nil)
}

func PltSetLeft(c *gin.Context) {
	c.HTML(200, "plt_set_left.html", nil)
}
func PltUpdateLeft(c *gin.Context) {
	c.HTML(200, "plt_update_left.html", nil)
}
func PltAuditLeft(c *gin.Context) {
	c.HTML(200, "plt_audit_left.html", nil)
}
func PltListLeft(c *gin.Context) {
	c.HTML(200, "plt_list_left.html", nil)
}
