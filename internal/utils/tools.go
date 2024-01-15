package utils

import (
	"os"
	"path"
	"time"
)

var (
	err        error
	timeFormat = time.Now().Format("2006_01_02_15_04_05")
	fn         = "./data/" + timeFormat + ".txt"
)

func DownFile(s string) error {
	err = os.MkdirAll(path.Dir(fn), os.ModePerm)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	_, err = file.Write([]byte(s))
	if err != nil {
		return err
	}
	return nil
}
