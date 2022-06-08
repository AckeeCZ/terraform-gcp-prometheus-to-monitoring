package p

import (
	"context"
	"fmt"
	"math"
	"os"
	"time"

	monitoring "cloud.google.com/go/monitoring/apiv3"
	googlepb "github.com/golang/protobuf/ptypes/timestamp"
	dto "github.com/prometheus/client_model/go"
	zlog "github.com/rs/zerolog/log"
	"google.golang.org/genproto/googleapis/api/distribution"
	metricpb "google.golang.org/genproto/googleapis/api/metric"
	monitoredrespb "google.golang.org/genproto/googleapis/api/monitoredres"
	monitoringpb "google.golang.org/genproto/googleapis/monitoring/v3"
)

func sendToGCPMonitoring(mf map[string]*dto.MetricFamily, config map[string]*string) {
	ctx := context.Background()
	ts := []*monitoringpb.TimeSeries{}
	start := time.Now()
	end := start.Add(time.Second)
	for _, v := range mf {
		if v.GetType() == dto.MetricType_GAUGE {
			for _, m := range v.GetMetric() {
				if m == nil {
					continue
				}
				labels := make(map[string]string)
				for _, l := range m.GetLabel() {
					labels[*l.Name] = *l.Value
				}
				ts = append(ts, &monitoringpb.TimeSeries{
					Metric: &metricpb.Metric{
						Type:   fmt.Sprintf("custom.googleapis.com/%s/%s", *config["SERVICE"], *v.Name),
						Labels: labels,
					},
					Resource: &monitoredrespb.MonitoredResource{
						Type: "global",
					},
					Points: []*monitoringpb.Point{
						{
							Interval: &monitoringpb.TimeInterval{
								EndTime: &googlepb.Timestamp{
									Seconds: end.Unix(),
								},
							},
							Value: &monitoringpb.TypedValue{
								Value: &monitoringpb.TypedValue_DoubleValue{
									DoubleValue: m.Gauge.GetValue(),
								},
							},
						},
					},
					MetricKind: metricpb.MetricDescriptor_GAUGE,
				})
			}
		} else if v.GetType() == dto.MetricType_COUNTER {
			for _, m := range v.GetMetric() {
				if m == nil {
					continue
				}
				labels := make(map[string]string)
				for _, l := range m.GetLabel() {
					labels[*l.Name] = *l.Value
				}
				ts = append(ts, &monitoringpb.TimeSeries{
					Metric: &metricpb.Metric{
						Type:   fmt.Sprintf("custom.googleapis.com/%s/%s", *config["SERVICE"], *v.Name),
						Labels: labels,
					},
					Resource: &monitoredrespb.MonitoredResource{
						Type: "global",
					},
					Points: []*monitoringpb.Point{
						{
							Interval: &monitoringpb.TimeInterval{
								EndTime: &googlepb.Timestamp{
									Seconds: end.Unix(),
								},
							},
							Value: &monitoringpb.TypedValue{
								Value: &monitoringpb.TypedValue_DoubleValue{
									DoubleValue: m.Counter.GetValue(),
								},
							},
						},
					},
					MetricKind: metricpb.MetricDescriptor_GAUGE,
				})
			}
		} else if v.GetType() == dto.MetricType_HISTOGRAM {
			for _, m := range v.GetMetric() {
				if m == nil {
					continue
				}
				labels := make(map[string]string)
				for _, l := range m.GetLabel() {
					labels[*l.Name] = *l.Value
				}
				var bucketCounts []int64
				var bounds []float64
				var count int64
				for _, v := range m.Histogram.GetBucket() {
					if v.GetUpperBound() != math.Inf(1) {
						bucketCounts = append(bucketCounts, int64(v.GetCumulativeCount()))
						bounds = append(bounds, v.GetUpperBound())
						count += int64(v.GetCumulativeCount())
					}
				}
				ts = append(ts, &monitoringpb.TimeSeries{
					Metric: &metricpb.Metric{
						Type:   fmt.Sprintf("custom.googleapis.com/%s/%s", *config["SERVICE"], *v.Name),
						Labels: labels,
					},
					Resource: &monitoredrespb.MonitoredResource{
						Type: "global",
					},
					Points: []*monitoringpb.Point{
						{
							Interval: &monitoringpb.TimeInterval{
								StartTime: &googlepb.Timestamp{
									Seconds: start.Unix(),
								},
								EndTime: &googlepb.Timestamp{
									Seconds: end.Unix(),
								},
							},
							Value: &monitoringpb.TypedValue{
								Value: &monitoringpb.TypedValue_DistributionValue{
									DistributionValue: &distribution.Distribution{
										BucketOptions: &distribution.Distribution_BucketOptions{
											Options: &distribution.Distribution_BucketOptions_ExplicitBuckets{
												ExplicitBuckets: &distribution.Distribution_BucketOptions_Explicit{
													Bounds: bounds,
												},
											},
										},
										BucketCounts: bucketCounts,
										Count:        count,
									},
								},
							},
						},
					},
					MetricKind: metricpb.MetricDescriptor_CUMULATIVE,
				})
			}
		}
	}
	for i := 0; i < len(ts); i += 200 {
		end := i + 200
		if end > len(ts) {
			end = len(ts)
		}
		client, err := monitoring.NewMetricClient(ctx)
		if err != nil {
			zlog.Error().Msg(fmt.Sprintf("Failed to create client: %v", err))
			os.Exit(1)
		}
		if err := client.CreateTimeSeries(ctx, &monitoringpb.CreateTimeSeriesRequest{
			Name:       fmt.Sprintf("projects/%s", *config["PROJECT_ID"]),
			TimeSeries: ts[i:end],
		}); err != nil {
			zlog.Error().Msg(fmt.Sprintf("Failed to write time series data: %v", err))
			os.Exit(1)
		}
		if err := client.Close(); err != nil {
			zlog.Error().Msg(fmt.Sprintf("Failed to write time series data: %v", err))
			os.Exit(1)
		}
	}
}
