package goadmin_os

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
)

// 创建文件
func createFile(path string) (*os.File, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	return f, nil
}

// 创建文件
// 如果文件存在,不创建
// 不存在,则创建,"创建文件成功"
func AddFile(filePath string) (res string, err error) {
	exists, err := PathExists(filePath)
	if err != nil {
		return "", errors.New(fmt.Sprintf("exec PathExists fail:%s", err.Error()))
	}
	if exists {
		return "文件已存在", errors.New("文件已存在")
	} else {
		_, err = createFile(filePath)
		if err != nil {
			return "创建文件失败", err
		}
		return "创建文件成功", nil
	}
}

// 读文件
func ReadFile(path string) ([]byte, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// 复制文件
func CopyFile(oldFilePath, bakFilePath string) (res string, err error) {
	bakf, err := os.OpenFile(bakFilePath, os.O_CREATE|os.O_RDWR, 0777)
	defer bakf.Close()
	if err != nil {
		return "", errors.New(fmt.Sprintf("打开backFile文件失败:%s", err.Error()))
	}

	f, err := os.OpenFile(oldFilePath, os.O_CREATE|os.O_RDWR, 0777)
	defer f.Close()
	if err != nil {
		return "", errors.New(fmt.Sprintf("打开oldFile失败:%s", err.Error()))
	}

	_, err = io.Copy(bakf, f)
	if err != nil {
		return "", errors.New(fmt.Sprintf("备份失败:%s", err.Error()))
	}

	return "备份文件成功", nil
}

// 文件写入
func WriteFile(filePath string, content string) (res string, err error) {
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0777)
	defer f.Close()
	if err != nil {
		return "", errors.New(fmt.Sprintf("打开更新的文件失败:%s", err.Error()))

	}
	//fmt.Println(written)

	writer := bufio.NewWriter(f)
	_, err = writer.WriteString(content)

	if err != nil {
		return "", errors.New(fmt.Sprintf("无法写入文件:%s", err.Error()))
	}
	//f.Seek(0, 0)
	writer.Flush()

	return "文件写入成功", nil
}

// 更新文件内容
func UpdateRowFileWithContent(filePath string) (scanner *bufio.Scanner, writer *bufio.Writer, file *os.File, err error) {
	f, err := os.OpenFile(filePath, os.O_RDWR, 0777)
	if err != nil {
		return nil, nil, nil, errors.New(fmt.Sprintf("打开更改的文件失败:%s", err.Error()))
	}
	//defer f.Close()
	sca := bufio.NewScanner(f)
	buf := make([]byte, 64*1024)
	sca.Buffer(buf, bufio.MaxScanTokenSize)

	wri := bufio.NewWriterSize(f, bufio.MaxScanTokenSize)

	return sca, wri, f, nil

	/* 调用
		scanner, writer, file, err := goadmin_os.UpdateRowFileWithContent(tmpFile)
	defer file.Close()
	if err != nil {
		Failed(err.Error(), c)
		return
	}
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "GOLANG_KAFKA_IPS") {
			line = strings.Replace(line, "GOLANG_KAFKA_IPS", v.KafkaIPs, -1)
		}
		if strings.Contains(line, "GOLANG_ES_IPS") {
			line = strings.Replace(line, "GOLANG_ES_IPS", v.EsIPs, -1)
		}
		if strings.Contains(line, "GOLANG_Zone") {
			line = strings.Replace(line, "GOLANG_Zone", v.Zone, -1)
		}
		if strings.Contains(line, "GOLANG_ES_USERNAME") {
			line = strings.Replace(line, "GOLANG_ES_USERNAME", v.EsUserName, -1)
		}
		if strings.Contains(line, "GOLANG_ES_PASSWORD") {
			line = strings.Replace(line, "GOLANG_ES_PASSWORD", v.EsPassword, -1)
		}
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			Failed("写入失败"+err.Error(), c)
			return
		}
	}

	file.Seek(0, 0)
	writer.Flush()
	*/
	//file.Seek(0, 0)
	//writer.Flush()
	//
	//return "写入成功", nil
}
