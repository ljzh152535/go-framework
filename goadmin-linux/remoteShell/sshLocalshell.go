package remoteShell

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

// 登录远程机器执行本地shell脚本
func LoginRemoteExecLoaclShell(username, ip, localScriptPathWithParmas string) (shellCommand, resOut string) {
	command := fmt.Sprintf("ssh -o ConnectTimeout=2 %s@%s -C \"/bin/bash -s\" < %s", username, ip, localScriptPathWithParmas)
	cmd := exec.Command("sh", "-c", command)
	//fmt.Println("conmmand:", command)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	var errLog string
	for {
		line, errLine := stderr.ReadString('\n')
		if errLine != nil || io.EOF == errLine {
			break
		}
		errLog = line
	}
	if err != nil {
		return command, fmt.Sprintf("Execute Shell failed with error:%s",
			strings.Replace(errLog, "\n", "", -1))
	}

	//reg := regexp.MustCompile(`( )+|(\n)+`)
	//after := reg.ReplaceAllString(string(out.String()), "$1$2")
	// 去除返回结果的换行符
	return command, strings.Replace(string(out.String()), "\n", "", -1)
	//return command, string(out.String())
}
