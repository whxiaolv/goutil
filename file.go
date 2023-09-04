package goutil

import (
	"crypto/md5"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
)

type File struct {
}

func (this *File) CopyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}
	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}
	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()
	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func (this *File) MD5(str string) string {
	md := md5.New()
	md.Write([]byte(str))
	return fmt.Sprintf("%x", md.Sum(nil))
}

func (this *File) GetFileMd5(file *os.File) string {
	file.Seek(0, 0)
	md5h := md5.New()
	io.Copy(md5h, file)
	sum := fmt.Sprintf("%x", md5h.Sum(nil))
	return sum
}

func (this *File) Sha1(str string) string {
	md5h := sha1.New()
	md5h.Write([]byte(str))
	return fmt.Sprintf("%x", md5h.Sum(nil))
}

func (this *File) GetFileSha1Sum(file *os.File) string {
	file.Seek(0, 0)
	md5h := sha1.New()
	io.Copy(md5h, file)
	sum := fmt.Sprintf("%x", md5h.Sum(nil))
	return sum
}

func (this *File) GetFileSum(file *os.File, alg string) string {
	alg = strings.ToLower(alg)
	if alg == "sha1" {
		return this.GetFileSha1Sum(file)
	} else {
		return this.GetFileMd5(file)
	}
}
func (this *File) GetFileSumByName(filepath string, alg string) (string, error) {
	var (
		err  error
		file *os.File
	)
	file, err = os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	alg = strings.ToLower(alg)
	if alg == "sha1" {
		return this.GetFileSha1Sum(file), nil
	} else {
		return this.GetFileMd5(file), nil
	}
}

func (this *File) WriteFileByOffSet(filepath string, offset int64, data []byte) error {
	var (
		err   error
		file  *os.File
		count int
	)
	file, err = os.OpenFile(filepath, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	count, err = file.WriteAt(data, offset)
	if err != nil {
		return err
	}
	if count != len(data) {
		return errors.New(fmt.Sprintf("write %s error", filepath))
	}
	return nil
}
func (this *File) ReadFileByOffSet(filepath string, offset int64, length int) ([]byte, error) {
	var (
		err    error
		file   *os.File
		result []byte
		count  int
	)
	file, err = os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	result = make([]byte, length)
	count, err = file.ReadAt(result, offset)
	if err != nil {
		return nil, err
	}
	if count != length {
		return nil, errors.New("read error")
	}
	return result, nil
}

func (this *File) FileExists(fileName string) bool {
	_, err := os.Stat(fileName)
	return err == nil
}
func (this *File) WriteFile(path string, data string) bool {
	if err := ioutil.WriteFile(path, []byte(data), 0775); err == nil {
		return true
	} else {
		return false
	}
}
func (this *File) WriteBinFile(path string, data []byte) bool {
	if err := ioutil.WriteFile(path, data, 0775); err == nil {
		return true
	} else {
		return false
	}
}

func (this *File) ReadBinFile(path string) ([]byte, error) {
	if this.IsExist(path) {
		fi, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer fi.Close()
		return ioutil.ReadAll(fi)
	} else {
		return nil, errors.New("not found")
	}
}

func (this *File) IsExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func (this *File) RemoveEmptyDir(pathname string) {
	defer func() {
		if re := recover(); re != nil {
			buffer := debug.Stack()
			log.Print(string(buffer))
		}
	}()
	handlefunc := func(file_path string, f os.FileInfo, err error) error {
		if f.IsDir() {
			files, _ := ioutil.ReadDir(file_path)
			if len(files) == 0 && file_path != pathname {
				os.Remove(file_path)
			}
		}
		return nil
	}
	fi, _ := os.Stat(pathname)
	if fi.IsDir() {
		filepath.Walk(pathname, handlefunc)
	}
}
