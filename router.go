package main

import (
	"ops/pkg/ctl"
	"ops/pkg/ah"
	"ops/pkg/auth"
	"ops/pkg/stc"


	"github.com/gin-gonic/gin"
)

func router() {
	var (
		r   *gin.Engine
		r1  *gin.RouterGroup
		err error
	)
	gin.SetMode(gin.ReleaseMode)
	r = gin.New()
	r.Use(gin.Recovery())

	r.Use(gin.Logger())
	r.LoadHTMLGlob("html/*.html")
	r.Static("/js", "js")
	r.Static("/font", "font")
	r.Static("/images", "images")
	r.Static("/css", "css")
	r.GET("/login", stc.Login)
	r.POST("/auth", auth.GetAuth)
	r1 = r.Group("/v1")
	r1.Use(auth.JWT())
	{
		r1.Static("/js", "js")
		r1.Static("/font", "font")
		r1.Static("/images", "images")
		r1.Static("/css", "css")
		r1.GET("/index", stc.Index)
		r1.GET("/header", stc.Header)
		r1.GET("/default", stc.DefaultHtml)
		r1.GET("/plt/:plt_type", stc.PltHeader)
		r1.GET("/plt", stc.PltHeader)
		r1.GET("/display", stc.DisplayHtml)

		r1.GET("/pop/:name/:action", stc.Pop)
		r1.GET("/pop/:name/:action/*sub", stc.Pop)

		r1.GET("/plt/:plt_type/list/list", stc.PltList)
		r1.POST("/plt/:plt_type/list/data", ah.ListData)
		r1.POST("/plt/:plt_type/list/switch_frt", ah.Switch)
		r1.POST("/plt/:plt_type/list/switch_sed", ah.Switch)
		r1.POST("/plt/:plt_type/list/restart", ah.Restart)
		r1.GET("/plt/:plt_type/update/add", ah.AddQueue)
		r1.POST("/plt/:plt_type/update/form", ah.AddQueueForm)

		r1.GET("/plt/:plt_type/audit/list", stc.PltAudit)
		r1.POST("/plt/:plt_type/audit/data", ah.AuditData)
		r1.POST("/plt/:plt_type/audit/pass", ah.AuditPass)
		r1.POST("/plt/:plt_type/audit/del_queue", ah.AuditDelQueue)

		r1.GET("/plt/:plt_type/set/fun", stc.PltSetFun)
		r1.POST("/plt/:plt_type/set/fun/data", ah.SetFunData)
		r1.GET("/plt/:plt_type/set/group", stc.PltSetGroup)
		r1.POST("/plt/:plt_type/set/group/data", ah.SetGroupData)

		r1.GET("/plt/:plt_type/update/left", stc.PltUpdateLeft)
		r1.GET("/plt/:plt_type/list/left", stc.PltListLeft)
		r1.GET("/plt/:plt_type/set/left", stc.PltSetLeft)
		r1.GET("/plt/:plt_type/audit/left", stc.PltAuditLeft)

		r1.GET("/plt/:plt_type/update/list", stc.PltQueue)
		r1.POST("/plt/:plt_type/update/data", ah.QueueData)
		r1.POST("/plt/:plt_type/update/update_single", ah.Update)
		r1.POST("/plt/:plt_type/update/update_multi", ah.Update)
		r1.GET("/plt/:plt_type/set/app", stc.PltSetApp)
		r1.GET("/plt/:plt_type/set/server", stc.PltSetServer)
		r1.GET("/plt/:plt_type/set/nginx", stc.PltSetNginx)
		r1.POST("/plt/:plt_type/set/app/data", ah.SetAppData)
		r1.POST("/plt/:plt_type/set/server/data", ah.SetServerData)
		r1.POST("/plt/:plt_type/set/nginx/data", ah.SetNginxData)

		r1.POST("/plt/:plt_type/set/nginx/del", ah.SetNginxDel)

		r1.POST("/pop/plt/:game_type/set/fun/edit", ah.SetFunEdit)
		r1.POST("/pop/plt/:game_type/set/group/edit", ah.SetGroupEdit)

		r1.POST("/pop/plt/:game_type/set/app/add", ah.SetAppAdd)
		r1.POST("/pop3/plt/:game_type/set/app/add", ah.SetAppAdd)
		r1.POST("/pop/plt/:game_type/set/server/add", ah.SetServerAdd)
		r1.POST("/pop3/plt/:game_type/set/server/add", ah.SetServerAdd)
		r1.POST("/pop/plt/:game_type/set/nginx/add", ah.SetNginxAdd)
		r1.POST("/pop3/plt/:game_type/set/nginx/add", ah.SetNginxAdd)

		r1.GET("/pop2/pop/:name/:action", stc.Pop2)
		r1.GET("/pop2/pop/:name/:action/*sub", stc.Pop2)

		r1.GET("/pop2/game/:name/:action", stc.Pop2)
		r1.GET("/pop2/game/:name/:action/*sub", stc.Pop2)

		r1.GET("/pop2/plt/:name/:action", stc.Pop2)
		r1.GET("/pop2/plt/:name/:action/*sub", stc.Pop2)

		r1.POST("/pop2/plt/:game_type/set/app/edit", ah.SetAppEdit)
		r1.POST("/pop2/plt/:game_type/set/server/edit", ah.SetServerEdit)
		r1.POST("/pop2/plt/:game_type/set/nginx/edit", ah.SetNginxEdit)
		r1.POST("/pop2/plt/:game_type/set/fun/add", ah.SetFunAdd)

		r1.POST("/pop2/plt/:game_type/set/group/add", ah.SetGroupAdd)

		r1.POST("/pop3/plt/:game_type/set/group/edit_fun", ah.SetGroupEditFun)
		r1.POST("/pop4/plt/:game_type/set/group/edit_group", ah.SetGroupEditGroup)

		r1.GET("/pop3/plt/:name/:action", stc.Pop3)
		r1.GET("/pop3/plt/:name/:action/*sub", stc.Pop3)

		r1.GET("/pop4/plt/:name/:action", stc.Pop4)
		r1.GET("/pop4/plt/:name/:action/*sub", stc.Pop4)

	}
	err = r.Run(":8888")
	ctl.Debug(err)

}
