package ctl

import (
	"bufio"
	"errors"
	"io"
	"io/ioutil"

	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gogf/gf/os/gfile"
	"github.com/pkg/sftp"
)

func UnixDir(path string) string {
	path = filepath.Dir(path)
	path = strings.ReplaceAll(path, "\\", "/")
	if path[len(path)-1] != '/' {
		path += "/"
	}
	return path
}

func UnixSubDir(path string) string {
	if path[len(path)-1] != '/' {
		path += "/"
	}
	return strings.Split(path, "/")[len(strings.Split(path, "/"))-2]
}

func UnixPrtDir(path string) string {
	path = filepath.Dir(path)
	path = UnixDir(path)
	return path
}
func UnixPath(path string) string {
	path = strings.ReplaceAll(path, "\\", "/")
	if gfile.IsDir(path) {
		if path[len(path)-1] != '/' {
			path += "/"
		}
	}
	return path
}
func CheckFile(path string) bool {
	var exist = true
	if _, err := os.Stat(path); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func CheckFilePath(path string) bool {
	var exist = true
	if info, err := os.Stat(path); os.IsNotExist(err) || info.IsDir() {
		exist = false
	}
	return exist
}

func CheckDir(path string) bool {
	var exist = true
	if info, err := os.Stat(path); os.IsNotExist(err) || !info.IsDir() {
		exist = false
	}
	return exist
}
func Join(p ...string) string {
	return filepath.Join(p...)
}
func GetFileList(path string) (int, []string, error) {
	var (
		fileList  []string
		pathCount int
	)

	path = filepath.ToSlash(path)
	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		fileList = append(fileList, path)
		return nil
	})
	if err != nil {
		return 0, nil, err
	}
	list := strings.Split(path, "/")
	if list[len(list)-1] == "" {
		pathCount = len(list) - 1
	} else {
		pathCount = len(list)
	}
	maxCount := 0
	for _, v := range fileList {
		v = filepath.ToSlash(v)
		listFile := strings.Split(v, "/")
		if len(listFile)-pathCount > maxCount {
			maxCount = len(listFile) - pathCount
		}
	}
	return maxCount, fileList, err
}

func ListDir(path string) ([]string, error) {
	var fileList []string
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, v := range files {
		fileList = append(fileList, filepath.Join(path, v.Name()))
	}

	return fileList, nil
}

func ListDirAll(path string) ([]os.FileInfo, error) {
	var fileList []os.FileInfo
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, v := range files {
		if v.IsDir() {
			fileList = append(fileList, v)
			fileTmp, _ := ListDirAll(filepath.Join(path, v.Name()))

			for _, t := range fileTmp {
				fileList = append(fileList, t)
			}
		} else {
			fileList = append(fileList, v)
		}
	}
	return fileList, nil
}

func ListPath(path string) ([]string, error) {
	var (
		fileList, fileTmp []string
		files             []os.FileInfo
		err               error
	)
	files, err = ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, v := range files {
		if v.IsDir() {
			fileList = append(fileList, filepath.Join(path, v.Name()))
			fileTmp, _ = ListPath(filepath.Join(path, v.Name()))

			for _, t := range fileTmp {

				fileList = append(fileList, t)
			}
		} else {
			fileList = append(fileList, filepath.Join(path, v.Name()))
		}
	}
	return fileList, nil
}
func ListSftpPath(path string, client *sftp.Client) ([]string, error) {
	var (
		fileList, fileTmp []string
		files             []os.FileInfo
		err               error
	)
	files, err = client.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, v := range files {
		if v.IsDir() {
			fileList = append(fileList, strings.Replace(filepath.Join(path, v.Name()), "\\", "/", -1))
			fileTmp, _ = ListSftpPath(filepath.Join(path, v.Name()), client)

			for _, t := range fileTmp {
				fileList = append(fileList, strings.Replace(t, "\\", "/", -1))
			}
		} else {
			//fileList = append(fileList, filepath.Join(path, v.Name()))
			fileList = append(fileList, strings.Replace(filepath.Join(path, v.Name()), "\\", "/", -1))
		}
	}
	return fileList, nil
}

