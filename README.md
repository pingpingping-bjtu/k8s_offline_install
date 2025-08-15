# k8s_offline_install
1.根据实际修改config/config.yaml
2.编译
```shell

$env:CGO_ENABLED=0
$env:GOOS="linux"
$env:GOARCH="amd64"
# 编译
go build -o k8s-installer.exe ./cmd/k8sinstall
```
3.拷贝到实施主机上
```shell
chmod +x k8s-installer.exe
```
4.运行
```shell
./k8s-installer.exe
```