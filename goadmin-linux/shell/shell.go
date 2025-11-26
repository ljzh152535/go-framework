package shell

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
)

// 执行shell脚本
func ExecLocalShell(command string) string {
	cmd := exec.Command("sh", "-c", command)
	//cmd := exec.Command(command)

	stderr := &bytes.Buffer{} // make sure to import bytes
	//Stdout := &bytes.Buffer{} // make sure to import bytes
	cmd.Stderr = stderr
	//cmd.Stdout = Stdout
	//resOut := Stdout.String()

	output, err := cmd.Output()

	var errLog string
	for {
		line, errLine := stderr.ReadString('\n')
		if errLine != nil || io.EOF == errLine {
			break
		}
		errLog = line
	}

	if err != nil {
		return fmt.Sprintf("Execute Shell failed with error:%s", errLog)
	}

	return string(output)
}

// 在本地执行linxu命令
func ExecLocalCommand(command string) string {
	cmd := exec.Command("sh", "-c", command)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Sprintf("Execute failed with error:%s", err.Error())
	}

	return string(output)
}
