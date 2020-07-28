docker exec -t etcd-node1 etcdctl member list
docker exec -it etcd-node1 bin/sh
etcdctl set /user/101/name xiahualou
etcdctl get /user/101/name
etcdctl get /user/101 --prefix
etcdctl rm /user/101/name