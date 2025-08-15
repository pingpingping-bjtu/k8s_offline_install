package k8sinstall

import (
	"fmt"
	"k8s-offline-installer/internal/util"
)

func InstallContainerd() {
	util.Logger.Infof("开始安装containerd...")
	pkg := fmt.Sprintf("containerd-%s-linux-amd64.tar.gz", util.Cfg.ContainerdVer)
	downloadURL := fmt.Sprintf("https://github.com/containerd/containerd/releases/download/v%s/%s", util.Cfg.ContainerdVer, pkg)

	// 下载并上传安装包
	if err := util.DownloadFile(downloadURL, pkg); err != nil {
		util.Logger.Fatalf("下载containerd失败: %v", err)
	}
	if err := util.ScpToRemote(pkg); err != nil {
		util.Logger.Fatalf("上传containerd安装包失败: %v", err)
	}
	defer util.CleanLocalFiles(pkg) // 清理本地安装包

	// 解压到/usr/
	err := util.ExecRemote(fmt.Sprintf("cd %s && tar -zxvf %s -C /usr/", util.Cfg.TargetPath, pkg))
	if err != nil {
		return
	}

	// 生成配置文件（使用变量替换硬编码仓库地址）
	configScript := fmt.Sprintf(`
mkdir -p /etc/containerd/
containerd config default > /etc/containerd/config.toml
cat <<EOF | sudo tee /etc/containerd/config.toml > /dev/null
version = 2
root = "/data/containerd"
state = "/run/containerd"
oom_score = 0
[grpc]
  max_recv_message_size = 16777216
  max_send_message_size = 16777216
[debug]
  address = ""
  level = "info"
  format = ""
  uid = 0
  gid = 0
[metrics]
  address = ""
  grpc_histogram = false
[plugins]
  [plugins."io.containerd.cri.v1.runtime"]
    max_container_logger_line_size = 16384
    enable_unprivileged_ports = false
    enable_unprivileged_icmp = false
    enable_selinux = false
    disable_apparmor = false
    tolerate_missing_hugetlb_controller = true
    disable_hugetlb_controller = true
   [plugins."io.containerd.cri.v1.runtime".containerd]
     default_runtime_name = "runc"
     [plugins."io.containerd.cri.v1.runtime".containerd.runtimes]
       [plugins."io.containerd.cri.v1.runtime".containerd.runtimes.runc]
         runtime_type = "io.containerd.runc.v2"
         runtime_engine = ""
         runtime_root = ""
         base_runtime_spec = "/etc/containerd/cri-base.json"
         [plugins."io.containerd.cri.v1.runtime".containerd.runtimes.runc.options]
           SystemdCgroup = true
           BinaryName = "/usr/bin/runc"
  [plugins."io.containerd.cri.v1.images"]
    snapshotter = "overlayfs"
    discard_unpacked_layers = true
    image_pull_progress_timeout = "5m"
  [plugins."io.containerd.cri.v1.images".pinned_images]
    sandbox = "registry.k8s.io/pause:3.10" #自己的镜像仓库名/pause:3.10
  [plugins."io.containerd.cri.v1.images".registry]
    config_path = "/etc/containerd/cert.d"
  [plugins."io.containerd.nri.v1.nri"]
    disable = false
#安装harbor或其他私有镜像仓库后添加配置
   #   [plugins."io.containerd.grpc.v1.cri".registry.mirrors."%s"]
   #     endpoint = ["http://%s"]    #填自己的镜像仓库
#[plugins."io.containerd.grpc.v1.cri".registry.configs."%s".auth]
  # 配置访问镜像仓库的用户名
 # username = "%s"
  # 配置访问镜像仓库的密码（默认密码，如果之前改了这里也要改）
 # password = "%s"
EOF
	`, util.Cfg.ImageRepo, util.Cfg.ImageRepo, util.Cfg.ImageRepo, util.Cfg.RepoUsername, util.Cfg.RepoPassword)
	err = util.ExecRemote(configScript)
	if err != nil {
		return
	}
	//配置"/etc/containerd/cri-base.json"所提及的cri-base.json文件
	configScript2 := `
cat <<EOF | sudo tee /etc/containerd/cri-base.json > /dev/null
{
  "ociVersion": "1.2.0",
  "process": {
    "user": {
      "uid": 0,
      "gid": 0
    },
    "cwd": "/",
    "capabilities": {
      "bounding": [
        "CAP_CHOWN",
        "CAP_DAC_OVERRIDE",
        "CAP_FSETID",
        "CAP_FOWNER",
        "CAP_MKNOD",
        "CAP_NET_RAW",
        "CAP_SETGID",
        "CAP_SETUID",
        "CAP_SETFCAP",
        "CAP_SETPCAP",
        "CAP_NET_BIND_SERVICE",
        "CAP_SYS_CHROOT",
        "CAP_KILL",
        "CAP_AUDIT_WRITE"
      ],
      "effective": [
        "CAP_CHOWN",
        "CAP_DAC_OVERRIDE",
        "CAP_FSETID",
        "CAP_FOWNER",
        "CAP_MKNOD",
        "CAP_NET_RAW",
        "CAP_SETGID",
        "CAP_SETUID",
        "CAP_SETFCAP",
        "CAP_SETPCAP",
        "CAP_NET_BIND_SERVICE",
        "CAP_SYS_CHROOT",
        "CAP_KILL",
        "CAP_AUDIT_WRITE"
      ],
      "permitted": [
        "CAP_CHOWN",
        "CAP_DAC_OVERRIDE",
        "CAP_FSETID",
        "CAP_FOWNER",
        "CAP_MKNOD",
        "CAP_NET_RAW",
        "CAP_SETGID",
        "CAP_SETUID",
        "CAP_SETFCAP",
        "CAP_SETPCAP",
        "CAP_NET_BIND_SERVICE",
        "CAP_SYS_CHROOT",
        "CAP_KILL",
        "CAP_AUDIT_WRITE"
      ]
    },
    "rlimits": [
      {
        "type": "RLIMIT_NOFILE",
        "hard": 65535,
        "soft": 65535
      }
    ],
    "noNewPrivileges": true
  },
  "root": {
    "path": "rootfs"
  },
  "mounts": [
    {
      "destination": "/proc",
      "type": "proc",
      "source": "proc",
      "options": ["nosuid", "noexec", "nodev"]
    },
    {
      "destination": "/dev",
      "type": "tmpfs",
      "source": "tmpfs",
      "options": ["nosuid", "strictatime", "mode=755", "size=65536k"]
    },
    {
      "destination": "/dev/pts",
      "type": "devpts",
      "source": "devpts",
      "options": ["nosuid", "noexec", "newinstance", "ptmxmode=0666", "mode=0620", "gid=5"]
    },
    {
      "destination": "/dev/shm",
      "type": "tmpfs",
      "source": "shm",
      "options": ["nosuid", "noexec", "nodev", "mode=1777", "size=65536k"]
    },
    {
      "destination": "/dev/mqueue",
      "type": "mqueue",
      "source": "mqueue",
      "options": ["nosuid", "noexec", "nodev"]
    },
    {
      "destination": "/sys",
      "type": "sysfs",
      "source": "sysfs",
      "options": ["nosuid", "noexec", "nodev", "ro"]
    },
    {
      "destination": "/run",
      "type": "tmpfs",
      "source": "tmpfs",
      "options": ["nosuid", "strictatime", "mode=755", "size=65536k"]
    }
  ],
  "linux": {
    "resources": {
      "devices": [
        {
          "allow": false,
          "access": "rwm"
        }
      ]
    },
    "cgroupsPath": "/default",
    "namespaces": [
      {
        "type": "pid"
      },
      {
        "type": "ipc"
      },
      {
        "type": "uts"
      },
      {
        "type": "mount"
      },
      {
        "type": "network"
      }
    ],
    "maskedPaths": [
      "/proc/acpi",
      "/proc/asound",
      "/proc/kcore",
      "/proc/keys",
      "/proc/latency_stats",
      "/proc/timer_list",
      "/proc/timer_stats",
      "/proc/sched_debug",
      "/sys/firmware",
      "/sys/devices/virtual/powercap",
      "/proc/scsi"
    ],
    "readonlyPaths": [
      "/proc/bus",
      "/proc/fs",
      "/proc/irq",
      "/proc/sys",
      "/proc/sysrq-trigger"
    ]
  }
}

EOF
	`
	err = util.ExecRemote(configScript2)
	if err != nil {
		return
	}
	// 配置系统服务并启动
	serviceScript := `
cat <<EOF | sudo tee /usr/lib/systemd/system/containerd.service > /dev/null
[Unit]
Description=containerd container runtime
After=network.target
[Service]
ExecStartPre=/sbin/modprobe overlay
ExecStart=/usr/bin/containerd --config /etc/containerd/config.toml
Restart=always
RestartSec=5
LimitNOFILE=infinity
[Install]
WantedBy=multi-user.target
EOF
systemctl daemon-reload && systemctl enable --now containerd && systemctl restart containerd
`
	err = util.ExecRemote(serviceScript)
	if err != nil {
		return
	}
	util.Logger.Infof("containerd安装完成")
}
