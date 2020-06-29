go mod init go-pkg-utils

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

