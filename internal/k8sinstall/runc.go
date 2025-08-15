package k8sinstall

import (
	"fmt"
	"k8s-offline-installer/internal/util"
)

// InstallRunc 安装runc
func InstallRunc() {
	util.Logger.Infof("开始安装runc...")
	pkg := "runc.amd64"
	downloadURL := "https://github.com/opencontainers/runc/releases/download/v1.3.0/" + pkg

	if err := util.DownloadFile(downloadURL, pkg); err != nil {
		util.Logger.Fatalf("下载runc失败: %v", err)
	}
	if err := util.ScpToRemote(pkg); err != nil {
		util.Logger.Fatalf("上传runc失败: %v", err)
	}
	defer util.CleanLocalFiles(pkg)

	// 远程配置
	err := util.ExecRemote(fmt.Sprintf(
		"cd %s && cp runc.amd64 /usr/bin/runc && chmod +x /usr/bin/runc",
		util.Cfg.TargetPath,
	))
	if err != nil {
		return
	}
	util.Logger.Infof("runc安装完成")
}
