package ctl

import (
	"bytes"
	"io/ioutil"
	"os"
	"regexp"
)

func ReplaceAll(path, oldStr, newStr string) error {
	var (
		err           error
		file, newByte []byte
		regStr        *regexp.Regexp
	)
	regStr, err = regexp.Compile(oldStr)
	if err != nil {
		return err
	}
	file, err = ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	newByte = regStr.ReplaceAll(file, []byte(newStr))
	err = ioutil.WriteFile(path, newByte, os.ModePerm)
	return err
}

func Replace(path, oldStr, newStr, reg string, num int) error {
	var (
		err             error
		file            *os.File
		regStr, regLine *regexp.Regexp
		i               int
		newFile         bytes.Buffer
	)
	regStr, err = regexp.Compile(oldStr)
	if err != nil {
		return err
	}
	regLine, err = regexp.Compile(reg)
	if err != nil {
		return err
	}
	if num < 0 {
		num = 1000000000
	}
	err = ReadBytes(path, func(line []byte) error {
		if i < num {
			if regLine.Match(line) {
				if regStr.Match(line) {
					newFile.Write(regStr.ReplaceAll(line, []byte(newStr)))
					i += 1
				} else {
					newFile.Write(line)
				}
			} else {
				newFile.Write(line)
			}
		} else {
			newFile.Write(line)
		}
		return nil
	})
	if err != nil {
		return err
	}
	file, err = os.OpenFile(path, os.O_RDWR|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	_, err = newFile.WriteTo(file)
	file.Close()
	return err
}
