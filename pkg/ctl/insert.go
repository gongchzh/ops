package ctl

import (
	"bytes"
	"os"
	"regexp"
)

func Add(path, line string) error {
	var (
		err  error
		file *os.File
	)
	file, err = os.OpenFile(path, os.O_APPEND, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(line)
	file.WriteString("\n")
	file.Close()
	return err
}

func InsertLast(path, reg, newLine string, num int) error {
	var (
		err     error
		file    *os.File
		regStr  *regexp.Regexp
		newFile bytes.Buffer
		i       int
	)
	regStr, err = regexp.Compile(reg)
	if err != nil {
		return err
	}
	if num < 0 {
		num = 2000000000
	}
	err = ReadBytes(path, func(line []byte) error {
		if i < num {
			if regStr.Match(line) {
				newFile.WriteString(newLine)
				newFile.WriteByte(10)
				i += 1
			}
		}
		newFile.Write(line)
		return nil
	})
	if err != nil {
		return err
	}
	file, err = os.OpenFile(path, os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	_, err = newFile.WriteTo(file)
	file.Close()
	return err
}

func InsertNext(path, reg, newLine string, num int) error {
	var (
		err     error
		file    *os.File
		regStr  *regexp.Regexp
		newFile bytes.Buffer
		i       int
	)
	regStr, err = regexp.Compile(reg)
	if err != nil {
		return err
	}
	if num < 0 {
		num = 2000000000
	}
	err = ReadBytes(path, func(line []byte) error {
		newFile.Write(line)
		if i < num {
			if regStr.Match(line) {
				newFile.WriteString(newLine)
				newFile.WriteByte(10)
				i += 1
			}
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
