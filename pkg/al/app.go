package al

import (
	"ops/pkg/ctl"
	"strings"
)

func (app *App) NewBack() string {
	var (
		basep string
		back  string
	)
	switch app.AppTypeId {
	case 1:
		if app.BasePath[len(app.BasePath)] == '/' {
			basep = app.BasePath[:len(app.BasePath)-1]
		} else {
			basep = app.BasePath
		}
		bases := strings.Split(basep, "/")
		back = bases[len(bases)-1]
		back += ".bak." + ctl.NowMinute()
	case 2:
		back = app.AppProgram + ".bak." + ctl.NowMinute()
	}
	return back
}
func (app *App) ParseScript(script string) string {
	script = strings.ReplaceAll(script, "{Port}", ctl.Itoa(app.Port))
	script = strings.ReplaceAll(script, "{AppName}", app.AppName)
	script = strings.ReplaceAll(script, "{port}", ctl.Itoa(app.Port))
	return script
}
