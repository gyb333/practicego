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

![Image text](https://github.com/gyb333/practicego/middleware/Nginx/images/articlex.png)
