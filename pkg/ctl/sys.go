package ctl

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"strconv"
	"time"
)

func Errorf(format string, a ...interface{}) error {
	return fmt.Errorf(format, a...)
}
func ErrorAf(format string, err error, a ...interface{}) error {
	var (
		a1 []interface{}
	)
	if err != nil {
		a1 = append(a1, err)
		a1 = append(a1, a...)
		return fmt.Errorf(format, a1...)
	}
	//a1 = append(a1, fmt.Errorf(""))
	a1 = append(a1, "")
	a1 = append(a1, a...)
	return fmt.Errorf(format, a1...)
}

func Pkill(exe string) error {
	var (
		process []*os.Process
		err     error
	)
	process, err = FindProcess(exe)
	if err != nil {
		return err
	}
	for _, v := range process {
		err = v.Kill()
		if err != nil {
			return err
		}
	}
	return nil
}
func Printerrln(d ...interface{}) {
	fmt.Fprintln(os.Stderr, d...)
}
func Printerr(d ...interface{}) {
	fmt.Fprint(os.Stderr, d...)
}
func Sprintf(format string, d ...interface{}) string {
	return fmt.Sprintf(format, d...)
}
func Printf(format string, d ...interface{}) (int, error) {
	return fmt.Printf(format+"\n", d...)
}
func Print(d ...interface{}) (int, error) {
	return fmt.Print(d...)
}

func Println(d ...interface{}) (int, error) {
	return fmt.Println(d...)
}

func FindProcess(exe string) ([]*os.Process, error) {
	var (
		process       []*os.Process
		processTmp    *os.Process
		err           error
		d             *os.File
		pid           int
		program, name string
		fis           []os.FileInfo
	)
	d, err = os.Open("/proc")
	if err != nil {
		return nil, err
	}
	defer d.Close()
	for {
		fis, err = d.Readdir(10)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		for _, fi := range fis {
			processTmp = nil
			if !fi.IsDir() {
				continue
			}
			name = fi.Name()
			if name[0] < '0' || name[0] > '9' {
				continue
			}
			pid, err = strconv.Atoi(name)
			if err != nil {
				continue
			}
			program, _ = os.Readlink(fmt.Sprintf("/proc/%d/exe", pid))
			if program == exe {
				processTmp, err = os.FindProcess(pid)
				if err != nil {
					return nil, err
				}
				process = append(process, processTmp)
			}
		}
	}
	return process, nil
}

func GetIp() ([]string, error) {
	var (
		ip  []string
		err error
	)
	addr, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}
	for _, v := range addr {
		if ipnet, ok := v.(*net.IPNet); ok {
			if ipnet.IP.IsLoopback() || ipnet.IP.IsLinkLocalUnicast() || ipnet.IP.To4() == nil || !ipnet.IP.IsGlobalUnicast() {
				continue
			}
			ip = append(ip, ipnet.IP.String())

		}
	}
	return ip, err
}

func GetOutIp() ([]string, error) {
	var (
		ip  []string
		err error
	)
	addr, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	for _, v := range addr {
		if ipnet, ok := v.(*net.IPNet); ok {
			if ipnet.IP.IsLoopback() || ipnet.IP.IsLinkLocalUnicast() || regLocal.MatchString(ipnet.IP.String()) || ipnet.IP.To4() == nil || !ipnet.IP.IsGlobalUnicast() {
				continue
			}
			ip = append(ip, ipnet.IP.String())

		}
	}
	return ip, err
}

func IsLocalIp(ip string) bool {
	if !RegIp.MatchString(ip) {
		return false
	}
	if ip[:3] == "10." {
		return true
	}
	if ip[:4] == "172." {
		return true
	}
	if ip[:8] == "192.168." {
		return true
	}

	return false
}

func InString(sl []string, s string) bool {
	for _, v := range sl {
		if v == s {
			return true
		}
	}
	return false
}

func GetLocalIp() (string, error) {
	ipl := ""
	ips, err := GetIp()
	if err != nil {
		return ipl, err
	}

	for _, v := range ips {
		if ipl != "" && IsLocalIp(v) {
			return ipl, Errorf("内网IP多于一个")
		}
		if IsLocalIp(v) {
			ipl = v
		}

	}
	if ipl == "" {
		return ipl, Errorf("内网IP不存在")
	}
	return ipl, err
}

func Run(str string, arg ...string) (string, error) {
	var (
		cmd            *exec.Cmd
		err            error
		stdout, stderr bytes.Buffer
		out            bytes.Buffer
	)
	cmd = exec.Command(str, arg...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	out.Write(stdout.Bytes())
	out.Write(stderr.Bytes())
	if len(stderr.Bytes()) != 0 && err == nil {
		err = Errorf(stderr.String())
	}
	return out.String(), err
}

func RunStart(str string, arg ...string) (string, error) {
	var (
		cmd            *exec.Cmd
		err            error
		stdout, stderr bytes.Buffer
		out            bytes.Buffer
	)
	cmd = exec.Command(str, arg...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Start()
	time.Sleep(time.Second)
	out.Write(stdout.Bytes())
	out.Write(stderr.Bytes())
	return out.String(), err
}

func Panic(err error, errStr ...interface{}) {
	if err != nil {
		fmt.Fprint(os.Stderr, "[错误]")
		fmt.Fprint(os.Stderr, errStr...)
		fmt.Fprint(os.Stderr, err)
		os.Stderr.WriteString("\n")
		os.Exit(1)
	}
}
