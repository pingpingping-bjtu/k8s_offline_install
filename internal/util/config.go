package util

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config 配置结构体（与configs/config.yaml对应）
type Config struct {
	TargetIP      string        `yaml:"target_ip,omitempty"`      // 目标主机IP
	TargetUser    string        `yaml:"target_user,omitempty"`    // 目标主机用户
	Port          string        `yaml:"port,omitempty"`           // SSH端口
	Password      string        `yaml:"password,omitempty"`       // SSH密码
	TargetPath    string        `yaml:"target_path,omitempty"`    // 目标主机临时目录
	K8sVersion    string        `yaml:"k8s_version,omitempty"`    // K8s版本
	ContainerdVer string        `yaml:"containerd_ver,omitempty"` // Containerd版本
	CRISocket     string        `yaml:"cri_socket,omitempty"`     // CRI socket路径
	ServiceCIDR   string        `yaml:"service_cidr,omitempty"`   // 服务网段
	PodCIDR       string        `yaml:"pod_cidr,omitempty"`       // Pod网段
	ImageRepo     string        `yaml:"image_repo,omitempty"`     // 镜像仓库地址（如私有仓库）
	RepoUsername  string        `yaml:"repo_username,omitempty"`  // 仓库用户名
	RepoPassword  string        `yaml:"repo_password,omitempty"`  // 仓库密码
	MaxRetries    int           `yaml:"max_retries,omitempty"`    // 最大重试次数
	RetryDelay    time.Duration `yaml:"retry_delay,omitempty"`    // 重试延迟时间
}

var Cfg Config // 全局配置实例

// LoadConfig 从yaml文件加载配置
func LoadConfig(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		Logger.Fatalf("加载配置文件失败: %v", err)
	}
	if err := yaml.Unmarshal(data, &Cfg); err != nil {
		Logger.Fatalf("解析配置文件失败: %v", err)
	}
}