func IsDir(path string) bool {
	var (
		err  error
		info os.FileInfo
	)
	info, err = os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func IsFile(path string) bool {
	var (
		err  error
		info os.FileInfo
	)
	info, err = os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()

}

func IsZipFile(path string) bool {
	regZip, _ := regexp.Compile(".*\\.zip$")
	if regZip.MatchString(path) {
		return true
	} else {
		return false
	}
}
func PathUinx(path string) string {

	path = strings.Replace(path, "\\", "/", -1)
	if path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}
	return path
}

func ReadLine(path string, fc func(line string) error) error {
	var (
		err  error
		buf  *bufio.Reader
		f    *os.File
		line string
	)
	f, err = os.Open(path)
	if err != nil {
		return err
	}
	buf = bufio.NewReader(f)
	defer f.Close()
	for {
		line, err = buf.ReadString(10)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		err = fc(line)
		if err != nil {
			return err
		}
	}
	f.Close()
	return nil
}

func LastLine(p string) (string, error) {
	var (
		f   *os.File
		err error
		l   string
		buf []byte
		fl  []string
		li  int
	)
	f, err = os.Open(p)
	if err != nil {
		return l, err
	}
	buf = make([]byte, 512)
	f.Seek(-512, 2)
	li, err = f.Read(buf)
	if err != nil {
		return l, err
	}
	buf = buf[:li]
	fl = strings.Split(string(buf), "\n")
	if len(fl) < 2 {
		buf = make([]byte, 2048)
		f.Seek(-2048, 2)
		li, err = f.Read(buf)
		if err != nil {
			return l, err
		}
		buf = buf[:li]
		fl = strings.Split(string(buf), "\n")
	}
	f.Close()
	if fl[len(fl)-1] == "" {
		return fl[len(fl)-2], err
	}
	return fl[len(fl)-1], err

}

func ReadLineSlice(path string) ([]string, error) {
	var (
		err   error
		buf   *bufio.Reader
		f     *os.File
		line  string
		lines []string
	)
	f, err = os.Open(path)
	if err != nil {
		return lines, err
	}
	buf = bufio.NewReader(f)
	defer f.Close()
	for {
		line, err = buf.ReadString(10)
		if err != nil {
			if err == io.EOF {
				err = nil
				line = strings.Replace(line, "\r\n", "", -1)
				line = strings.Replace(line, "\n", "", -1)
				if line != "" {
					lines = append(lines, line)
				}
				break
			}
			return lines, err
		}
		line = strings.Replace(line, "\r\n", "", -1)
		line = strings.Replace(line, "\n", "", -1)
		lines = append(lines, line)
	}
	f.Close()
	return lines, err
}

func Append(path string, data interface{}) error {
	var (
		err  error
		file *os.File
	)
	file, err = os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	switch value := data.(type) {
	case string:
		_, err = file.WriteString(value)
		if err != nil {

			return err
		}
	case []byte:
		_, err = file.Write(value)
		if err != nil {
			return err
		}
	}
	return err
}

func ReadBytes(path string, fc func(line []byte) error) error {
	var (
		err  error
		buf  *bufio.Reader
		f    *os.File
		line []byte
	)
	f, err = os.Open(path)
	if err != nil {
		return err
	}
	buf = bufio.NewReader(f)
	defer f.Close()
	for {
		line, err = buf.ReadBytes(10)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		err = fc(line)
		if err != nil {
			return err
		}
	}
	f.Close()
	return nil
}

