package gofile

import (
	"encoding/json"
	"os"
)

func WriteJson(path string, d any) (err error) {
	dBytes, err := json.Marshal(d)
	if err != nil {
		return err
	}
	err = os.WriteFile(path, dBytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

func WriteIndentedJson(path string, d any) (err error) {
	dBytes, err := json.MarshalIndent(d, "", "\t")
	if err != nil {
		return err
	}
	err = os.WriteFile(path, dBytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

func Write(path string, d []byte) (err error) {
	return os.WriteFile(path, d, 0644)
}

func Append(path string, d []byte) (err error) {
	f, err := os.OpenFile(path,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.Write(d); err != nil {
		return err
	}

	return nil
}

func AppendString(path string, s string) (err error) {
	f, err := os.OpenFile(path,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.WriteString(s); err != nil {
		return err
	}
	return nil
}

func Rm(path string) (err error) {
	return os.Remove(path)
}

func RmAll(path string) (err error) {
	return os.RemoveAll(path)
}
