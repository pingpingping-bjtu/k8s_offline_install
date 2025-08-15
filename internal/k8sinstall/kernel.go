package k8sinstall

import "k8s-offline-installer/internal/util"

// ConfigKernelParams 配置系统内核参数（关闭防火墙、SELinux、swap等）
func ConfigKernelParams() {
	util.Logger.Infof("开始配置内核参数...")
	// 原configKernelParams中的script逻辑（使用util.ExecRemote执行）
	script := `
systemctl stop firewalld
systemctl disable firewalld
sed -i 's/enforcing/disabled/' /etc/selinux/config
setenforce 0
swapoff -a
sed -ri 's/.*swap.*/#&/' /etc/fstab

cat <<EOF | sudo tee /etc/sysctl.d/k8s.conf > /dev/null
net.bridge.bridge-nf-call-ip6tables = 1
net.bridge.bridge-nf-call-iptables  = 1
net.ipv4.ip_forward                 = 1
EOF
modprobe br_netfilter
sysctl -p /etc/sysctl.d/k8s.conf
`
	if err := util.ExecRemote(script); err != nil {
		util.Logger.Fatalf("内核参数配置失败: %v", err)
	}
	util.Logger.Infof("内核参数配置完成")
}
