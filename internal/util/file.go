package util

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"golang.org/x/crypto/ssh"

	"github.com/pkg/sftp"
)

// DownloadFile 下载文件（通用封装，支持重试）
func DownloadFile(url, filename string) error {

	// 检查文件是否已存在
	if _, err := os.Stat(filename); err == nil {
		Logger.Infof("文件 %s 已存在，跳过下载", filename)
		return nil
	}

	Logger.Infof("开始下载: %s", url)
	for i := 0; i < Cfg.MaxRetries; i++ {
		if err := ExecLocal("wget", "--user-agent=\"Mozilla\"", "-c", url, "-O", filename); err != nil {
			Logger.Errorf("下载失败(第%d次重试): %v", i+1, err)
			continue
		}
		Logger.Infof("下载完成: %s", filename)
		return nil
	}
	return fmt.Errorf("文件 %s 下载失败", filename)
}

// ScpToRemote 上传文件到远程主机（通用封装）
func ScpToRemote(localFile string) error {
	client, err := NewSSHClient() // 复用上面的 SSH 客户端
	if err != nil {
		return fmt.Errorf("SSH 连接失败: %v", err)
	}
	defer func(client *ssh.Client) {
		err := client.Close()
		if err != nil {

		}
	}(client)

	// 创建 SFTP 客户端
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return fmt.Errorf("创建 SFTP 客户端失败: %v", err)
	}
	defer func(sftpClient *sftp.Client) {
		err := sftpClient.Close()
		if err != nil {

		}
	}(sftpClient)

	// 打开本地文件
	localF, err := os.Open(localFile)
	if err != nil {
		return fmt.Errorf("打开本地文件失败: %v", err)
	}
	defer func(localF *os.File) {
		err := localF.Close()
		if err != nil {

		}
	}(localF)

	// 获取远程文件路径（目标目录 + 文件名）
	fileName := filepath.Base(localFile)
	remotePath := filepath.Join(Cfg.TargetPath, fileName)

	// 创建远程文件
	remoteF, err := sftpClient.Create(remotePath)
	if err != nil {
		return fmt.Errorf("创建远程文件失败: %v", err)
	}
	defer func(remoteF *sftp.File) {
		err := remoteF.Close()
		if err != nil {

		}
	}(remoteF)

	// 传输文件内容
	Logger.Infof("上传文件: %s -> %s", localFile, remotePath)
	if _, err := io.Copy(remoteF, localF); err != nil {
		return fmt.Errorf("文件传输失败: %v", err)
	}

	// 保持文件权限（可选）
	localInfof, _ := localF.Stat()
	err = sftpClient.Chmod(remotePath, localInfof.Mode())
	if err != nil {
		return err
	}

	return nil
}
