package util

import (
	"fmt"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
)

// NewSSHClient 创建 SSH 客户端连接
func NewSSHClient() (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		User: Cfg.TargetUser,
		Auth: []ssh.AuthMethod{
			ssh.Password(Cfg.Password), // 也可替换为密钥认证
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // 生产环境建议验证主机密钥
		Timeout:         30 * time.Second,
	}
	addr := fmt.Sprintf("%s:%s", Cfg.TargetIP, Cfg.Port)
	return ssh.Dial("tcp", addr, config)
}

// ExecRemote 执行远程命令（通用封装，支持重试）
func ExecRemote(script string) error {
	client, err := NewSSHClient()
	if err != nil {
		return fmt.Errorf("SSH 连接失败: %v", err)
	}
	defer func(client *ssh.Client) {
		err := client.Close()
		if err != nil {

		}
	}(client)

	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("创建会话失败: %v", err)
	}
	defer func(session *ssh.Session) {
		err := session.Close()
		if err != nil {

		}
	}(session)

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	Logger.Infof("执行远程命令: ssh %v", script)

	// 重试机制（最多3次）
	for i := 0; i < Cfg.MaxRetries; i++ {

		if err := session.Run(script); err != nil {
			if i == Cfg.MaxRetries-1 {
				return fmt.Errorf("命令执行失败（重试%d次）: %v", Cfg.MaxRetries, err)
			}
			Logger.Errorf("远程命令执行失败（第%d次重试）: %v", i+1, err)
			time.Sleep(Cfg.RetryDelay)
			continue
		}
		return nil
	}
	return fmt.Errorf("远程命令多次执行失败: %s", script)

}
