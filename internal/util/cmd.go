package util

import (
	"os"
	"os/exec"
)

// ExecLocal 执行本地命令（通用封装）
func ExecLocal(args ...string) error {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	Logger.Infof("执行本地命令: %v", args)
	if err := cmd.Run(); err != nil {
		Logger.Errorf("本地命令执行失败: %v", err)
		return err
	}
	return nil
}

// CleanLocalFiles 清理本地临时文件
func CleanLocalFiles(files ...string) {
	for _, f := range files {
		if err := os.Remove(f); err != nil {
			Logger.Errorf("清理文件 %s 失败: %v", f, err)
			continue
		}
		Logger.Infof("已清理本地文件: %s", f)
	}
}
