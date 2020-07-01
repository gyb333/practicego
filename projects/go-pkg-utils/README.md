go mod init go-pkg-utils

docker ps -a    查看docker容器内的ID
docker inspect id 查看网络信息
解决方法：
　　1. 首先要保证在虚拟机中能够连接到Docker容器中，用ping测试是否通畅
　　2. 关闭虚拟中的防火墙： systemctl stop firewalld.service
　　3. 打开宿主机（windows）的cmd,在其中添加通往192.168.1.0/24网络的路由。
　　通往192.168.1.0/24网络的数据包由172.20.1.12来转发
    route add 192.168.1.0 mask 255.255.255.0 172.20.1.12
    
在Windows宿主机中连接虚拟机中的Docker容器
route add -p 172.20.0.0 mask 255.255.255.0 192.168.56.100
#route delete 172.20.0.0
route print 172.20.0.0
ping 172.20.0.1

docker run -di --name=redis -v  /usr/local/docker/redis/redis.conf:/usr/redis/redis.conf \
-v /usr/local/docker/redis/data:/data \
-p 6379:6379 redis redis-server /usr/redis/redis.conf

docker logs redis

docker exec -it redis redis-cli
    config get requirepass
    config set requirepass qwer.1234
    auth qwer.1234  身份认证
    
基础教程：http://redisdoc.com/string/index.html
    select 1    切换数据库

docker 安装zookeeper
//检查网络
docker network ls
新建网络
docker network create --driver bridge --subnet 172.23.0.0/25 --gateway 172.23.0.1  zookeeper_network

// 校验配置文件，不打印
$ docker-compose -f zookeeper-compose.yml config -q

docker-compose.yml 当前目录下运行:

COMPOSE_PROJECT_NAME=zk docker-compose  up -d
COMPOSE_PROJECT_NAME=zk docker-compose ps

COMPOSE_PROJECT_NAME=zk docker-compose stop

COMPOSE_PROJECT_NAME=zk docker-compose rm

COMPOSE_PROJECT_NAME=zk docker-compose -f zookeeper-compose.yml up
    
docker exec -it zkMaster /bin/bash

docker run -it --rm \
        --net network \
        zookeeper zkCli.sh -server zkMaster:2181,zkSecod:2181,zkSlave:2181  
         
rmr /brokers/topics/<topic_name>        
         
docker 安装kafka
ocker pull wurstmeister/kafka

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
bin/kafka-console-producer.sh --broker-list kafka:9092 --topic gyb 
bin/kafka-console-producer.sh --bootstrap-server localhost:9092 --topic test
//创建消费者
bin/kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic gyb --from-beginning
bin/kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic test --from-beginning

bin/kafka-topics.sh  --zookeeper zkMaster:2181,zkSecond:2181,zkSlave:2181 --delete --topic gyb

docker 安装RocketMQ

docker-compose.yml 当前目录下运行:
COMPOSE_PROJECT_NAME=RocketMQ docker-compose  up -d
RocketMQ 控制台
默认访问 http://rmqIP:8080 登入控制台 
http://hadoop:8888


COMPOSE_PROJECT_NAME=hadoop docker-compose  up -d
docker inspect id 查看网络信息
route add -p 172.23.0.0 mask 255.255.255.0 192.168.56.100

