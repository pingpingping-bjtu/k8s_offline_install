package util

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// PullImageWithRetry 拉取镜像并添加重试机制
func PullImageWithRetry(image string) error {
	Logger.Infof("开始拉取镜像: %s", image)
	// 准备拉取命令
	pullCmd := []string{
		"docker", "pull", image,
	}
	// 带重试的拉取操作
	for i := 0; i < Cfg.MaxRetries; i++ {
		if err := ExecLocal(pullCmd...); err != nil {
			// 最后一次失败则返回错误
			if i == Cfg.MaxRetries-1 {
				return fmt.Errorf("镜像 %s 拉取失败（已尝试%d次）: %v", image, Cfg.MaxRetries, err)
			}
			Logger.Errorf("镜像 %s 拉取失败（第%d次重试）: %v", image, i+1, err)
			time.Sleep(Cfg.RetryDelay)
			continue
		}
		Logger.Infof("镜像 %s 拉取成功", image)
		return nil
	}

	return fmt.Errorf("镜像 %s 拉取失败（达到最大重试次数）", image)
}

// PullAndTransferImages 拉取、打包、上传镜像
func PullAndTransferImages(images []string) {
	Logger.Infof("开始处理镜像（拉取->打包->上传）...")
	localTempDir := "./k8s_images_temp"
	err := os.MkdirAll(localTempDir, 0755)
	if err != nil {
		return
	} // 创建临时目录
	for _, image := range images {
		// 1. 拉取镜像
		if err := PullImageWithRetry(image); err != nil {
			Logger.Fatalf("处理镜像 %s 失败: %v", image, err)
		}
		// 2. 保存镜像为tar包
		tarName := filepath.Base(image) + ".tar"
		localTarPath := filepath.Join(localTempDir, tarName)
		saveCmd := []string{"docker", "save", "-o", localTarPath, image}
		if err := ExecLocal(saveCmd...); err != nil {
			Logger.Fatalf("保存镜像 %s 失败: %v", image, err)
		}

		// 3. 上传到目标主机
		if err := ScpToRemote(localTarPath); err != nil {
			Logger.Fatalf("上传镜像 %s 失败: %v", tarName, err)
		}

		// 4. 远程导入镜像
		remoteImportCmd := fmt.Sprintf(
			"cd %s && ctr -n k8s.io images import %s && rm -f %s",
			Cfg.TargetPath, tarName, tarName,
		)
		if err := ExecRemote(remoteImportCmd); err != nil {
			Logger.Fatalf("远程导入镜像 %s 失败: %v", tarName, err)
		}
		Logger.Infof("镜像 %s 处理完成", image)
	}

	// 清理本地临时文件
	err = os.RemoveAll(localTempDir)
	if err != nil {
		return
	}
	Logger.Infof("镜像处理完成，本地临时文件已清理")

}
