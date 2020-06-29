package test

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/rcrowley/go-metrics"
	"os"
	"testing"
)

func TestBroker(t *testing.T)  {
	Broker()
}

func Broker()  {
	broker := sarama.NewBroker("hadoop:9092")
	err := broker.Open(nil)
	if err != nil {
		panic(err)
	}

	request := sarama.MetadataRequest{Topics: []string{"myTopic"}}
	response, err := broker.GetMetadata(&request)
	if err != nil {
		_ = broker.Close()
		panic(err)
	}

	fmt.Println("There are", len(response.Topics), "topics active in the cluster.")

	if err = broker.Close(); err != nil {
		panic(err)
	}
}

func TestMetrics(t *testing.T)  {
	Metrics()
}

func Metrics()  {
	// Our application registry
	appMetricRegistry := metrics.NewRegistry()
	appGauge := metrics.GetOrRegisterGauge("m1", appMetricRegistry)
	appGauge.Update(1)

	config := sarama.NewConfig()
	// Use a prefix registry instead of the default local one
	config.MetricRegistry = metrics.NewPrefixedChildRegistry(appMetricRegistry, "sarama.")

	// Simulate a metric created by sarama without starting a broker
	saramaGauge := metrics.GetOrRegisterGauge("m2", config.MetricRegistry)
	saramaGauge.Update(2)

	metrics.WriteOnce(appMetricRegistry, os.Stdout)
}