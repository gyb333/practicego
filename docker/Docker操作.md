#查看进程，发现相关的容器并没有在运行，而 docker-proxy 却依然绑定着端口：
#docker  ps -aux  | grep docker-proxy
#service docker stop
#rm /var/lib/docker/network/files/local-kv.db
#service docker start


docker中 启动所有的容器命令
docker start $(docker ps -a | awk '{ print $1}' | tail -n +2)

docker中 关闭所有的容器命令
docker stop $(docker ps -a | awk '{ print $1}' | tail -n +2)

docker中 删除所有的容器命令
docker rm $(docker ps -a | awk '{ print $1}' | tail -n +2)

docker中 删除所有的镜像
docker rmi $(docker images | awk '{print $3}' |tail -n +2)

//检查网络
docker network ls
新建网络
docker network create --driver bridge --subnet 172.23.0.0/25 --gateway 172.23.0.1  network
 

docker ps -a    查看docker容器内的ID
docker inspect id 查看网络信息
解决方法：
　　1. 首先要保证在虚拟机中能够连接到Docker容器中，用ping测试是否通畅
　　2. 关闭虚拟中的防火墙： systemctl stop firewalld.service
　　3. 打开宿主机（windows）的cmd,在其中添加通往192.168.1.0/24网络的路由。
　　通往192.168.1.0/24网络的数据包由172.20.1.12来转发
    route add 192.168.1.0 mask 255.255.255.0 172.20.1.12
    
在Windows宿主机中连接虚拟机中的Docker容器
route add -p 172.23.0.0 mask 255.255.255.0 192.168.56.100
route add -p 172.23.0.0 mask 255.255.255.0 169.254.23.215
#route delete 172.23.0.0
route print 172.23.0.0
ping 172.23.0.1
需要中间192.168.56.100机器开启IP路由转发功能
要把B机器上的IP转发功能打开，临时的打开方法是“echo 1 > /proc/sys/net/ipv4/ip_forward”，永久的修改，需要修改/etc/sysctl.conf文件。
#vi /etc/sysctl.conf
#net.ipv4.ip_forward=1 
sysctl –p
init 6 
iptables -L
如果192.168.56.100是windows服务器,IP Forwarding打开
#reg query HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Services\Tcpip\Parameters /v IPEnableRouter
#reg add HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Services\Tcpip\Parameters /v IPEnableRouter /t REG_DWORD /d 1 /f
   

docker run -di --name=redis --restart=always -v  /usr/local/docker/redis/redis.conf:/usr/redis/redis.conf \
-v /usr/local/docker/redis/data:/data \
-p 6379:6379 redis redis-server /usr/redis/redis.conf

docker logs redis

docker exec -it redis redis-cli
    config get requirepass
    config set requirepass qwer.1234
    auth qwer.1234  身份认证
    
基础教程：http://redisdoc.com/string/index.html
    select 1    切换数据库

docker 安装mysql主从复制
docker-compose  up -d
docker exec -it mysql-Master bash
docker exec -it mysql-Slave bash
mysql -u root -proot
show master status;
 
CHANGE MASTER TO
    MASTER_HOST='mysql-Master',
    MASTER_USER='root',
    MASTER_PASSWORD='root',
    MASTER_LOG_FILE='replicas-mysql-bin.000002',
    MASTER_LOG_POS=156;
start slave
show slave status\G;
查看下面两项值均为Yes，即表示设置从服务器成功。
Slave_IO_Running: Yes
Slave_SQL_Running: Yes

#create database test;
#SET FOREIGN_KEY_CHECKS=0;

 
docker 安装zookeeper集群

// 校验配置文件，不打印
#docker-compose -f zookeeper-compose.yml config -q

docker-compose.yml 当前目录下运行:

COMPOSE_PROJECT_NAME=zk docker-compose  up -d
COMPOSE_PROJECT_NAME=zk docker-compose ps

COMPOSE_PROJECT_NAME=zk docker-compose stop

COMPOSE_PROJECT_NAME=zk docker-compose rm

COMPOSE_PROJECT_NAME=zk docker-compose -f zookeeper-compose.yml up

docker exec -t zkMaster zkServer.sh status    
docker exec -it zkMaster /bin/bash

docker run -it --rm \
        --net network \
        zookeeper zkCli.sh -server zkMaster:2181,zkSecod:2181,zkSlave:2181  
         
rmr /brokers/topics/<topic_name>        

docker 安装kafka集群
COMPOSE_PROJECT_NAME=kafka docker-compose  up -d
docker-compose ps    
     
docker 安装kafka docker pull wurstmeister/kafka

//启动kafka #--link zookeeper \
docker run -d --name kafka --publish 9092:9092 \
--net network \
--env KAFKA_ZOOKEEPER_CONNECT=zkMaster:2181,zkSecond:2181,zkSlave:2181 \
--env KAFKA_ADVERTISED_HOST_NAME=192.168.56.100 \
--env KAFKA_ADVERTISED_PORT=9092  \
--volume /etc/localtime:/etc/localtime \
wurstmeister/kafka    

