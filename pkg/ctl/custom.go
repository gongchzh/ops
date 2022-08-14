package ctl

import (
	"archive/zip"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/smtp"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"github.com/admpub/mail"

	"github.com/kardianos/service"
)

func ReadToMap(path string) (map[string]bool, error) {
	var lines map[string]bool
	lines = make(map[string]bool)
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()
	err = ReadLine(path, func(line string) error {
		if _, ok := lines[line]; ok {
			err = errors.New("键重复了")
			return err
		} else {
			if line != "" {
				lines[line] = true
			}
		}
		return nil
	})
	return lines, nil
}

func IsInSlice(slice []string, arg ...string) bool {
	var isin int

	for _, v := range slice {
		for _, t := range arg {
			if v == t {
				isin += 1
			}
		}
	}
	if len(arg) == isin {
		return true
	}
	return false
}

func GetIni(path string) (map[string]map[string]string, error) {
	var (
		err        error
		regSection *regexp.Regexp
		regOption  *regexp.Regexp
		sectionTmp string
		key        string
		value      string
		con        map[string]map[string]string
	)
	regSection, _ = regexp.Compile("^\\[.*\\]$")
	regOption, _ = regexp.Compile("^[^;].*=.*")
	con = make(map[string]map[string]string)
	err = ReadLine(path, func(line string) error {
		line = strings.Replace(line, string([]byte{239, 187, 191}), "", 1)
		line = strings.Replace(line, "\r\n", "", -1)
		line = strings.Replace(line, "\n", "", -1)
		if regSection.MatchString(line) {
			if sectionTmp != "" {
				//				con[sectionTmp] = tmpMap
			}
			sectionTmp = strings.Split(strings.Split(line, "[")[1], "]")[0]
			con[sectionTmp] = make(map[string]string)
		}
		if regOption.MatchString(line) {
			key = strings.Split(line, "=")[0]
			if len(strings.Split(line, "=")) == 2 {
				value = strings.Split(line, "=")[1]
			} else if len(strings.Split(line, "=")) > 2 {
				value = ""
				for s, t := range strings.Split(line, "=") {
					if s == 0 {
						continue
					}
					if s == len(strings.Split(line, "="))-1 {
						value += t
					} else {
						value += t + "="
					}
				}
			}
			if sectionTmp == "" {
				return nil
			}
			con[sectionTmp][key] = value
		}
		return nil
	})
	return con, err
}

func ZipInfoList(path string) ([]os.FileInfo, error) {
	var fileInfos []os.FileInfo
	zipReader, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}
	for k, _ := range zipReader.File {
		fileInfos = append(fileInfos, zipReader.File[k].FileInfo())
	}
	zipReader.Close()
	return fileInfos, nil
}

func ZipFileExist(zipPath string, fileInfo os.FileInfo) bool {
	zipReader, err := zip.OpenReader(zipPath)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	for _, v := range zipReader.File {
		if fileInfo.IsDir() {
			if strings.Split(v.Name, "/")[0] == fileInfo.Name() {
				return true
			}
		}
		if v.FileInfo().Name() == fileInfo.Name() && v.FileInfo().IsDir() == fileInfo.IsDir() && v.FileInfo().Size() == fileInfo.Size() {
			return true
		}
	}
	zipReader.Close()
	return false
}

func ZipInfo(path string) ([]os.FileInfo, error) {
	var fileInfos []os.FileInfo
	zipReader, err := zip.OpenReader(path)

	if err != nil {
		return nil, err
	}
	for k, _ := range zipReader.File {
		fileInfos = append(fileInfos, zipReader.File[k].FileInfo())
	}
	zipReader.Close()
	return fileInfos, nil
}

func GetPathUnix(path string) string {
	if filepath.Separator == 92 {
		path = filepath.ToSlash(path)
	}
	if []byte(path)[len([]byte(path))-1] == 47 {
		path = string([]byte(path)[:len([]byte(path))-1])
	}
	return path
}

type ZipFiles struct {
	ZipPath  string
	FilePath string
}

func CheckZipFile(zipFiles *ZipFiles) bool {
	if CheckFile(zipFiles.FilePath) && CheckFile(zipFiles.ZipPath) {
		file, err := os.Lstat(zipFiles.FilePath)
		if err != nil {
			return false
		}
		if file.IsDir() {
			infoZip, err := ZipInfo(zipFiles.ZipPath)
			zipCount := 0
			if err != nil {
				return false
			}
			var zipSize, fileSize int64
			for _, t := range infoZip {
				zipSize += t.Size()
				zipCount += 1
			}
			infoFile, err := ListDirAll(zipFiles.FilePath)
			fileCount := 0
			for _, t := range infoFile {
				fileSize += t.Size()
				fileCount += 1
			}
			if zipSize == fileSize {
				return true
			}
			return false
		} else {
			infoZip, err := ZipInfo(zipFiles.ZipPath)
			if err != nil {
				return false
			}
			if len(infoZip) == 1 && file.Size() == infoZip[0].Size() {
				return true
			}
			return false
		}
	}
	return false
}

func RanderHtml(w http.ResponseWriter, r *http.Request, html string) error {
	var (
		t   *template.Template
		err error
	)

	t, err = template.ParseFiles(html) //html是具体的html文件路径
	if err != nil {
		return err
	}
	t.Execute(w, nil)
	return nil

}

