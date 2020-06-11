Nginx 是一个 Web Server，可以用作反向代理、负载均衡、邮件代理、TCP / UDP、HTTP 服务器等等，它拥有很多吸引人的特性，例如：

以较低的内存占用率处理 10,000 多个并发连接（每10k非活动HTTP保持活动连接约2.5 MB ）
静态服务器（处理静态文件）
正向、反向代理
负载均衡
通过OpenSSL 对 TLS / SSL 与 SNI 和 OCSP 支持
FastCGI、SCGI、uWSGI 的支持
WebSockets、HTTP/1.1 的支持
Nginx + Lua

常用命令
nginx：启动 Nginx
nginx -s stop：立刻停止 Nginx 服务
nginx -s reload：重新加载配置文件
nginx -s quit：平滑停止 Nginx 服务
nginx -t：测试配置文件是否正确
nginx -v：显示 Nginx 版本信息
nginx -V：显示 Nginx 版本信息、编译器和配置参数的信息
涉及配置
1、 proxy_pass：配置反向代理的路径。需要注意的是如果 proxy_pass 的 url 最后为
/，则表示绝对路径。否则（不含变量下）表示相对路径，所有的路径都会被代理过去

2、 upstream：配置负载均衡，upstream 默认是以轮询的方式进行负载，另外还支持四种模式，分别是：

（1）weight：权重，指定轮询的概率，weight 与访问概率成正比

（2）ip_hash：按照访问 IP 的 hash 结果值分配

（3）fair：按后端服务器响应时间进行分配，响应时间越短优先级别越高

（4）url_hash：按照访问 URL 的 hash 结果值分配

部署
nginx.conf 进行配置，如果你不知道对应的配置文件是哪个，可执行 nginx -t 看一下
反向代理
反向代理是指以代理服务器来接受网络上的连接请求，然后将请求转发给内部网络上的服务器，并将从服务器上得到的结果返回给请求连接的客户端，此时代理服务器对外就表现为一个反向代理服务器。（来自百科）

![Image text](https://github.com/gyb333/practicego/blob/master/middleware/Nginx/images/articlex.png?raw=true)

配置 hosts
由于需要用本机作为演示，因此先把映射配上去，打开 /etc/hosts，增加内容：

127.0.0.1 api.blog.com
配置 nginx.conf
打开 nginx 的配置文件 nginx.conf（我的是 /usr/local/etc/nginx/nginx.conf），我们做了如下事情：

增加 server 片段的内容，设置 server_name 为 api.blog.com 并且监听 8081 端口，将所有路径转发到 http://127.0.0.1:8000/ 下

worker_processes  1;

events {
    worker_connections  1024;
}


http {
    include       mime.types;
    default_type  application/octet-stream;

    sendfile        on;
    keepalive_timeout  65;

    server {
        listen       8081;
        server_name  api.blog.com;

        location / {
            proxy_pass http://127.0.0.1:8000/;
        }
    }
}
验证
启动 go-gin-example
回到 go-gin-example 的项目下，执行 make，再运行 ./go-gin-exmaple

$ make
重启 nginx
$ nginx -t
 $ nginx -s reload
 
 配置 nginx.conf
 回到 nginx.conf 的老地方，增加负载均衡所需的配置。新增 upstream 节点，设置其对应的 2 个后端服务，最后修改了 proxy_pass 指向（格式为 http:// + upstream 的节点名称）
 
 worker_processes  1;
 
 events {
     worker_connections  1024;
 }
 
 
 http {
     include       mime.types;
     default_type  application/octet-stream;
 
     sendfile        on;
     keepalive_timeout  65;
 
     upstream api.blog.com {
         server 127.0.0.1:8001;
         server 127.0.0.1:8002;
     }
 
     server {
         listen       8081;
         server_name  api.blog.com;
 
         location / {
             proxy_pass http://api.blog.com/;
         }
     }
 }
 重启 nginx