go mod init go-pkg-utils

docker run -di --name=redis -v  /usr/local/conf/redis/redis.conf:/usr/redis/redis.conf -p 6379:6379 redis redis-server /usr/redis/redis.conf
docker exec -it redis redis-cli
    config get requirepass
    config set requirepass qwer.1234
    auth qwer.1234  身份认证
    
http://redisdoc.com/string/index.html
    select 1    切换数据库