package al

import (
	"ops/pkg/ctl"
	"ops/pkg/sh"
)

func (srv *AppServer) ShAppServer() sh.AppServer {
	ctl.Debug(sh.AppServer(*srv))
	return sh.AppServer(*srv)
}
