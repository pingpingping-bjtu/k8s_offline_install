package k8sinstall

import "k8s-offline-installer/internal/util"

// K8sImages K8s v1.32.2所需镜像列表（根据官方默认清单）
var K8sImages = []string{
	"registry.k8s.io/kube-apiserver:v1.32.2",
	"registry.k8s.io/kube-controller-manager:v1.32.2",
	"registry.k8s.io/kube-scheduler:v1.32.2",
	"registry.k8s.io/kube-proxy:v1.32.2",
	"registry.k8s.io/coredns/coredns:v1.11.3",
	"registry.k8s.io/pause:3.10",
	"registry.k8s.io/etcd:3.5.16-0",
}

// PullAndTransferK8sImages 拉取并传输K8s镜像
func PullAndTransferK8sImages() {
	util.PullAndTransferImages(K8sImages)
}
