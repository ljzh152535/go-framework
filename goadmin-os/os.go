package goadmin_os

import (
	"errors"
	"fmt"
	"os"
)

// 获得工作目录
func GetWd() (s string, e error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return dir, nil
}

// 修改工作目录
func ChWd(path string) error {
	err := os.Chdir(path)
	if err != nil {
		return errors.New("修改工作目录失败")
	}
	return nil
}

// 获得临时目录
func GetTemp() {
	s := os.TempDir()
	fmt.Printf("s: %v\n", s)
}

// 重命名文件
func RenameFile(oldPath string, newPath string) error {
	//err := os.Rename("test.txt", "test2.txt")
	err := os.Rename(oldPath, newPath)
	if err != nil {
		return errors.New("重命名文件失败")
	}
	return nil
}

// 写文件
func writeFile() {
	s := "hello world"
	os.WriteFile("test2.txt", []byte(s), os.ModePerm)
}

// 判断文件是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path) // 获取目录信息,如果返回的 error 为空，那么说明文件是存在的，如果返回的错误信息是 os.IsNotExist 说明文件是不存在的。
	if err == nil {
		return true, nil // 文件存在
	}
	if os.IsNotExist(err) {
		return false, nil // 文件不存在
	}
	return false, err //返回其他错误信息
}
