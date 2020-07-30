#!/bin/bash
ifconfig lo:0 172.23.0.100  broadcast 172.23.0.100  netmask 255.255.255.255 up
route add -host 172.23.0.100  dev lo:0
echo "1" > /proc/sys/net/ipv4/conf/lo/arp_ignore
echo "2" > /proc/sys/net/ipv4/conf/lo/arp_announce
echo "1" > /proc/sys/net/ipv4/conf/all/arp_ignore
echo "2" > /proc/sys/net/ipv4/conf/all/arp_announce
sysctl -p &>/dev/null