package gofile

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"

	"github.com/BurntSushi/toml"
)

func Bytes(path string) ([]byte, error) {
	fileData, err := os.ReadFile(path)
	if err != nil {
		return []byte{}, err
	}
	return fileData, nil
}

func String(path string) (string, error) {
	fileData, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(fileData), nil
}

func Csv(path string) ([][]string, error) {
	var res [][]string
	f, err := os.Open(path)
	if err != nil {
		return res, fmt.Errorf("error opening the file: %s", err.Error())
	}

	defer f.Close()

	csvReader := csv.NewReader(f)
	res, err = csvReader.ReadAll()
	if err != nil {
		return res, fmt.Errorf("error parsing the csv: %s", err.Error())
	}

	return res, nil
}

func SimpleCsv(path string) ([]string, error) {
	var res []string
	f, err := os.Open(path)
	if err != nil {
		return res, fmt.Errorf("error opening the file: %s", err.Error())
	}

	defer f.Close()

	csvReader := csv.NewReader(f)
	unparsedRes, err := csvReader.ReadAll()
	if err != nil {
		return res, fmt.Errorf("error parsing the csv: %s", err.Error())
	}
	for _, line := range unparsedRes {
		res = append(res, line[0])
	}

	return res, nil
}

func Toml(path string, pointer interface{}) error {
	tomlData, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("error laoding the toml file: %s", err.Error())
	}
	_, err = toml.Decode(string(tomlData), pointer)
	if err != nil {
		return fmt.Errorf("error decoding the toml file: %s", err.Error())
	}
	return nil
}

func Json(path string, pointer interface{}) error {
	jsonStruct := pointer
	jsonData, err := os.ReadFile(path)

	if err != nil {
		return fmt.Errorf("error loading the json file: %s", err.Error())
	}

	err = json.Unmarshal(jsonData, jsonStruct)
	if err != nil {
		return fmt.Errorf("error unmarshalling the json data: %s", err.Error())
	}

	return nil
}

// directory functions

type ReadDirOpts struct {
	suffix string
	filter ReadDirFilterType
}

type ReadDirOptFunc func(*ReadDirOpts)

type ReadDirFilterType string

const (
	DirFilter  ReadDirFilterType = "dir"
	FileFilter ReadDirFilterType = "file"
	NoneFilter ReadDirFilterType = "none"
)

func DefaultReadDirOpts() ReadDirOpts {
	return ReadDirOpts{
		suffix: "",
		filter: "none",
	}
}

func SuffixReadRirOpt(suffix string) ReadDirOptFunc {
	return func(rdo *ReadDirOpts) {
		rdo.suffix = suffix
	}
}

func FilterReadDirOpt(filter ReadDirFilterType) ReadDirOptFunc {
	return func(rdo *ReadDirOpts) {
		rdo.filter = filter
	}
}

func filterFilesSuffix(files []fs.DirEntry, suffix string) (resFiles []fs.DirEntry) {
	for _, file := range files {
		fn := file.Name()
		if len(fn) < len(suffix) {
			continue
		}
		if fn[len(fn)-len(suffix):] == suffix {
			resFiles = append(resFiles, file)
		}
	}
	return resFiles
}

func filterFilesType(files []fs.DirEntry, fileType ReadDirFilterType) (resFiles []fs.DirEntry) {
	if fileType == NoneFilter {
		return files
	}
	if fileType == FileFilter {
		for _, file := range files {
			if !file.IsDir() {
				resFiles = append(resFiles, file)
			}
		}
	} else if fileType == DirFilter {
		for _, file := range files {
			if file.IsDir() {
				resFiles = append(resFiles, file)
			}
		}
	}
	return resFiles
}

func ReadDir(path string, opts ...ReadDirOptFunc) (files []fs.DirEntry, err error) {
	rdOpts := DefaultReadDirOpts()
	for _, opt := range opts {
		opt(&rdOpts)
	}

	unfilteredFiles, err := os.ReadDir(path)
	if err != nil {
		return files, err
	}

	if len(rdOpts.suffix) != 0 {
		unfilteredFiles = filterFilesSuffix(unfilteredFiles, rdOpts.suffix)
	}
	files = filterFilesType(unfilteredFiles, rdOpts.filter)

	return files, nil
}

func ReadDirNames(path string, opts ...ReadDirOptFunc) (fileNames []string, err error) {
	rdOpts := DefaultReadDirOpts()
	for _, opt := range opts {
		opt(&rdOpts)
	}

	unfilteredFiles, err := os.ReadDir(path)
	if err != nil {
		return fileNames, err
	}

	var files []fs.DirEntry
	if len(rdOpts.suffix) != 0 {
		unfilteredFiles = filterFilesSuffix(unfilteredFiles, rdOpts.suffix)
	}
	files = filterFilesType(unfilteredFiles, rdOpts.filter)

	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}

	return fileNames, nil
}