测试kafka
docker exec -it kafka /bin/bash
cd opt/kafka
//创建topic
bin/kafka-topics.sh --create --zookeeper zkMaster:2181,zkSecond:2181,zkSlave:2181  --replication-factor 1 --partitions 1 --topic gyb
bin/kafka-topics.sh --create --bootstrap-server localhost:9092 --replication-factor 1 --partitions 1 --topic test
//查看topic
bin/kafka-topics.sh --list --zookeeper zkMaster:2181,zkSecond:2181,zkSlave:2181 
bin/kafka-topics.sh --list --bootstrap-server localhost:9092
//创建生产者
bin/kafka-console-producer.sh --broker-list localhost:9092 --topic gyb 
bin/kafka-console-producer.sh --bootstrap-server localhost:9092 --topic test
//创建消费者
bin/kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic gyb --from-beginning
bin/kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic test --from-beginning

bin/kafka-topics.sh  --zookeeper zkMaster:2181,zkSecond:2181,zkSlave:2181 --delete --topic gyb
bin/kafka-topics.sh  --bootstrap-server localhost:9092 --delete --topic test


docker 安装RocketMQ集群
docker-compose.yml 当前目录下运行:
COMPOSE_PROJECT_NAME=RocketMQ docker-compose  up -d
RocketMQ 控制台
默认访问 http://rmqIP:8080 登入控制台 
http://hadoop:8888


 
docker inspect id 查看网络信息
route add -p 172.23.0.0 mask 255.255.255.0 192.168.56.100
route add -p 172.23.0.0 mask 255.255.255.0 169.254.23.215


docker network create --driver bridge --subnet 172.23.0.0/25 --gateway 172.23.0.1  network


docker 安装ETCD集群
docker-compose -f docker-compose.yml config -q
COMPOSE_PROJECT_NAME=etcd docker-compose  up -d
docker-compose ps
 
#route add -p 172.23.0.0 mask 255.255.255.0 192.168.56.100
创建键值 curl http://172.23.0.20:2379/v2/keys/cqh -XPUT -d value="陈琼和1"
创建目录 curl http://172.23.0.20:2379/v2/keys/gym -XPUT -d dir=true
获取键值 curl http://172.23.0.20:2379/v2/keys/cqh
创建键值带ttl curl http://172.23.0.20:2379/v2/keys/hero -XPUT -d value="超人" -d ttl=5
创建有序键值
curl http://172.23.0.20:2379/v2/keys/fitness -XPOST -d value="bench_press"
curl http://172.23.0.20:2379/v2/keys/fitness -XPOST -d value="dead_lift"
curl http://172.23.0.20:2379/v2/keys/fitness -XPOST -d value="deep_squat"
curl http://172.23.0.20:2379/v2/keys/fitness
删除键 curl http://172.23.0.20:2379/v2/keys/cqh -XDELETE
列出所有集群成员 curl http://172.23.0.20:2379/v2/members
统计信息-查看leader curl http://172.23.0.20:2379/v2/stats/leader
节点自身信息 curl http://172.23.0.20:2379/v2/stats/self
查看集群运行状态 curl http://172.23.0.20:2379/v2/stats/store


#docker exec -t etcd-node1  etcdctl --endpoints=$ETCDENDPOINTS member list
export ETCDCTL_API=3 
ETCDENDPOINTS=127.23.0.20:2379,127.23.0.21:2379,127.23.0.22:2379
#docker exec -t etcd-node1 etcdctl member list
#docker exec -t etcd-node1 etcdctl watch key -f
#docker exec -t etcd-node1 etcdctl set /key value
etcdctl member list
etcdctl set /cqh muscle
etcdctl watch key -forever
#docker exec -t etcd-node1  etcdctl --endpoints=http://127.23.0.20:2379 set lmh "lmh1"

#docker exec -t etcd-node1 etcdctl exec-watch key -- sh -c 'pwd'
#docker exec -t etcd-node1 etcdctl cluster-health

docker run -it -d --name etcdkeeper \
--net network --ip 172.23.0.50  -p 18080:8080 \
deltaprojects/etcdkeeper
访问http://gyb333:18080/etcdkeeper/，输入etcd的地址,看到如下界面


go的etcd v3安装包
go get -u -v go.etcd.io/etcd/clientv3

可视化界面：搭建etcd-browser和etcdkeeper，两者功能大同小异，不同的是etcdkeeper支持v3的api

etcd-browser
#--rm
docker run -itd --name etcd-browser \
--net network --ip 172.23.0.51 -p 18000:8000 \
--env ETCD_HOST=172.23.0.20 \
--env ETCD_PORT=2379 \
buddho/etcd-browser
运行后访问http://gyb333:18000/


安装keepalived：yum install keepalived
systemctl enable keepalived # 开机自启动
systemctl start keepalived     # 启动
systemctl stop keepalived     # 暂停
systemctl restart keepalived  # 重启
systemctl status keepalived   # 查看状态  
tail -f /var/log/messages

docker pause  nginx_master

构建centos镜像：yum install bind bind-utils -y
docker build -t gyb333/dns .
docker build -t gyb333/centos .

安装DNS服务器
docker exec -it dns_master /bin/bash
连接外网：
vi /etc/resolv.conf
nameserver 119.29.29.29
nameserver 114.114.114.114

解决Failed to set locale, defaulting to C.UTF-8：
  echo "export LC_ALL=en_US.UTF8" >> /etc/profile 
  source /etc/profile
  
systemctl enable named
systemctl start named
dig masterdns.gyb333.com

nslookup masterdns.gyb333.com
host masterdns.gyb333.com


docker exec -it lvs01 /bin/bash
ipvsadm -Ln

docker exec -it resty01 /bin/bash
docker exec -it resty02 /bin/bash
ipvsadm -Ln
systemctl enable keepalived
systemctl start keepalived