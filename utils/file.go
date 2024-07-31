package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

func Open(name string) (f *os.File, err error) {
	_, err = os.Stat(name)

	if os.IsNotExist(err) {
		return os.OpenFile(name, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	}
	if err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("file %s already exists", name)
}

func OpenFile(file string) ([]byte, error) {
	_file, err := os.Open(file)

	if err != nil {
		return nil, err
	}
	defer _file.Close()

	return io.ReadAll(_file)
}

func PathExists(path string) (bool, error) {
	s, err := os.Stat(path)

	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) || !s.IsDir() {
		return false, nil
	}
	return false, err
}

func FileExists(filepath string) (bool, error) {
	s, err := os.Stat(filepath)

	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) || s.IsDir() {
		return false, nil
	}
	return false, err
}

func IsDir(path string) bool {
	s, err := os.Stat(path)

	if err != nil {
		return false
	}
	return s.IsDir()
}

func Mkdir(path string) error {
	return os.Mkdir(path, os.ModePerm)
}

func MkdirAll(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

func Remove() {

}

func RemoveAll(pwd string, patterns ...string) error {
	infos, err := os.ReadDir(pwd)

	if err != nil {
		return err
	}
	for _, info := range infos {
		if info.IsDir() {
			continue
		}
		for _, pattern := range patterns {
			regex, _ := regexp.Compile(pattern)
			isExist := regex.MatchString(info.Name())

			if isExist {
				_ = os.Remove(path.Join(pwd, info.Name()))
			}
		}
	}
	return nil
}

func Create(filepath string) (*os.File, error) {
	_path := path.Dir(filepath)

	isExist, err := PathExists(_path)

	if err != nil {
		return nil, err
	} else if !isExist {
		if err = MkdirAll(_path); err != nil {
			return nil, err
		}
	}
	return os.Create(filepath)
}

func Touch(filepath string, content []byte) error {
	_path := path.Dir(filepath)

	isExist, err := PathExists(_path)

	if err != nil {
		return err
	} else if !isExist {
		if err = MkdirAll(_path); err != nil {
			return err
		}
	}
	var file *os.File

	if file, err = os.Create(filepath); err != nil {
		return err
	}
	defer file.Close()

	if content != nil {
		_, _ = file.Write(content)
	}
	return nil
}

func TouchJson(filepath string, content []byte) error {
	_path := path.Dir(filepath)

	isExist, err := PathExists(_path)

	if err != nil {
		return err
	} else if !isExist {
		if err = MkdirAll(_path); err != nil {
			return err
		}
	}
	var file *os.File

	if file, err = os.Create(filepath); err != nil {
		return err
	}
	defer file.Close()

	if content != nil {
		var out bytes.Buffer

		if err = json.Indent(&out, content, "", "\t"); err != nil {
			return err
		}
		file.Write([]byte(out.String()))
	}
	return nil
}

func PrepareOutput(path string) error {
	fi, err := os.Stat(path)

	if err != nil {
		if os.IsExist(err) && !fi.IsDir() {
			return err
		}
		if os.IsNotExist(err) {
			err = os.MkdirAll(path, 0777)
		}
	}
	return err
}

type Catalogue struct {
	Path string `json:"path"`
	File string `json:"file"`
}

// GetCatalogueFiles 获取目录下所有文件
func GetCatalogueFiles(catalogue string) []*Catalogue {
	pwd := ""

	if catalogue != "" {
		pwd = catalogue
	} else {
		pwd, _ = os.Getwd()
	}
	// 目录，文件名
	out := make([]*Catalogue, 0)

	if pwd == "" {
		return out
	}
	_ = filepath.Walk(pwd, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info != nil && !info.IsDir() {
			out = append(out, &Catalogue{Path: strings.ReplaceAll(path, "\\", "/"), File: info.Name()})
		}
		return nil
	})
	return out
}

// GetCurrentCatalogueFiles 获取当前目录下所有文件
func GetCurrentCatalogueFiles(catalogue string) []*Catalogue {
	pwd := ""

	if catalogue != "" {
		pwd = catalogue
	} else {
		pwd, _ = os.Getwd()
	}
	out := make([]*Catalogue, 0)

	infos, err := os.ReadDir(pwd)

	if err != nil {
		return out
	}
	for _, info := range infos {
		if info.IsDir() {
			continue
		}
		out = append(out, &Catalogue{
			File: info.Name(),
		})
	}
	return out
}

func Rename(OldPath string, NewPath string) error {
	return os.Rename(OldPath, NewPath)
}

type File struct {
	Path   string // 文件目录
	Name   string // 文件名
	Suffix string // 文件后缀
}

func FileInfo(file string) File {
	dir := path.Dir(file)
	base := path.Base(file)
	suffix := path.Ext(file)
	name := base[0 : len(base)-len(suffix)]
	return File{
		Path: dir, Name: name, Suffix: suffix,
	}
}

// FileCopy 文件复制
func FileCopy(distFile, newFile string) error {
	source, err := os.Open(distFile)

	if err != nil {
		return err
	}
	defer source.Close()

	var current *os.File

	if current, err = os.Create(newFile); err != nil {
		return err
	}
	defer current.Close()

	_, err = io.Copy(current, source)

	return nil
}
