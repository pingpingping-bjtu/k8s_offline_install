package k8sinstall

import (
	"fmt"
	"k8s-offline-installer/internal/util"
	"os/exec"
	"strings"
)

// InstallFlannel 安装flannel网络插件
func InstallFlannel() {
	util.Logger.Infof("开始安装flannel网络插件...")
	// 1. 安装CNI插件
	cniVersion := "v1.7.1"
	cniPkg := "cni-plugins-linux-amd64-" + cniVersion + ".tgz"
	cniURL := "https://github.com/containernetworking/plugins/releases/download/" + cniVersion + "/" + cniPkg

	if err := util.DownloadFile(cniURL, cniPkg); err != nil {
		util.Logger.Fatalf("下载CNI插件失败: %v", err)
	}
	if err := util.ScpToRemote(cniPkg); err != nil {
		util.Logger.Fatalf("上传CNI插件失败: %v", err)
	}
	defer util.CleanLocalFiles(cniPkg)
	err := util.ExecRemote(fmt.Sprintf("cd %s && mkdir -p /opt/cni/bin && tar -zxvf %s -C /opt/cni/bin", util.Cfg.TargetPath, cniPkg))
	if err != nil {
		return
	}

	// 2. 部署flannel
	flannelYaml := "kube-flannel.yml"
	flannelURL := "https://github.com/flannel-io/flannel/releases/latest/download/" + flannelYaml
	if err := util.DownloadFile(flannelURL, flannelYaml); err != nil {
		util.Logger.Fatalf("下载flannel配置失败: %v", err)
	}
	if err := util.ScpToRemote(flannelYaml); err != nil {
		util.Logger.Fatalf("上传flannel配置失败: %v", err)
	}
	defer util.CleanLocalFiles(flannelYaml)

	// 3. 拉取flannel镜像并部署
	// 先获取需要的镜像列表
	imageListCmd := exec.Command("sh", "-c", fmt.Sprintf("grep 'image:' %s | awk -F 'image: ' '{print $2}' | sort | uniq", flannelYaml))
	imageListOut, err := imageListCmd.CombinedOutput()
	if err != nil {
		util.Logger.Fatalf("获取flannel镜像列表失败: %v", err)
	}
	flannelImages := strings.Split(strings.TrimSpace(string(imageListOut)), "\n")
	fmt.Println(flannelImages[1])
	util.PullAndTransferImages(flannelImages)

	// 部署flannel
	err = util.ExecRemote(fmt.Sprintf("cd %s && kubectl apply -f %s", util.Cfg.TargetPath, flannelYaml))
	if err != nil {
		return
	}
	util.Logger.Infof("flannel安装完成")
}
