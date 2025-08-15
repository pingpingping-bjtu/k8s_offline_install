package main

import (
	"k8s-offline-installer/internal/k8sinstall"
	"k8s-offline-installer/internal/util"
)

func main() {
	// 初始化（日志、配置）
	util.InitLog()
	util.LoadConfig("configs/config.yaml") // 从配置文件加载参数（替代硬编码）
	cmd := "mkdir -p /data"
	err := util.ExecRemote(cmd)
	if err != nil {
		return
	}

	// 执行安装流程（按顺序调用k8sinstall模块的函数）
	k8sinstall.ConfigKernelParams()       // 配置内核参数
	k8sinstall.InstallContainerd()        // 安装containerd
	k8sinstall.InstallCrictl()            // 安装crictl
	k8sinstall.InstallRunc()              // 安装runc
	k8sinstall.InstallK8sComponents()     // 安装k8s组件
	k8sinstall.PullAndTransferK8sImages() // 拉取并传输K8s镜像
	k8sinstall.InitK8sCluster()           // 初始化集群
	k8sinstall.InstallFlannel()           // 安装flannel

	util.Logger.Infof("K8s集群离线安装完成！")
}
