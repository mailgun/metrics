// Package metrics provides tools for emitting metrics to StatsD.
package metrics

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mailgun/go-statsd-client/statsd"
	"github.com/mailgun/gotools-log"
)

type MetricsService struct {
	// config information
	period time.Duration

	// statsd remote endpoint
	client *statsd.Client
	url    string
}

func NewMetricsService(host string, port int, period time.Duration, id string) (*MetricsService, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	// format parameters
	hostPort := fmt.Sprintf("%v:%v", host, port)
	prefix := fmt.Sprintf("%v.%v", id, strings.Replace(hostname, ".", "_", -1))

	// start service
	client, err := statsd.New(hostPort, prefix)
	if err != nil {
		return nil, err
	}

	ms := &MetricsService{
		url:    hostPort,
		client: client,
		period: period,
	}

	log.Infof("[+] Started metrics service, emitting metrics to: %v [%v period]",
		hostPort, period)

	return ms, nil
}

func (ms *MetricsService) Stop() error {
	return ms.client.Close()
}

func (ms *MetricsService) GetEmitPeriod() time.Duration {
	if ms.period == 0 {
		return 30 * time.Second
	}
	return ms.period
}

func (ms *MetricsService) EmitGauge(bucket string, value int64) error {
	if ms.client == nil {
		return fmt.Errorf("[-] Metrics service is not started")
	}

	// send metric
	err := ms.client.Gauge(bucket, value, 1.0)
	if err != nil {
		return err
	}

	return nil
}

func (ms *MetricsService) EmitTimer(bucket string, value time.Duration) error {
	if ms.client == nil {
		return fmt.Errorf("[-] Metrics service is not started")
	}

	// send metric in milliseconds (time.Duration is nanoseconds)
	err := ms.client.Timing(bucket, int64(value)/1000000, 1.0)
	if err != nil {
		return err
	}

	return nil
}

func (ms *MetricsService) EmitCounter(bucket string, value int64) error {
	if ms.client == nil {
		return fmt.Errorf("[-] Metrics service is not started")
	}

	// send metric
	err := ms.client.Inc(bucket, value, 1.0)
	if err != nil {
		return err
	}

	return nil
}

