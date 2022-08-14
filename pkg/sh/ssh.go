package sh

import (
	"bytes"
	"ops/pkg/ctl"
	"io"
	"net"
	"ops/pkg/db"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type SSHHost struct {
	db.Host                     //继承主机库的结构
	HostName  string            //主机名
	Ip        string            //Ip地址
	User      string            //ssh账户
	Password  string            //ssh密码
	Config    *ssh.ClientConfig //ssh配置项
	SshClient *ssh.Client       //ssh连接
	Lock      *sync.Mutex       //同步锁
}

func (h *SSHHost) GetRemote() string {
	var (
		host db.Host
	)
	db.Db.Where("host_name=?", h.HostName).Get(&host)
	return strings.Split(host.Remote, ":")[0]
}
func (h *SSHHost) GetIp(name string) error {
	var (
		err  error
		h1   db.Host
		u    db.HostUser
		memu db.HostGroupUserMember
		memg db.HostGroupMember
		memh db.HostUserMember
	)
	_, err = db.Db.Where("host_name=?", name).Get(&h1)
	//ctl.Debug(h)
	if err != nil {
		ctl.Debug(name, h1)
		ctl.Debug(err)
		return err
	}
	h.HostName = h1.HostName
	h.HostId = h1.HostId
	if runtime.GOOS == "linux" {
		h.Ip = h1.Local
	} else {
		h.Ip = h1.Remote
	}
	_, err = db.Db.Where("host_id=?", h.HostId).Get(&memh)
	if err != nil {
		ctl.Debug(err)
	}
	if memh.UserId == 0 {
		_, err = db.Db.Where("host_id=?", h.HostId).Get(&memg)
		if err != nil {
			ctl.Debug(err)
		}
		if memg.HostGroupId == 0 {
			ctl.Debug(h.HostName, h.HostId)
			ctl.Debug(h1)
			ctl.Debug(name)

			ctl.Debug(memh)
			ctl.Debug(memg)
			ctl.Debug(h.HostName, "没有用户")
			os.Exit(1)
			return ctl.Errorf(h.HostName + "没有用户")
		}
		db.Db.Where("host_group_id=?", memg.HostGroupId).Get(&memu)
		db.Db.Where("user_id=?", memu.UserId).Get(&u)
	} else {
		db.Db.Where("user_id=?", memh.UserId).Get(&u)
	}
	if u.UserId == 0 {
		ctl.FatalErr(err)
	}
	h.User = u.Name
	h.Password = db.Aes.Decode(u.Password)
	h.Config = &ssh.ClientConfig{
		User: h.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(h.Password),
		},
		Timeout: time.Hour * 10,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	return err
}

func (h *SSHHost) Error(err interface{}) string {
	switch value := err.(type) {
	case error:
		if value != nil {
			return h.HostName + " " + value.Error()
		}
	case string:
		return h.HostName + " " + value
	}
	return h.HostName + " is null"
}
func (h *SSHHost) Keep(sleep time.Duration) {
	var (
		err error
		i   int
	)
	for {
		time.Sleep(time.Second * 10)
		_, _, err = h.SshClient.SendRequest("", false, nil)
		if err != nil {
			ctl.Debug(h.Error(err))
			ctl.Log.Error(h.Error(err))
			h.State = 0
			if i%15 == 1 {
				//	h.SshClient.Close()
				h.SshDial()
			}
		}
		time.Sleep(time.Second * sleep)
		i++
	}
}
func (h *SSHHost) ConFail() bool {
	var (
		err error
	)
	if h == nil {
		return true
	}
	if h.SshClient == nil {
		return true
	}
	_, _, err = h.SshClient.SendRequest("", false, nil)

	if err == nil {
		return false
	}
	return true
}

func (h *SSHHost) SshDial() error {
	var (
		err error
	)
	//	ctl.Debug(h)
	//	ctl.Debug(h.HostName)
	err = h.GetIp(h.HostName)
	if err != nil {
		return err
	}
	h.Lock.Lock()
	defer h.Lock.Unlock()
	h.SshClient, err = ssh.Dial("tcp", h.Ip, h.Config)
	//	h.sshClient.
	if err != nil {
		h.State = 0
		return err
	}
	h.State = 1
	return err
}

func (h *SSHHost) SshCmd(cmd string) string {
	var (
		err    error
		cmdout []byte
		sess   *ssh.Session
	)

	if h.ConFail() {
		ctl.Debug(h)
		if h == nil {
			return ""
		} else {
			return "错误:" + h.HostName + "连接异常"
		}
	}
	sess, err = h.SshClient.NewSession()
	if err != nil {
		ctl.Log.Debug(h.Error(err))
		return err.Error()
	}
	defer sess.Close()
	cmdout, err = sess.CombinedOutput(cmd)
	if err != nil {
		ctl.Log.Debug(h.HostName, err)
		ctl.Log.Debug(h.HostName, cmd)
		ctl.Log.Debug(h.HostName, string(cmdout))
	}
	if err != nil && cmdout == nil {
		ctl.Log.Debug(h.HostName, err)
		ctl.Log.Debug(h.HostName, string(cmdout))
		return err.Error()
	}

	return string(cmdout)
}
func (h *SSHHost) RunCmd(cmd string) (string, error) {
	var (
		err    error
		cmdout []byte
		sess   *ssh.Session
	)

	if h.ConFail() {
		ctl.Debug(h)
		if h == nil {
			return "", ctl.Errorf("连接为空")
		} else {
			return "", ctl.Errorf(h.HostName + "连接异常")
		}
	}
	sess, err = h.SshClient.NewSession()
	if err != nil {
		ctl.Log.Debug(h.Error(err))
		return "", err
	}
	defer sess.Close()
	cmdout, err = sess.CombinedOutput(cmd)

	return string(cmdout), err
}

func (h *SSHHost) SshHostChCmdErr(ch chan db.Result, cmd string) {
	ch <- h.SshHostCmdErr(cmd)
}

func (h *SSHHost) SshHostCmdErr(cmd string) (res db.Result) {
	var (
		err          error
		hostFormat   string
		cmdout, line []byte
		out, newOut  bytes.Buffer
		sess         *ssh.Session
	)
	if h.ConFail() {
		ctl.Debug(h)
		if h == nil {
			res.Err = ctl.Errorf("链接为空")
			return
		} else {
			res.HostName = h.HostName
			res.HostId = h.HostId
			res.Err = ctl.Errorf("[" + h.HostName + "]" + "错误:连接异常")
			return
		}
	}
	hostFormat = "[" + h.HostName + "]"
	res.HostName = h.HostName
	res.HostId = h.HostId
	sess, err = h.SshClient.NewSession()
	if err != nil {
		res.Err = ctl.Errorf(hostFormat + "错误:" + err.Error())
		return
	}
	defer sess.Close()
	cmdout, err = sess.CombinedOutput(cmd)
	if err != nil && cmdout == nil {
		res.Err = ctl.Errorf(hostFormat + "错误:" + err.Error())
		return
	}
	if err != nil {
		res.Err = ctl.Errorf(hostFormat + ":错误" + err.Error() + "\n")
	}
	out.Write(cmdout)
	for {
		line, err = out.ReadBytes(10)
		if err != nil {
			if err == io.EOF {
				break
			}
			if res.Err == nil {
				res.Err = err
			} else {
				res.Err = ctl.Errorf(res.Err.Error()+",%s", err)
			}
		}
		newOut.WriteString(hostFormat)
		newOut.Write(line)
	}
	res.Result = newOut.String()
	return
	//	return strings.Replace(newOut.String(), "\n", "</br>", -1)
}

func (h *SSHHost) SshAppChCmd(ch chan db.AppResult, srv AppServer, cmd string) {
	ch <- h.SshAppCmd(srv, cmd)
}

func (h *SSHHost) SshAppCmd(srv AppServer, cmd string) (res db.AppResult) {
	var (
		err          error
		hostFormat   string
		cmdout, line []byte
		out, newOut  bytes.Buffer
		sess         *ssh.Session
	)
	res.AppId = srv.AppId
	res.AppServerId = srv.AppServerId
	if h.ConFail() {
		ctl.Debug(h)
		if h == nil {
			res.Err = ctl.Errorf("链接为空")
			return
		} else {
			res.HostName = h.HostName
			res.HostId = h.HostId
			res.Err = ctl.Errorf("[" + h.HostName + "]" + "错误:连接异常")
			return
		}
	}
	hostFormat = "[" + h.HostName + "]"
	res.HostName = h.HostName
	res.HostId = h.HostId
	sess, err = h.SshClient.NewSession()
	if err != nil {
		res.Err = ctl.Errorf(hostFormat + "错误:" + err.Error())
		return
	}
	defer sess.Close()
	ctl.Debug(cmd)
	cmdout, err = sess.CombinedOutput(cmd)
	if err != nil && cmdout == nil {
		res.Err = ctl.Errorf(hostFormat + "错误:" + err.Error())
		return
	}
	if err != nil {
		res.Err = ctl.Errorf(hostFormat + ":错误" + err.Error() + "\n")
	}
	out.Write(cmdout)
	for {
		line, err = out.ReadBytes(10)
		if err != nil {
			if err == io.EOF {
				break
			}
			if res.Err == nil {
				res.Err = err
			} else {
				res.Err = ctl.Errorf(res.Err.Error()+",%s", err)
			}
		}
		newOut.WriteString(hostFormat)
		newOut.Write(line)
	}
	res.Result = newOut.String()
	return
	//	return strings.Replace(newOut.String(), "\n", "</br>", -1)
}

func (h *SSHHost) SshHostCmd(cmd string) string {
	var (
		err          error
		hostFormat   string
		cmdout, line []byte
		out, newOut  bytes.Buffer
		sess         *ssh.Session
	)
	hostFormat = "[" + h.HostName + "]"
	if h.ConFail() {
		ctl.Debug(h)
		if h == nil {
			return ""
		} else {
			return hostFormat + "错误:连接异常"
		}
	}
	sess, err = h.SshClient.NewSession()
	if err != nil {
		return hostFormat + "错误:" + err.Error()
	}
	defer sess.Close()
	cmdout, err = sess.CombinedOutput(cmd)
	if err != nil && cmdout == nil {
		return hostFormat + "错误:" + err.Error()
	}
	if err != nil {
		newOut.WriteString(hostFormat + ":错误" + err.Error() + "\n")
	}
	out.Write(cmdout)
	for {
		line, err = out.ReadBytes(10)
		if err != nil {
			break
		}
		newOut.WriteString(hostFormat)
		newOut.Write(line)
	}
	return newOut.String()
	//	return strings.Replace(newOut.String(), "\n", "</br>", -1)
}

func (h *SSHHost) SshHostHtmlCmd(cmd string) string {
	var (
		err          error
		hostFormat   string
		cmdout, line []byte
		out, newOut  bytes.Buffer
		sess         *ssh.Session
	)
	ctl.Debug(h)
	hostFormat = "[" + h.HostName + "]"
	if h.ConFail() {
		ctl.Debug(h)
		if h == nil {
			return ""
		} else {
			ctl.Debug(h.HostId, h.SshClient)
			return "错误:" + h.HostName + "连接异常"
		}
	}
	sess, err = h.SshClient.NewSession()
	if err != nil {
		return h.HostName + ":" + err.Error()
	}
	defer sess.Close()
	cmdout, err = sess.CombinedOutput(cmd)
	if err != nil && cmdout == nil {
		return hostFormat + ":" + err.Error()
	}
	out.Write(cmdout)
	for {
		line, err = out.ReadBytes(10)
		if err != nil {
			break
		}
		newOut.WriteString(hostFormat)
		newOut.Write(line)
	}
	//	return newOut.String()
	return strings.Replace(newOut.String(), "\n", "</br>", -1)
}

func (h *SSHHost) SshChCmd(ch chan string, cmd string) {
	ch <- h.SshCmd(cmd)
}
func (h *SSHHost) SshHostChCmd(ch chan string, cmd string) {
	ch <- h.SshHostCmd(cmd)
}
func (h *SSHHost) SshHostHtmlChCmd(ch chan string, cmd string) {
	ch <- h.SshHostHtmlCmd(cmd)
}
func (h *SSHHost) GetInfo(ch chan []byte) {
	var (
		err error
		s   *ssh.Session
		out []byte
	)
	defer func() {
		ch <- out
	}()
	if h.ConFail() {
		ctl.Debug(h.HostName)
		return
	}
	s, err = h.SshClient.NewSession()
	if err != nil {
		ctl.Debug(h.Error(err))
		ctl.Log.Error(h.Error(err))
		return
	}
	defer s.Close()
	out, err = s.CombinedOutput("/home/app/newgame/gameinfo")
	if err != nil {
		ctl.Log.Error(h.Error(err))
		return
	}
}

func (h *SSHHost) WriteFile(p string, b []byte) (int, error) {
	var (
		err    error
		line   int
		client *sftp.Client
		file   *sftp.File
	)
	if h.ConFail() {
		return 0, ctl.Errorf("错误:" + h.HostName + "连接异常")
	}
	client, err = sftp.NewClient(h.SshClient)
	if err != nil {
		return 0, err
	}
	defer client.Close()
	file, err = client.OpenFile(p, os.O_CREATE|os.O_WRONLY|os.O_TRUNC)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	line, err = file.Write(b)
	return line, err
}
func (h *SSHHost) Stat(p string) (os.FileInfo, error) {
	var (
		err    error
		info   os.FileInfo
		client *sftp.Client
	)
	if h.ConFail() {
		return nil, ctl.Errorf("错误:" + h.HostName + "连接异常")
	}

	client, err = sftp.NewClient(h.SshClient)
	if err != nil {
		return nil, err
	}
	defer client.Close()
	info, err = client.Stat(p)
	return info, err
}

func (h *SSHHost) NewSftpClient() (*sftp.Client, error) {
	return sftp.NewClient(h.SshClient)
}
func (h *SSHHost) OpenFile(p string, f int) (*sftp.File, error) {
	var (
		err    error
		client *sftp.Client
	)
	if h.ConFail() {
		return nil, ctl.Errorf("错误:" + h.HostName + "连接异常")
	}
	client, err = sftp.NewClient(h.SshClient)
	if err != nil {
		return nil, err
	}
	defer client.Close()
	return client.OpenFile(p, f)
}
func (h *SSHHost) Open(p string) (*sftp.File, error) {
	return h.OpenFile(p, os.O_RDONLY)
}

func (h *SSHHost) Sync(src, dst string) (err error) {
	if h.ConFail() {
		err = ctl.Errorf(h.HostName + "ssh连接异常")
		return
	}
	sinfo, err := os.Stat(src)
	if err != nil {
		return
	}
	sess, err := sftp.NewClient(h.SshClient)
	if err != nil {
		return
	}
	defer sess.Close()
	dinfo, err := sess.Stat(dst)
	if err != nil {
		if _, err = sess.Stat(ctl.UnixDir(dst)); err != nil {
			err = sess.Mkdir(ctl.UnixDir(dst))
			if err != nil {
				return
			}
		}
		return h.Copy(src, dst)
	}
	if sinfo.Size() == dinfo.Size() && sinfo.ModTime() == dinfo.ModTime() {
		return
	}
	return h.Copy(src, dst)
}
func (h *SSHHost) Copy(src, dst string) (err error) {
	var (
		bak string
	)
	if h.ConFail() {
		err = ctl.Errorf(h.HostName + "ssh连接异常")
		return
	}
	bak = dst + ".bak." + ctl.Nowf()
	info, err := os.Stat(src)
	if err != nil {
		return
	}
	sess, err := sftp.NewClient(h.SshClient)
	if err != nil {
		return
	}
	defer sess.Close()
	fs, err := os.Open(src)
	if err != nil {
		return
	}
	defer fs.Close()
	_, err = sess.Stat(dst)
	if err == nil {
		err = sess.Rename(dst, bak)
		if err != nil {
			return
		}
	}
	fd, err := sess.Create(dst)
	if err != nil {
		return
	}
	defer fd.Close()
	ln, err := io.Copy(fd, fs)
	if err != nil {
		return
	}
	if ln != info.Size() {
		err = ctl.Errorf("copy file size diff")
		return
	}
	err = sess.Chtimes(dst, info.ModTime(), info.ModTime())
	if err != nil {
		return
	}
	err = sess.Chmod(dst, info.Mode())
	if err != nil {
		return
	}
	err = sess.Remove(bak)
	return
}

func (h *SSHHost) Sftp(ch chan string, files []db.HtmlFile, backup bool) {
	var (
		err    error
		file   *sftp.File
		line   int
		client *sftp.Client
	)
	if h.ConFail() {
		ch <- "错误:" + h.HostName + "连接异常"
		return
	}
	if err != nil {
		ch <- h.HostName + ":" + err.Error()
		return
	}
	ctl.Debug(h.HostName)
	client, err = sftp.NewClient(h.SshClient)
	if err != nil {
		ch <- h.HostName + ":" + err.Error()
		return
	}
	defer client.Close()
	for _, v := range files {
		_, err = client.Stat(v.Dst)
		if err != nil {
			ch <- h.HostName + ":" + err.Error()
			return
		}
		if backup {
			_, err = client.Stat(v.PathName)
			ctl.Debug(err)
			if err == nil {
				err = client.Rename(v.PathName, v.PathName+".bak."+time.Now().Format("2006-01-02_15:04:05"))
				ctl.Debug(err)
				if err != nil {
					ch <- h.HostName + ":" + err.Error()
					return
				}
			}
		}
		file, err = client.OpenFile(v.PathName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE)
		if err != nil {
			ch <- h.HostName + ":" + err.Error()
			return
		}
		defer file.Close()
		line, err = file.Write(v.Data.Bytes())
		if err != nil {
			ch <- h.HostName + ":" + err.Error()
			return
		}
		if int64(line) != v.Size {
			ch <- h.HostName + ":size error"
			return
		}
	}
	ch <- h.HostName + ":上传成功\n"
	return
}
