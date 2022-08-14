package agt

import (
	"ctl"
	"io"
	"ops/pkg/al"
	"os"
	"strings"

	"github.com/gogf/gf/os/gproc"
)

type AppInfo al.AppInfo

type AppList al.AppList

type NginxInfo al.NginxInfo

func AppUpdate() {
	var (
		app AppInfo
		err error
	)
	if len(os.Args) != 5 {
		ctl.Panic(ctl.Errorf(""), "输入参数不正确")
	}
	app, err = GetApp(os.Args[2])
	ctl.FatalErr(err)
	err = app.Update(os.Args[3], os.Args[4])
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}
}

func (app *AppInfo) GetNewMd5() string {
	md5, _ := ctl.Md5sum(app.UpdatePath + "/" + app.UpdateProgram)
	return md5
}
func (app *AppInfo) GetOldMd5() string {
	md5, _ := ctl.Md5sum(app.ProgramDir + "/" + app.AppProgram)
	return md5
}

func (app *AppInfo) Printf(format string, d ...interface{}) (int, error) {
	return ctl.Printf(app.AppName+" "+ctl.Itoa(app.AppServerId)+":"+format, d...)
}

func (app *AppInfo) CheckOpt() error {
	var (
		err error
	)
	if !ctl.CheckDir(app.UpdatePath) {
		return ctl.Errorf("检查应用信息:应用更新目录不存在")
	}
	if !ctl.CheckDir(app.BasePath) {
		return ctl.Errorf("检查应用信息:主目录不存在")
	}
	if !ctl.CheckDir(app.BackDir) {
		return ctl.Errorf("检查应用信息:备份目录不存在")
	}
	if !ctl.CheckDir(app.ProgramDir) {
		return ctl.Errorf("检查应用信息:应用上级目录不存在")
	}
	if !ctl.CheckFile(app.ProgramDir + "/" + app.AppProgram) {
		return ctl.Errorf("检查应用信息:应用文件不存在")
	}
	switch app.AppTypeId {
	case 1:
	case 2:

	}
	app.Printf("检查应用配置项成功")
	return err
}

func (app *AppInfo) CheckUpdate(md5 string, back string) error {
	var (
		err error
	)
	if len(md5) != 32 {
		return ctl.Errorf("输入md5长度不对")
	}
	err = app.CheckOpt()
	if err != nil {
		return err
	}
	md5n := app.GetNewMd5()
	if md5n != md5 {
		return ctl.Errorf("输入md5与新md5不一致")
	}
	app.Printf("检查更新md5成功,md5:%s", md5n)
	if back == "" {
		return ctl.Errorf("检查备份格式异常")
	}
	return err
}

func (app *AppInfo) Kill() error {
	var (
		res string
		err error
	)
	pid := app.GetProcess()
	if pid != 0 {
		app.Printf("开始杀掉旧进程,进程ID:%d", pid)
		res, err = ctl.Run("kill", "-9", ctl.Itoa(pid))
		if err != nil {
			return ctl.Errorf("杀掉进程异常:%s\n%s", res, err.Error())
		}
		app.Printf("杀掉旧进程成功,进程ID:%d", pid)
	} else {
		app.Printf("旧进程未找到")
	}
	return err
}

func (app *AppInfo) Back(back string) error {
	var (
		err error
	)
	app.Printf("开始备份旧应用...")
	switch {
	case app.AppTypeId == 1:
		err = os.Rename(app.BasePath+"/"+app.AppProgram, app.BackDir+"/"+back)
		if err != nil {
			return ctl.Errorf("备份当前应用:%s", err.Error())
		}
	case app.AppTypeId == 2:
		err = os.Rename(app.BasePath, app.BackDir+"/"+back)
		if err != nil {
			return ctl.Errorf("备份当前应用:%s", err.Error())
		}
	case app.AppTypeId == 3:
		err = os.Rename(app.ProgramDir+"/"+app.AppProgram, app.BackDir+"/"+back)
		if err != nil {
			return ctl.Errorf("备份当前应用:%s", err.Error())
		}
	}
	app.Printf("备份旧应用成功")
	return err
}

