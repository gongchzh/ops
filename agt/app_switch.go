package agt

import (
	"ctl"
	"encoding/json"
	"io/ioutil"
	"os"
)

func GetNgxs() ([]NginxInfo, error) {
	var (
		ngxs []NginxInfo
	)
	b, err := ioutil.ReadFile(NginxInfoFile)
	if err != nil {

		return ngxs, err
	}
	err = json.Unmarshal(b, &ngxs)
	return ngxs, err
}
//要修改
func GetNgx(info string) (NginxInfo,error){
	var (
		ngx NginxInfo
	)
	b, err := ioutil.ReadFile(NginxInfoFile)
	if err != nil {
		return ngx,err
	}
	err = json.Unmarshal(b, &ngx)
	return ngx,err
}

func GetNginx(name string) (NginxInfo, error) {
	var (
		ngxs []NginxInfo
		ngx  NginxInfo
		err  error
	)
	ngxs, err = GetNgxs()
	if err != nil {
		return ngx, err
	}
	for _, v := range ngxs {

		if v.AppNginx.AppNginxName == name {
			ngx = v
		}
	}
	if ngx.AppNginx.AppNginxName == "" || ngx.AppNginx.AppNginxId == 0 {
		return ngx, ctl.Errorf("获取nginx信息异常")
	}
	return ngx, err
}

func AppSwitch() {
	var (
		ngx NginxInfo
		err error
	)
	if len(os.Args) != 5 {
		ctl.Panic(ctl.Errorf(""), "输入参数不正确")
	}
	ngx, err = GetNgx(os.Args[2])
	ctl.FatalErr(err)
	err = ngx.Switch(os.Args[3], os.Args[4])
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}
}

func (ngx *NginxInfo) Switch(md5, back string) error {
	var (
		err error
		app AppInfo  //要修改
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
