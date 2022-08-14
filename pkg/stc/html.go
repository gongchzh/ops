package stc

import (
	"ops/pkg/ctl"
	"ops/pkg/cfg"

	"github.com/gin-gonic/gin"
)

type HeaderHtml struct {
	DftState string
	State    cfg.HeaderHtml
}
type HeaderState struct {
	StateName string
	FullName  string
}

func PltHeader(c *gin.Context) {
	var (
		h HeaderHtml
	)
	if _, ok := cfg.PdtState[c.Param("plt_type")]; ok {
		h.DftState = c.Param("plt_type")
		ctl.Debug(h.DftState)
		h.State = cfg.GameHeaderHtml
		c.HTML(200, "plt.html", &h)
	}
}

func Login(c *gin.Context) {
	c.HTML(200, "login.html", nil)
}

type DftState struct {
	State     int
	StateName string
}

func Index(c *gin.Context) {
	ctl.Debug("start index")
	var (
		dftState DftState
	)
	dftState.State = cfg.Deploy.DftState
	dftState.StateName = cfg.GameState[cfg.Deploy.DftState].StateName
	c.HTML(200, "index.html", &dftState)
}
func Header(c *gin.Context) {
	ctl.Debug("start index")
	cookie, err := c.Cookie("user")
	ctl.Debug(cookie)
	ctl.Debug(err)
	ctl.Debug(c.GetString("user"))
	ctl.Debug(c.Keys)
	ctl.Debug(c.GetStringMap("user"))
	ctl.Debug(&c)
	if c.GetString("user") == "gongchunzheng" {
		c.HTML(200, "header.html", &cfg.GameHeaderHtml)

	} else {
		c.HTML(200, "game.html", &cfg.GameHeaderHtml)
	}
}
func Common(c *gin.Context) {
	ctl.Debug("start index")
	c.HTML(200, "common.html", &cfg.GameHeaderHtml)
}

func Pop(c *gin.Context) {
	ctl.Debug(c.Param("name"))
	ctl.Debug(c.Param("action"))
	c.HTML(200, "pop.html", nil)
}
func Pop2(c *gin.Context) {
	ctl.Debug(c.Param("name"))
	ctl.Debug(c.Param("action"))
	c.HTML(200, "pop2.html", nil)
}

func Pop3(c *gin.Context) {
	ctl.Debug(c.Param("name"))
	ctl.Debug(c.Param("action"))
	c.HTML(200, "pop3.html", nil)
}

func Pop4(c *gin.Context) {
	ctl.Debug(c.Param("name"))
	ctl.Debug(c.Param("action"))
	c.HTML(200, "pop4.html", nil)
}
func DisplayLeft(c *gin.Context) {
	c.HTML(200, "display_left.html", nil)
}
func DefaultHtml(c *gin.Context) {
	ctl.Debug("start default")
	c.HTML(200, "default.html", nil)
}
func DisplayHtml(c *gin.Context) {
	c.HTML(200, "display.html", &cfg.GameHeaderHtml)
}
