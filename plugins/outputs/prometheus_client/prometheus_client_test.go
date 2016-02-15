package prometheus_client

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs/prometheus"
	"github.com/influxdata/telegraf/testutil"
)


// As there's currently no (easy) way to instantiate several PrometheusClient
// (because of single /metrics Handle registration), need a single
// TestAll() (rather than TestFooBar for each needed test)

func TestAll(t *testing.T) {
	var pTesting *PrometheusClient
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	pTesting = &PrometheusClient{Listen: "localhost:9127"}
	err := pTesting.Start()
	time.Sleep(time.Millisecond * 200)
	require.NoError(t, err)
	defer pTesting.Stop()
	p := &prometheus.Prometheus{
		Urls: []string{"http://localhost:9127/metrics"},
	}
	testPrometheusTagsFields(t, pTesting, p)
	testPrometheusWritePointEmptyTag(t, pTesting, p)
}

func testPrometheusTagsFields(t *testing.T, pTesting *PrometheusClient, p *prometheus.Prometheus) {
	tags := map[string]string{
		"service":  "foo",
	}
	fields := map[string]interface{}{
		"read":   111,
		"write":  222,
	}
	m1, _ := telegraf.NewMetric("my_metric", tags, fields)
	var metrics = []telegraf.Metric{
		m1,
	}
	require.NoError(t, pTesting.Write(metrics))
	expected := []struct {
		name  string
		value float64
		tags  map[string]string
	}{
		{"my_metric_read", 111, map[string]string{"service": "foo"}},
		{"my_metric_write", 222, map[string]string{"service": "foo"}},
	}
	var acc testutil.Accumulator
	require.NoError(t, p.Gather(&acc))
	for _, e := range expected {
		acc.AssertContainsFields(t, "prometheus_"+e.name,
		map[string]interface{}{"value": e.value})
	}
}

func testPrometheusWritePointEmptyTag(t *testing.T, pTesting *PrometheusClient, p *prometheus.Prometheus) {
	tags := make(map[string]string)
	pt1, _ := telegraf.NewMetric(
		"test_point_1",
		tags,
		map[string]interface{}{"value": 0.0})
	pt2, _ := telegraf.NewMetric(
		"test_point_2",
		tags,
		map[string]interface{}{"value": 1.0})
	var metrics = []telegraf.Metric{
		pt1,
		pt2,
	}
	require.NoError(t, pTesting.Write(metrics))

	expected := []struct {
		name  string
		value float64
		tags  map[string]string
	}{
		{"test_point_1_value", 0.0, tags},
		{"test_point_2_value", 1.0, tags},
	}

	var acc testutil.Accumulator

	require.NoError(t, p.Gather(&acc))
	for _, e := range expected {
		acc.AssertContainsFields(t, "prometheus_"+e.name,
			map[string]interface{}{"value": e.value})
	}

	tags = make(map[string]string)
	tags["testtag"] = "testvalue"
	pt3, _ := telegraf.NewMetric(
		"test_point_3",
		tags,
		map[string]interface{}{"value": 0.0})
	pt4, _ := telegraf.NewMetric(
		"test_point_4",
		tags,
		map[string]interface{}{"value": 1.0})
	metrics = []telegraf.Metric{
		pt3,
		pt4,
	}
	require.NoError(t, pTesting.Write(metrics))

	expected2 := []struct {
		name  string
		value float64
	}{
		{"test_point_3_value", 0.0},
		{"test_point_4_value", 1.0},
	}

	require.NoError(t, p.Gather(&acc))
	for _, e := range expected2 {
		acc.AssertContainsFields(t, "prometheus_"+e.name,
			map[string]interface{}{"value": e.value})
	}
}

