package k8sinstall

import (
	"fmt"
	"k8s-offline-installer/internal/util"
)

// InstallCrictl 安装crictl工具
func InstallCrictl() {
	util.Logger.Infof("开始安装crictl...")
	pkg := "crictl-v1.32.0-linux-amd64.tar.gz"
	downloadURL := "https://github.com/kubernetes-sigs/cri-tools/releases/download/v1.32.0/" + pkg

	if err := util.DownloadFile(downloadURL, pkg); err != nil {
		util.Logger.Fatalf("下载crictl失败: %v", err)
	}
	if err := util.ScpToRemote(pkg); err != nil {
		util.Logger.Fatalf("上传crictl失败: %v", err)
	}
	defer util.CleanLocalFiles(pkg)

	// 远程解压并配置
	err := util.ExecRemote(fmt.Sprintf(
		"cd %s && tar -zxvf %s -C /usr/bin && "+
			"echo 'runtime-endpoint: %s' | tee /etc/crictl.yaml",
		util.Cfg.TargetPath, pkg, util.Cfg.CRISocket,
	))
	if err != nil {
		return
	}
	util.Logger.Infof("crictl安装完成")
}
