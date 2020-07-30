FROM centos:7
ENV container docker
RUN (cd /lib/systemd/system/sysinit.target.wants/; for i in *; do [ $i == \
systemd-tmpfiles-setup.service ] || rm -f $i; done); \
rm -f /lib/systemd/system/multi-user.target.wants/*;\
rm -f /etc/systemd/system/*.wants/*;\
rm -f /lib/systemd/system/local-fs.target.wants/*; \
rm -f /lib/systemd/system/sockets.target.wants/*udev*; \
rm -f /lib/systemd/system/sockets.target.wants/*initctl*; \
rm -f /lib/systemd/system/basic.target.wants/*;\
rm -f /lib/systemd/system/anaconda.target.wants/*;

RUN yum install bind bind-utils ipvsadm keepalived net-tools -y
#RUN yum clean all; 
#RUN systemctl enable keepalived;

ADD ./start.sh /usr/local/start.sh
RUN chmod +x /usr/local/start.sh
RUN sh /usr/local/start.sh

VOLUME [ "/sys/fs/cgroup" ]
CMD ["/usr/sbin/init"]