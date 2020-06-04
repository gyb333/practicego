module GYB.Middleware

go 1.14

require (
	github.com/garyburd/redigo v1.6.0
	github.com/streadway/amqp v0.0.0-20200108173154-1c71cc93ed71
	gopkg.in/olivere/elastic.v5 v5.0.85

	GYB.Common v0.0.0-incompatible
)
replace GYB.Common => ../Common