func SendMail(user, passwd, host, to, header, body, mailType string) error {
	var (
		err          error
		content_type string
		msg          []byte
		send_to      []string
		auth         smtp.Auth
	)

	auth = smtp.PlainAuth("", user, passwd, strings.Split(host, ":")[0])
	if mailType == "html" {
		content_type = "Content-Type: text/" + mailType + "; charset=UTF-8"
	} else {
		content_type = "Content-Type: text/plain" + "; charset=UTF-8"
	}

	//msg = []byte("To: " + to + "\r\nFrom: " + user + ">\r\nSubject: " + header + "\r\n" + content_type + "\r\n\r\n" + body)
	msg = []byte("To: " + to + "\r\nSubject: " + header + "\r\n" + content_type + "\r\n\r\n" + body)

	send_to = strings.Split(to, ";")
	err = smtp.SendMail(host, auth, user, send_to, msg)
	return err
}
func SendMailSSL(user, passwd, host, to, header, body, mailType string) error {
	var (
		conf *mail.SMTPConfig
		err  error
		con  mail.SMTPClient
		m    mail.Mail
		h1   []string
		port int
	)
	h1 = strings.Split(host, ":")
	if len(h1) != 2 {
		return Errorf("host error")
	}
	port, err = strconv.Atoi(h1[1])
	if err != nil {
		return err
	}
	conf = &mail.SMTPConfig{
		Username: user,
		Password: passwd,
		Host:     h1[0],
		Port:     port,
		Secure:   "SSL",
	}
	con = mail.NewSMTPClient(conf)
	m = mail.NewMail()
	for _, v := range strings.Split(to, ";") {
		if v == "" {
			continue
		}
		m.AddTo(v)
	}

	m.AddFrom("龚纯振 <" + user + ">")
	m.AddSubject(header)
	if mailType == "html" {
		m.AddHTML(body)
	} else {
		m.AddText(body)
	}
	//m.Charset = "UTF-8"
	//	Debug(m.AddHeaders())
	//m.AddHeaders("charset=UTF-8")
	//return err
	err = con.Send(m)
	return err
}

func tlsDial(addr string) (*smtp.Client, error) {
	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		return nil, err
	}
	//分解主机端口字符串
	host, _, _ := net.SplitHostPort(addr)
	return smtp.NewClient(conn, host)
}

func sendMailTLS(addr string, auth smtp.Auth, from string,
	to []string, msg []byte) (err error) {

	//create smtp client
	con, err := tlsDial(addr)
	if err != nil {
		return err
	}
	defer con.Close()
	if auth != nil {
		if ok, _ := con.Extension("AUTH"); ok {
			if err = con.Auth(auth); err != nil {
				return err
			}
		}
	}
	if err = con.Mail(from); err != nil {
		return err
	}

	for _, addr := range to {
		if err = con.Rcpt(addr); err != nil {
			return err
		}
	}

	w, err := con.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return con.Quit()
}

func ChToNum(num string) int {
	var numMap = map[string]int{"一": 1, "二": 2, "三": 3, "四": 4, "五": 5, "六": 6, "七": 7, "八": 8, "九": 9}
	if _, ok := numMap[num]; ok {
		return numMap[num]
	}
	return 0
}

func GetInfo(host, url string, data interface{}) (interface{}, error) {
	var (
		err  error
		conn net.Conn
		bt   bytes.Buffer
	)
	conn, err = net.Dial("tcp", host)
	if err != nil {
		return data, err
	}
	_, err = conn.Write([]byte(url))
	if err != nil {
		return data, err
	}
	_, err = bt.ReadFrom(conn)
	if err != nil {
		return data, err
	}
	err = json.Unmarshal(bt.Bytes(), &data)
	if err != nil {
		return data, err
	}
	return data, err
}

func WinService(prg service.Interface, name ...string) error {
	var (
		srvConfig *service.Config
	)
	if len(name) == 3 {
		srvConfig = &service.Config{
			Name:        name[0], //服务显示名称
			DisplayName: name[1], //服务名称
			Description: name[2], //服务描述
		}

	} else if name != nil {
		return errors.New("service option error")
	} else {
		configPath := DirStand(filepath.Dir(os.Args[0])) + "/config.ini"
		config, err := GetIni(configPath)
		if err != nil {
			os.Stderr.WriteString("错误:加载配置文件失败\n")
			os.Exit(1)
		}
		srvConfig = &service.Config{
			Name:        config["main"]["servicename"], //服务显示名称
			DisplayName: config["main"]["displayname"], //服务名称
			Description: config["main"]["description"], //服务描述
		}
	}
	s, err := service.New(prg, srvConfig)
	if err != nil {
		return err
	}
	if len(os.Args) == 2 {
		if os.Args[1] == "install" {
			err = s.Install()
			if err != nil {
				fmt.Println("服务安装失败", err.Error())
			} else {
				fmt.Println("服务安装成功")
			}
			os.Exit(0)

		}
		if os.Args[1] == "remove" {
			err = s.Uninstall()
			if err != nil {
				fmt.Println("服务卸载失败", err.Error())
			} else {
				fmt.Println("服务卸载成功")
			}
			os.Exit(0)
		}
	}
	err = s.Run()
	return err
}

func StringSliceEnd(list []string) string {
	return list[len(list)-1]
}

var (
	Ecc1 = `-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIG3w1N3XL3D4D2S4/xgJW7CfxumBEd+Ngd6OLp6+vGyUoAoGCCqGSM49
AwEHoUQDQgAEYcTXpce+3Mq4io31oMAVTA4DiBLwThHaROKK7Jbzyro3Rdws/yNM
Hv0CRslfUbG0ncmx7aa7vC8mYSR0fe95cg==
-----END EC PRIVATE KEY-----`
)
