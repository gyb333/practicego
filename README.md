# practicego
golang 实践

下载go(我的当前目录是/data/work)
$wget https://studygolang.com/dl/golang/go1.10.1.linux-amd64.tar.gz
$tar -xvf go1.10.1.linux-amd64.tar.gz

设置环境变量
$vim /etc/profile
添加
export GOROOT=/data/work/go
export GOPATH=/data/work/gopath
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin

保存
esc
:wq