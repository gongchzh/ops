package ctl

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"runtime"
	"runtime/debug"
	"strings"
	"time"
)

var (
	fileErr, fileAll               *os.File
	LogAll, LogErr                 Logger
	Log                            Logger
	Level                          int
	allLog, errLog, CurDay, prefix string
	regLocal, _                    = regexp.Compile("^10\\.")
	Env                            bool
	End                            string
	RegGb, _                       = regexp.Compile("^GB")
	RegNum, _                      = regexp.Compile("^[0-9]+$")
	RegNumStr, _                   = regexp.Compile("[0-9]+")
	RegIso, _                      = regexp.Compile("^ISO-8859")
	RegIp, _                       = regexp.Compile("[1-2]?[0-9][0-9]?\\.[1-2]?[0-9][0-9]?\\.[1-2]?[0-9][0-9]?\\.[1-2]?[0-9][0-9]?")
	RegIpOn, _                     = regexp.Compile("^[1-2]?[0-9][0-9]?\\.[1-2]?[0-9][0-9]?\\.[1-2]?[0-9][0-9]?\\.[1-2]?[0-9][0-9]?$")
)

func Fprintln(w io.Writer, str ...interface{}) {
	Log.Info(str...)
	fmt.Fprintln(w, str...)
}
func Fprinterr(w io.Writer, str ...interface{}) {
	Fprintln(w, str...)
	LogErr.Println(str...)
}

type Logger struct {
	*log.Logger
}

func (l *Logger) FatalErr(err error) {
	if err != nil {
		var (
			file string
			line int
			arh  []interface{}
		)
		_, file, line, _ = runtime.Caller(1)
		file = strings.Split(file, "/")[len(strings.Split(file, "/"))-1]
		arh = append(arh, "[error]", file, line)
		arh = append(arh, " ", err.Error(), End)
		l.Print(arh...)
		os.Exit(1)
	}
}
func FatalErr(err error) {
	if err != nil {
		var (
			file string
			line int
			arh  []interface{}
		)
		_, file, line, _ = runtime.Caller(1)
		file = strings.Split(file, "/")[len(strings.Split(file, "/"))-1]
		arh = append(arh, Now(), "[error]", file, line)
		arh = append(arh, " ", err.Error(), "\n")
		fmt.Print(arh...)
		os.Exit(1)
	}
}

func (l *Logger) Fatal(arg ...interface{}) {
	if Level == 4 {
		return
	}
	var (
		file string
		line int
		arh  []interface{}
	)
	arh = append(arh, "[fatal]")
	arh = append(arh, arg...)
	for k, v := range []int{1, 2, 3, 4, 5, 6} {
		_, file, line, _ = runtime.Caller(v)
		file = strings.Split(file, "/")[len(strings.Split(file, "/"))-1]
		if k == 5 {
			arh = append(arh, file, line, " ")
		} else {
			arh = append(arh, file, line, End)
		}
		//fmt.Println([]byte(file))
	}

	//arh = append(arh, arg...)
	l.Print(arh...)
}

func (l *Logger) Debug(arg ...interface{}) {
	if Level == 4 {
		return
	}
	var (
		file string
		line int
		arh  []interface{}
	)
	_, file, line, _ = runtime.Caller(1)
	file = strings.Split(file, "/")[len(strings.Split(file, "/"))-1]
	arh = append(arh, "[debug]", file, line, " ")

	arh = append(arh, arg...)
	arh = append(arh, End)
	l.Print(arh...)
}

func (l *Logger) Stack(arg ...interface{}) {
	var (
		arh []interface{}
	)
	arh = append(arh, "[track]", string(debug.Stack()))
	arh = append(arh, arg...)
	arh = append(arh, End)
	l.Print(arh...)
}
func (l *Logger) Track(arg ...interface{}) {
	if Level == 4 {
		return
	}
	var (
		file1, file2 string
		line1, line2 int
		arh          []interface{}
	)
	_, file1, line1, _ = runtime.Caller(1)
	_, file2, line2, _ = runtime.Caller(2)
	arh = append(arh, "[track]", file1, line1, " ", file2, line2, " ")
	for _, v := range arg {
		arh = append(arh, v, " ")
	}
	//	arh = append(arh, arg...)
	arh = append(arh, End)

	l.Print(arh...)
}

func (l *Logger) Info(arg ...interface{}) {
	var (
		file string
		line int
		arh  []interface{}
	)
	_, file, line, _ = runtime.Caller(1)
	file = strings.Split(file, "/")[len(strings.Split(file, "/"))-1]
	arh = append(arh, "[info]", file, line, " ")
	arh = append(arh, arg...)
	arh = append(arh, End)
	l.Print(arh...)

}
func (l *Logger) Error(arg ...interface{}) {
	var (
		file string
		line int
		arh  []interface{}
	)
	_, file, line, _ = runtime.Caller(1)
	file = strings.Split(file, "/")[len(strings.Split(file, "/"))-1]
	arh = append(arh, "[error]", file, line, " ")
	arh = append(arh, arg...)
	arh = append(arh, End)
	l.Print(arh...)
}

