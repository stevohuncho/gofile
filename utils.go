package gofile

import "os"

func Exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func SafeDir(dir string) error {
	if Exists(dir) {
		return nil
	}
	return os.Mkdir(dir, os.ModePerm)
}
