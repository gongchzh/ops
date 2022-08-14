package ctl

import (
	"archive/zip"
	"errors"
	"io"
	"os"
	ph "path"
)

func ZipCmpr(paths []string, zipFile string) error {
	var (
		f     *os.File
		err   error
		files []*os.File
	)
	for _, v := range paths {
		f, err := os.Open(v)
		if err != nil {
			return err
		}
		files = append(files, f)
	}
	err = zipCompress(files, zipFile)
	if err != nil {
		return err
	}
	defer f.Close()
	return nil

}

func zipCompress(files []*os.File, dest string) error {
	var (
		d   *os.File
		err error
		w   *zip.Writer
	)
	d, err = os.Create(dest)
	if err != nil {
		return err
	}
	defer d.Close()

	w = zip.NewWriter(d)
	defer w.Close()
	for _, file := range files {
		err = zipCompressFile(file, "", w)
		if err != nil {
			return err
		}
	}
	return nil
}

func zipCompressFile(file *os.File, prefix string, zw *zip.Writer) error {
	var (
		info      os.FileInfo
		err       error
		fileInfos []os.FileInfo
		f         *os.File
	)
	info, err = file.Stat()
	if err != nil {
		return err
	}
	if info.IsDir() {
		if prefix != "" {
			prefix = prefix + info.Name() + "/"
		} else {
			prefix = info.Name() + "/"
		}
		fileInfos, err = file.Readdir(-1)
		if err != nil {
			return err
		}
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		if prefix != "" {
			header.Name = prefix
		}
		_, err = zw.CreateHeader(header)
		if err != nil {
			return err
		}
		file.Close()
		for _, fi := range fileInfos {
			f, err = os.Open(file.Name() + "/" + fi.Name())
			if err != nil {
				return err
			}
			err = zipCompressFile(f, prefix, zw)
			if err != nil {
				return err
			}
		}
	} else {
		header, err := zip.FileInfoHeader(info)
		if prefix == "" {
			header.Name = header.Name
		} else {
			header.Name = prefix + header.Name
		}
		if err != nil {
			return err
		}
		header.Method = 8
		writer, err := zw.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, file)
		if err != nil {
			return err
		}
		file.Close()
	}
	file.Close()
	return nil
}

func ZipDecmpr(src, dst string) error {
	var (
		dstName string
		err     error
		srcFile io.ReadCloser
		file    *zip.ReadCloser
		newFile *os.File
		line    int64
	)
	file, err = zip.OpenReader(src)
	if err != nil {
		return err
	}
	if dst == "" {
		dst = ph.Dir(src)
	} else {
		dst = DirStand(dst)
	}
	defer file.Close()
	for _, v := range file.File {
		dstName = dst + "/" + v.Name
		if false {
			Debug(dstName)
			if CheckFile(dstName) {
				err = errors.New("file or directory exists")
				return err
			}
		}
		if v.FileInfo().IsDir() {
			err = os.MkdirAll(dstName, v.FileInfo().Mode())
			if err != nil {
				return err
			}
		}
	}
	for _, v := range file.File {
		dstName = dst + "/" + v.Name
		if v.FileInfo().IsDir() {
			continue
		}
		if CheckFile(dstName) {
			err = errors.New("file or directory exists")
			return err
		}
		srcFile, err = v.Open()
		if err != nil {
			return err
		}
		defer srcFile.Close()
		newFile, err = os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, v.FileInfo().Mode())
		if err != nil {
			return err
		}
		defer newFile.Close()
		line, err = io.Copy(newFile, srcFile)
		if err != nil {
			return err
		}
		if line != v.FileInfo().Size() {
			return errors.New("copy file size error")
		}
		srcFile.Close()
		newFile.Close()
	}
	file.Close()
	return nil
}