func (l *Logger) Recover(arg ...interface{}) {
	var (
		file string
		line int
		arh  []interface{}
	)
	_, file, line, _ = runtime.Caller(-1)
	arh = append(arh, "[error]", file, line, "\n")
	_, file, line, _ = runtime.Caller(0)
	arh = append(arh, "[error]", file, line, "\n")
	_, file, line, _ = runtime.Caller(1)
	arh = append(arh, "[error]", file, line, "\n")
	_, file, line, _ = runtime.Caller(2)
	arh = append(arh, "[error]", file, line, "\n")
	_, file, line, _ = runtime.Caller(3)
	arh = append(arh, "[error]", file, line, "\n")
	_, file, line, _ = runtime.Caller(4)
	arh = append(arh, "[error]", file, line, "\n")
	arh = append(arh, arg...)
	arh = append(arh, End)
	l.Print(arh...)
}
func (l *Logger) Warning(arg ...interface{}) {
	var (
		file string
		line int
		arh  []interface{}
	)
	_, file, line, _ = runtime.Caller(1)
	arh = append(arh, "[warning]", file, line)
	arh = append(arh, arg...)
	arh = append(arh, End)
	l.Print(arh...)
}

func LogFile(p string, path string) {
	prefix = p
	CurDay = string(time.Now().Format("2006-01-02"))
	allLog = path + "all." + CurDay + ".txt"
	errLog = path + "error." + CurDay + ".txt"
	if CheckFile(allLog) {
		fileAll, _ = os.OpenFile(allLog, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	} else {
		fileAll, _ = os.Create(allLog)
	}
	if CheckFile(errLog) {
		fileErr, _ = os.OpenFile(errLog, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	} else {
		fileErr, _ = os.Create(errLog)
	}
	LogErr.Logger = log.New(fileErr, prefix, log.Ltime)
	LogAll.Logger = log.New(fileAll, prefix, log.Ltime)
}

func Rotate(sleep time.Duration) {
	for {

		//CurDay = time.Now().Add(time.Hour * -24).Format("2006-01-02")
		if CurDay != time.Now().Format("2006-01-02") {
			fileAll.Close()
			os.Rename(allLog, allLog+"."+CurDay)
			fileAll, _ = os.OpenFile(allLog, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
			Log.Logger = log.New(fileAll, prefix, log.Ltime)
			CurDay = time.Now().Format("2006-01-02")
		}
		time.Sleep(sleep)

	}
}

func InitLog(p string, path string, name string) {
	prefix = p
	CurDay = time.Now().Format("2006-01-02")
	allLog = path + name + ".log"
	fileAll, _ = os.OpenFile(allLog, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	Log.Logger = log.New(fileAll, prefix, log.Ltime)

	if runtime.GOOS == "linux" {
		Env = true
		End = "\n"
	} else {
		End = "\r\n"
	}
}
func InitLogDebug(p string, path string, name string, level int) error {
	var (
		err error
	)
	Level = level
	prefix = p
	if level == 4 {
		return Errorf("level error")
	}
	CurDay = string(time.Now().Format("2006-01-02"))
	allLog = path + name + "." + CurDay + ".log"
	fileAll, err = os.OpenFile(allLog, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	Log.Logger = log.New(fileAll, prefix, log.Ltime|log.Lshortfile)
	if runtime.GOOS == "linux" {
		Env = true
		End = "\n"
	} else {
		End = "\r\n"
	}
	return nil
}

func Debug(arg ...interface{}) {
	var (
		file string
		line int
		arh  []interface{}
	)
	_, file, line, _ = runtime.Caller(1)
	file = strings.Split(file, "/")[len(strings.Split(file, "/"))-1]
	arh = append(arh, Now(), file, line, " ")
	arh = append(arh, arg...)
	fmt.Println(arh...)
}

func Stack(arg ...interface{}) {
	var (
		arh []interface{}
	)
	arh = append(arh, "[track]", string(debug.Stack()))
	arh = append(arh, arg...)
	arh = append(arh, End)
	fmt.Println(arh...)
}
func Track(arg ...interface{}) {

	var (
		file1, file2 string
		line1, line2 int
		arh          []interface{}
	)
	_, file1, line1, _ = runtime.Caller(1)
	_, file2, line2, _ = runtime.Caller(2)
	arh = append(arh, "[track]", file1, line1, " ", file2, line2, " ")
	arh = append(arh, arg...)
	arh = append(arh, End)
	fmt.Println(arh...)
}
