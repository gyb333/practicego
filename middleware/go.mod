module GYB.Middleware

go 1.14

require (
	GYB.Common v0.0.0-incompatible
	github.com/Shopify/sarama v1.26.4
	github.com/apache/rocketmq-client-go/v2 v2.0.0
	github.com/bsm/sarama-cluster v2.1.15+incompatible
	github.com/garyburd/redigo v1.6.0
	github.com/onsi/ginkgo v1.13.0 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20190826022208-cac0b30c2563
	github.com/streadway/amqp v0.0.0-20200108173154-1c71cc93ed71
	gopkg.in/olivere/elastic.v5 v5.0.85
)

replace GYB.Common => ../Common
