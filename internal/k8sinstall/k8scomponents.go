package k8sinstall

import (
	"fmt"
	"k8s-offline-installer/internal/util"
	"strings"
)

// InstallK8sComponents 安装k8s组件（kubeadm/kubelet/kubectl）
func InstallK8sComponents() {
	util.Logger.Infof("开始安装K8s组件...")
	repoURL := "https://mirrors.tuna.tsinghua.edu.cn/kubernetes/core%3A/stable%3A/v1.32/rpm/x86_64/"
	packages := []string{
		"cri-tools-1.32.0-150500.1.1.x86_64.rpm",
		"kubernetes-cni-1.6.0-150500.1.1.x86_64.rpm",
		"kubeadm-" + strings.TrimPrefix(util.Cfg.K8sVersion, "v") + "-150500.1.1.x86_64.rpm",
		"kubectl-" + strings.TrimPrefix(util.Cfg.K8sVersion, "v") + "-150500.1.1.x86_64.rpm",
		"kubelet-" + strings.TrimPrefix(util.Cfg.K8sVersion, "v") + "-150500.1.1.x86_64.rpm",
	}

	// 下载并上传rpm包
	for _, pkg := range packages {
		if err := util.DownloadFile(repoURL+pkg, pkg); err != nil {
			util.Logger.Fatalf("下载%s失败: %v", pkg, err)
		}
		if err := util.ScpToRemote(pkg); err != nil {
			util.Logger.Fatalf("上传%s失败: %v", pkg, err)
		}
		util.CleanLocalFiles(pkg) // 清理本地文件
	}

	// 远程安装
	err := util.ExecRemote(fmt.Sprintf(
		"cd %s && yum -y localinstall %s",
		util.Cfg.TargetPath, strings.Join(packages, " "),
	))
	if err != nil {
		return
	}
	err = util.ExecRemote("systemctl enable kubelet")
	if err != nil {
		return
	}
	util.Logger.Infof("K8s组件安装完成")
}