func (app *AppInfo) Deploy() error {
	var (
		err error
	)
	app.Printf("开始部署新应用...")
	switch app.AppTypeId {
	case 1:
		fn, err := os.Create(app.BasePath + "/" + app.AppProgram)
		if err != nil {
			return ctl.Errorf("新建程序文件:%s", err.Error())
		}
		defer fn.Close()
		fos, err := os.Stat(app.UpdatePath + "/" + app.UpdateProgram)
		if err != nil {
			return ctl.Errorf("获取更新文件信息:%s", err.Error())
		}
		fo, err := os.Open(app.UpdatePath + "/" + app.UpdateProgram)
		if err != nil {
			return ctl.Errorf("打开更新文件:%s", err.Error())
		}
		defer fo.Close()
		ln, err := io.Copy(fn, fo)
		if err != nil {
			return ctl.Errorf("复制更新文件:%s", err.Error())
		}
		if ln != fos.Size() {
			return ctl.Errorf("复制更新文件大小不一致,原大小:%d,新大小:%d", fos.Size(), ln)
		}
		err = os.Chmod(app.BasePath+"/"+app.AppProgram, 0755)
		if err != nil {
			return ctl.Errorf("授权更新文件:" + err.Error())
		}
	case 2:
		err = os.Mkdir(app.BasePath, 0755)
		if err != nil {
			return ctl.Errorf("创建tomcat应用目录:%s", err.Error())
		}
		if app.BasePath != app.ProgramDir {
			err = os.MkdirAll(app.ProgramDir, 0755)
			if err != nil {
				return ctl.Errorf("创建war包上级目录:%s", err.Error())
			}
		}
		app.Printf("开始复制更新文件,源路径:%s,新路径:%s", app.UpdatePath+"/"+app.UpdateProgram, app.ProgramDir+"/"+app.AppProgram)
		fn, err := os.Create(app.ProgramDir + "/" + app.AppProgram)
		if err != nil {
			return ctl.Errorf("新建程序文件:%s", err.Error())
		}
		defer fn.Close()
		fos, err := os.Stat(app.UpdatePath + "/" + app.UpdateProgram)
		if err != nil {
			return ctl.Errorf("获取更新文件信息:%s", err.Error())
		}
		fo, err := os.Open(app.UpdatePath + "/" + app.UpdateProgram)
		if err != nil {
			return ctl.Errorf("打开更新文件:%s", err.Error())
		}
		defer fo.Close()
		ln, err := io.Copy(fn, fo)
		if err != nil {
			return ctl.Errorf("复制更新文件:%s", err.Error())
		}
		if ln != fos.Size() {
			return ctl.Errorf("复制更新文件大小不一致,原大小:%d,新大小:%d", fos.Size(), ln)
		}
	}
	app.Printf("部署新应用成功")
	return err
}

func (app *AppInfo) Start() error {
	var (
		err error
	)
	app.Printf("开始启动应用...")
	switch app.AppTypeId {
	case 1:
		cmd := "cd " + app.BasePath + " && " + app.BasePath + "/" + app.AppProgram + "  -d=true"
		res, err := gproc.ShellExec(cmd)
		ctl.Printf(res)
		if err != nil {
			return ctl.Errorf("重启脚本执行异常:%s", err.Error())
		}
	case 2:
		dirs := strings.Split(app.BasePath, "/")
		res, err := ctl.RunStart("/bin/bash", TomcatDir+dirs[len(dirs)-1]+"/bin/startup.sh")
		ctl.Printf(res)
		if err != nil {
			return ctl.Errorf("重启脚本执行异常:%s", err.Error())
		}
	}
	app.Printf("应用启动脚本己执行")
	return err
}

func (app *AppInfo) Update(md5, back string) error {
	var (
		err error
	)
	err = app.CheckUpdate(md5, back)
	if err != nil {
		return err
	}
	err = app.Kill()
	if err != nil {
		return err
	}
	err = app.Back(back)
	if err != nil {
		return ctl.Errorf("备份当前应用:%s", err.Error())
	}
	err = app.Deploy()
	if err != nil {
		return err
	}
	err = app.Start()
	if err != nil {
		return err
	}
	md5n := app.GetOldMd5()
	if md5n != md5 {
		return ctl.Errorf("更新文件md5与输入不一致")
	}
	app.Printf("升级后检查更新文件MD5成功 新MD5:%s", md5n)
	return err
}
