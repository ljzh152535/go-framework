package goadmin_os

import (
	"log"
	"math/rand"
	"os"
)

type FolderInfos struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	FolerType string `json:"folerType"`
	//ParentId  int
	IsParent bool          `json:"isParent"`
	Path     string        `json:"path"`
	Children []FolderInfos `json:"children"`
}

// 获取treetables列表工具
func TreeFolders(dstPath string) []FolderInfos {
	dirInfos, err := os.ReadDir(dstPath)
	PthSep := PathSeparator
	if err != nil {
		log.Fatal(err)
	}

	var infos []FolderInfos
	for _, dirInfo := range dirInfos {
		var folderInfos FolderInfos
		folderName := dirInfo.Name()
		tmp_dir := dstPath + PthSep + folderName

		folderInfos.Path = tmp_dir
		if dirInfo.IsDir() {
			//fmt.Println("dir:", folderName)
			folderInfos.ID = rand.Intn(20)
			folderInfos.Name = folderName
			folderInfos.FolerType = "d"
			folderInfos.IsParent = true
			//fmt.Println(folderInfos)
			tmpChildern := TreeFolders(tmp_dir)
			if len(tmpChildern) > 0 {
				folderInfos.Children = tmpChildern
			}
		} else {
			//ok := strings.HasSuffix(folderName, ".conf")
			//if ok {
			folderInfos.ID = rand.Intn(20)
			folderInfos.Name = dirInfo.Name()
			folderInfos.FolerType = "f"
			folderInfos.IsParent = false
			//}
		}
		infos = append(infos, folderInfos)
	}
	return infos
}