func Rename(src, dst string) error {
	var (
		err           error
		info, newInfo os.FileInfo
		newFile, file *os.File
	)
	err = os.Rename(src, dst)
	if err != nil && (strings.Contains(err.Error(), "invalid cross-device link") || strings.Contains(err.Error(), "The system cannot move the file to a different disk drive")) {
		info, err = os.Stat(src)
		if err != nil {
			return err
		}
		if info.IsDir() {
			src = DirStand(src)
			dst = DirStand(dst)
			err = filepath.Walk(src, func(name string, fi os.FileInfo, err error) error {
				name = strings.Replace(name, "\\", "/", -1)
				if fi.IsDir() {
					name = strings.Replace(name, src, dst, 1)
					err = os.Mkdir(name, fi.Mode())
					if err != nil {
						return err
					}
				} else {
					file, err = os.Open(name)
					name = strings.Replace(name, src, dst, 1)
					if err != nil {
						return err
					}
					defer file.Close()
					newFile, err = os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, fi.Mode())
					if err != nil {
						return err
					}
					defer newFile.Close()
					_, err = io.Copy(newFile, file)
					if err != nil {
						return err
					}
					file.Close()
					newFile.Close()
					newInfo, err = os.Stat(name)
					if err != nil {
						return err
					}
					if fi.Size() != newInfo.Size() {
						return io.ErrShortWrite
					}
				}
				return nil
			})
		} else {
			file, err = os.Open(src)
			if err != nil {
				return err
			}
			defer file.Close()
			newFile, err = os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, info.Mode())
			if err != nil {
				return err
			}
			defer newFile.Close()
			_, err = io.Copy(newFile, file)
			if err != nil {
				return err
			}
			file.Close()
			newFile.Close()
			newInfo, err = os.Stat(dst)
			if err != nil {
				return err
			}
			if info.Size() != newInfo.Size() {
				return io.ErrShortWrite
			}
		}
		err = os.RemoveAll(src)
	}
	return err
}

func CopyFile(src, dst string) error {
	var (
		err           error
		info, newInfo os.FileInfo
		newFile, file *os.File
	)
	if CheckFile(dst) {
		return errors.New("file or directory exists")
	}
	info, err = os.Stat(src)
	if err != nil {
		return err
	}
	if info.IsDir() {
		src = DirStand(src)
		dst = DirStand(dst)
		err = filepath.Walk(src, func(name string, fi os.FileInfo, err error) error {
			name = strings.Replace(name, "\\", "/", -1)
			if fi.IsDir() {
				name = strings.Replace(name, src, dst, 1)
				err = os.Mkdir(name, fi.Mode())
				if err != nil {
					return err
				}
			} else {
				file, err = os.Open(name)
				name = strings.Replace(name, src, dst, 1)
				if err != nil {
					return err
				}
				defer file.Close()
				newFile, err = os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, fi.Mode())
				if err != nil {
					return err
				}
				defer newFile.Close()
				_, err = io.Copy(newFile, file)
				if err != nil {
					return err
				}
				file.Close()
				newFile.Close()
				newInfo, err = os.Stat(name)
				if err != nil {
					return err
				}
				if fi.Size() != newInfo.Size() {
					return io.ErrShortWrite
				}
			}
			return nil
		})
	} else {
		file, err = os.Open(src)
		if err != nil {
			return err
		}
		defer file.Close()
		newFile, err = os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, info.Mode())
		if err != nil {
			return err
		}
		defer newFile.Close()
		_, err = io.Copy(newFile, file)
		if err != nil {
			return err
		}
		file.Close()
		newFile.Close()
		newInfo, err = os.Stat(dst)
		if err != nil {
			return err
		}
		if info.Size() != newInfo.Size() {
			return io.ErrShortWrite
		}
	}
	return err
}

func ReadFile(path string) ([]byte, error) {
	var (
		err error
		f   *os.File
		fd  []byte
	)
	f, err = os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	fd, err = ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	f.Close()
	return fd, nil
}

func DirStand(p string) string {
	p = strings.Replace(p, "\\", "/", -1)
	if []byte(p)[len(p)-1] == 47 {
		return string([]byte(p)[:len(p)-1])
	}
	return p
}

func CurDir() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return strings.Replace(dir, "\\", "/", -1)
}
