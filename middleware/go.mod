module GYB.Middleware

go 1.14

require (
	GYB.Common v0.0.0-incompatible
	github.com/Shopify/sarama v1.26.4
	github.com/apache/rocketmq-client-go/v2 v2.0.0
	github.com/bsm/sarama-cluster v2.1.15+incompatible
	github.com/coreos/bbolt v1.3.5 // indirect
	github.com/coreos/etcd v3.3.22+incompatible
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/garyburd/redigo v1.6.0
	github.com/golang/protobuf v1.4.2
	github.com/google/btree v1.0.0 // indirect
	github.com/google/uuid v1.1.1 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.0 // indirect
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.14.6 // indirect
	github.com/jonboulle/clockwork v0.2.0 // indirect
	github.com/onsi/ginkgo v1.13.0 // indirect
	github.com/prometheus/client_golang v1.7.1 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20190826022208-cac0b30c2563
	github.com/soheilhy/cmux v0.1.4 // indirect
	github.com/streadway/amqp v0.0.0-20200108173154-1c71cc93ed71
	github.com/tmc/grpc-websocket-proxy v0.0.0-20200427203606-3cfed13b9966 // indirect
	github.com/xiang90/probing v0.0.0-20190116061207-43a291ad63a2 // indirect
	go.etcd.io/bbolt v1.3.5 // indirect
	go.etcd.io/etcd v3.3.22+incompatible // indirect
	go.uber.org/zap v1.15.0 // indirect
	golang.org/x/net v0.0.0-20200520004742-59133d7f0dd7
	golang.org/x/time v0.0.0-20200630173020-3af7569d3a1e // indirect
	google.golang.org/grpc v1.29.1
	google.golang.org/protobuf v1.23.0 // indirect
	gopkg.in/olivere/elastic.v5 v5.0.85
	sigs.k8s.io/yaml v1.2.0 // indirect

)

replace go.etcd.io/bbolt v1.3.5 => github.com/coreos/bbolt v1.3.5

replace github.com/coreos/bbolt v1.3.5 => go.etcd.io/bbolt v1.3.5

replace google.golang.org/grpc v1.29.1 => google.golang.org/grpc v1.25.1

replace GYB.Common => ../Common
