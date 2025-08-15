package k8sinstall

import (
	"fmt"
	"k8s-offline-installer/internal/util"
)

// InitK8sCluster 初始化K8s集群
func InitK8sCluster() {
	util.Logger.Infof("开始初始化K8s集群...")
	initScript := fmt.Sprintf(`
		kubeadm init \
			--apiserver-advertise-address=%s \
			--control-plane-endpoint=%s \
			--kubernetes-version %s \
			--skip-phases=preflight \
			--service-cidr=%s \
			--pod-network-cidr=%s \
			--cri-socket=%s \
			--ignore-preflight-errors=Swap \
			--v=5
	`, util.Cfg.TargetIP, util.Cfg.TargetIP, util.Cfg.K8sVersion, util.Cfg.ServiceCIDR, util.Cfg.PodCIDR, util.Cfg.CRISocket)
	err := util.ExecRemote(initScript)
	if err != nil {
		return
	}

	// 配置kubectl
	err = util.ExecRemote(`
			mkdir -p $HOME/.kube
			cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
			chown $(id -u):$(id -g) $HOME/.kube/config
		`)
	if err != nil {
		return
	}

	// 去除master污点（允许调度pod）
	err = util.ExecRemote("kubectl taint nodes --all node-role.kubernetes.io/control-plane-")
	if err != nil {
		return
	}
	util.Logger.Infof("K8s集群初始化完成")
}
