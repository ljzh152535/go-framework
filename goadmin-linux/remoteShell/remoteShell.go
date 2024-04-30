package remoteShell

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"os"
	"path"
	"time"
)

const (
	SSHTypeP string = "password"
	SSHTypeK string = "key"
)

type Cli struct {
	IP       string
	Username string
	Password string
	Port     int
	SSHType  string // password或者key
	//client     *ssh.Client
	//SSHKeyPath string // ssh id_rsa.id路径
	LastResult string
}

//func New(ip string, username string, password string, SSHType string, port ...int) *Cli {
//	cli := new(Cli)
//	cli.IP = ip
//	cli.Username = username
//	cli.Password = password
//	cli.SSHType = SSHType
//	switch {
//	case len(port) <= 0:
//		cli.Port = 22
//	case len(port) > 0:
//		cli.Port = port[0]
//	}
//	return cli
//}

func formatError(err error, errStr string) string {
	if err != nil {
		return fmt.Sprintf("%s:%s", errStr, err.Error())
	}
	return ""
}

func publicKeyAuthFunc() (sshAuthMethod ssh.AuthMethod, errStr string) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return nil, formatError(err, "get UserHomeDir failed")
	}

	key, err := os.ReadFile(path.Join(homePath, ".ssh", "id_rsa"))
	if err != nil {
		return nil, formatError(err, "ssh key file read failed")
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, formatError(err, "ssh key signer failed")
	}
	return ssh.PublicKeys(signer), ""
}

func (c *Cli) ExecRemoteShell(cli *Cli, command string) string {
	//创建sshp登陆配置
	config := ssh.ClientConfig{
		User:            cli.Username,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // 这个可以,但是不够安全
		//HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		//	return nil
		//},
		Timeout: 1 * time.Second, //ssh 连接time out 时间一秒钟, 如果ssh验证错误 会在一秒内返回
	}

	if cli.SSHType == SSHTypeP {
		config.Auth = []ssh.AuthMethod{ssh.Password(cli.Password)}
	} else {
		publicKeyAuthFuncRes, str := publicKeyAuthFunc()
		if str == "" {
			config.Auth = []ssh.AuthMethod{publicKeyAuthFuncRes}
		} else {
			return str
		}
	}

	//dial 获取ssh client
	addr := fmt.Sprintf("%s:%d", cli.IP, cli.Port)
	sshClient, err := ssh.Dial("tcp", addr, &config)
	if err != nil {
		return formatError(err, "创建ssh client failed")
	}
	defer sshClient.Close()

	//创建ssh-session
	session, err := sshClient.NewSession()
	if err != nil {
		return formatError(err, "创建ssh session failed")
	}
	defer session.Close()

	stderr := &bytes.Buffer{} // make sure to import bytes
	//Stdout := &bytes.Buffer{} // make sure to import bytes
	session.Stderr = stderr

	// 执行远程命令
	//combo, err := session.Output("whoami; cd /; ls -al;echo https://github.com/libragen/felix")
	combo, err := session.Output(command)

	var errLog string
	for {
		line, errLine := stderr.ReadString('\n')
		if errLine != nil || io.EOF == errLine {
			break
		}
		errLog = line
	}
	if err != nil {
		if errLog != "" {
			return fmt.Sprintf("远程执行cmd failed:%s", errLog)
		}
		return fmt.Sprintf("远程执行cmd failed:%s", err.Error())
	}
	return fmt.Sprintf(string(combo))
}

func (c *Cli) ExecRemoteCommand(cli *Cli, command string) string {
	//创建sshp登陆配置
	config := ssh.ClientConfig{
		User:            cli.Username,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // 这个可以,但是不够安全
		//HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		//	return nil
		//},
		Timeout: 1 * time.Second, //ssh 连接time out 时间一秒钟, 如果ssh验证错误 会在一秒内返回
	}

	if cli.SSHType == SSHTypeP {
		config.Auth = []ssh.AuthMethod{ssh.Password(cli.Password)}
	} else {
		publicKeyAuthFuncRes, str := publicKeyAuthFunc()
		if str == "" {
			config.Auth = []ssh.AuthMethod{publicKeyAuthFuncRes}
		} else {
			return str
		}
	}

	//dial 获取ssh client
	addr := fmt.Sprintf("%s:%d", cli.IP, cli.Port)
	sshClient, err := ssh.Dial("tcp", addr, &config)
	if err != nil {
		return formatError(err, "创建ssh client failed")
	}
	defer sshClient.Close()

	//创建ssh-session
	session, err := sshClient.NewSession()
	if err != nil {
		return formatError(err, "创建ssh session failed")
	}
	defer session.Close()

	// 执行远程命令
	combo, err := session.CombinedOutput(command)

	if err != nil {
		return fmt.Sprintf("远程执行cmd failed:%s", string(combo))
	}
	return fmt.Sprintf("远程命令执行成功 success:%s", string(combo))
}
