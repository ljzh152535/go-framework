package goadmin_os

import (
	"errors"
	"fmt"
	"os"
)

// 创建目录
func CreateDir(path string) error {
	// 创建单个目录
	err := os.Mkdir(path, os.ModePerm)
	if err != nil {
		return err
		//log.Fatal("创建文件失败")
	}
	//err := os.MkdirAll("test/a/b", os.ModePerm)
	//if err != nil {
	//	fmt.Printf("err: %v\n", err)
	//}
	return nil
}

// 删除目录
func RemoveDir(path string) error {
	/* err := os.Remove("test.txt")
	if err != nil {
		fmt.Printf("err: %v\n", err)
	} */

	err := os.RemoveAll(path)
	if err != nil {
		return err
		//log.Fatal("删除文件失败")
	}
	return nil
}

// 创建文件夹
// 如果文件夹存在,不创建
// 不存在,则创建,"创建目录成功"
func AddFolder(folderPath string) (res string, err error) {
	exists, err := PathExists(folderPath)
	if err != nil {
		return "", errors.New(fmt.Sprintf("exec PathExists fail:%s", err.Error()))
	}
	if exists {
		return "目录已存在", errors.New("目录已存在")
	} else {
		err = CreateDir(folderPath)
		if err != nil {
			return "创建目录失败", err
		}
		return "创建目录成功", nil
	}
}

type Folders struct {
	Name string `json:"name"`
	//Path string `json:"path"`
}

func ListDirs(path string, isDir bool) ([]Folders, error) {
	dirInfos, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	var folders []Folders

	for _, dirInfo := range dirInfos {
		var folder Folders
		folder.Name = dirInfo.Name()
		if dirInfo.IsDir() {
			if isDir {
				folders = append(folders, folder)
			}
			continue
		} else {
			if isDir {
				continue
			}
			folders = append(folders, folder)
		}

	}

	return folders, nil
}
