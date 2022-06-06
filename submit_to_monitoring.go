package main

import (
	"context"
	"fmt"
	"log"
	"time"

	monitoring "cloud.google.com/go/monitoring/apiv3/v2"
	googlepb "github.com/golang/protobuf/ptypes/timestamp"
	dto "github.com/prometheus/client_model/go"
	"google.golang.org/genproto/googleapis/api/distribution"
	metricpb "google.golang.org/genproto/googleapis/api/metric"
	monitoredrespb "google.golang.org/genproto/googleapis/api/monitoredres"
	monitoringpb "google.golang.org/genproto/googleapis/monitoring/v3"
)

func sendToGCPMonitoring(mf map[string]*dto.MetricFamily, config map[string]*string) {
	ctx := context.Background()
	client, err := monitoring.NewMetricClient(ctx)
	ts := []*monitoringpb.TimeSeries{}
	start := time.Now()
	end := start.Add(time.Second)
	for metric, v := range mf {
		fmt.Println("KEY: ", metric)
		fmt.Println("VAL: ", v)
		fmt.Println(v.GetType())
		labels := make(map[string]string)

		// if v.GetType() == dto.MetricType_GAUGE {
		// 	for _, m := range v.GetMetric() {
		// 		if m == nil {
		// 			continue
		// 		}
		// 		labels := make(map[string]string)
		// 		for _, l := range m.GetLabel() {
		// 			labels[*l.Name] = *l.Value
		// 		}

		// 		client, err := monitoring.NewMetricClient(ctx)
		// 		if err != nil {
		// 			log.Fatalf("Failed to create client: %v", err)
		// 		}
		// 		if err := client.CreateTimeSeries(ctx, &monitoringpb.CreateTimeSeriesRequest{
		// 			Name: fmt.Sprintf("projects/%s", *config["PROJECT_ID"]),
		// 			TimeSeries: []*monitoringpb.TimeSeries{
		// 				{
		// 					Metric: &metricpb.Metric{
		// 						Type:   fmt.Sprintf("custom.googleapis.com/%s/%s", *config["SERVICE"], *v.Name),
		// 						Labels: labels,
		// 					},
		// 					Resource: &monitoredrespb.MonitoredResource{
		// 						Type: "global",
		// 					},
		// 					Points: []*monitoringpb.Point{
		// 						{
		// 							Interval: &monitoringpb.TimeInterval{
		// 								EndTime: &googlepb.Timestamp{
		// 									Seconds: time.Now().Unix(),
		// 								},
		// 							},
		// 							Value: &monitoringpb.TypedValue{
		// 								Value: &monitoringpb.TypedValue_DoubleValue{
		// 									DoubleValue: m.Gauge.GetValue(),
		// 								},
		// 							},
		// 						},
		// 					},
		// 					MetricKind: metricpb.MetricDescriptor_GAUGE,
		// 				},
		// 			},
		// 		}); err != nil {
		// 			log.Fatalf("Failed to write time series data: %v", err)
		// 		}
		// 		time.Sleep(1000)
		// 		if err := client.Close(); err != nil {
		// 			log.Fatalf("Failed to close client: %v", err)
		// 		}
		// 	}
		// }

		// if v.GetType() == dto.MetricType_COUNTER {
		// 	for _, m := range v.GetMetric() {
		// 		if m == nil {
		// 			continue
		// 		}
		// 		labels := make(map[string]string)
		// 		for _, l := range m.GetLabel() {
		// 			labels[*l.Name] = *l.Value
		// 		}

		// 		client, err := monitoring.NewMetricClient(ctx)
		// 		if err != nil {
		// 			log.Fatalf("Failed to create client: %v", err)
		// 		}
		// 		if err := client.CreateTimeSeries(ctx, &monitoringpb.CreateTimeSeriesRequest{
		// 			Name: fmt.Sprintf("projects/%s", *config["PROJECT_ID"]),
		// 			TimeSeries: []*monitoringpb.TimeSeries{
		// 				{
		// 					Metric: &metricpb.Metric{
		// 						Type:   fmt.Sprintf("custom.googleapis.com/%s/%s", *config["SERVICE"], *v.Name),
		// 						Labels: labels,
		// 					},
		// 					Resource: &monitoredrespb.MonitoredResource{
		// 						Type: "global",
		// 					},
		// 					Points: []*monitoringpb.Point{
		// 						{
		// 							Interval: &monitoringpb.TimeInterval{
		// 								StartTime: &googlepb.Timestamp{
		// 									Seconds: time.Now().Unix(),
		// 								},
		// 								EndTime: &googlepb.Timestamp{
		// 									Seconds: time.Now().Add(time.Second * time.Duration(1)).Unix(),
		// 								},
		// 							},
		// 							Value: &monitoringpb.TypedValue{
		// 								Value: &monitoringpb.TypedValue_DoubleValue{
		// 									DoubleValue: m.Counter.GetValue(),
		// 								},
		// 							},
		// 						},
		// 					},
		// 					MetricKind: metricpb.MetricDescriptor_CUMULATIVE,
		// 				},
		// 			},
		// 		}); err != nil {
		// 			log.Fatalf("Failed to write time series data: %v", err)
		// 		}
		// 		time.Sleep(1000)
		// 		if err := client.Close(); err != nil {
		// 			log.Fatalf("Failed to close client: %v", err)
		// 		}
		// 	}
		// }

		if v.GetType() == dto.MetricType_HISTOGRAM {
			for _, m := range v.GetMetric() {
				if m == nil {
					continue
				}

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
		}

		if v.GetType() == dto.MetricType_COUNTER {
			for _, m := range v.GetMetric() {
				if m == nil {
					continue
				}
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
		}

		if v.GetType() == dto.MetricType_HISTOGRAM {
			for _, m := range v.GetMetric() {
				if m == nil {
					continue
				}
				for _, l := range m.GetLabel() {
					labels[*l.Name] = *l.Value
				}

				fmt.Println(m.Histogram.GetBucket())

				var bucketCounts []int64
				var bounds []float64
				var count int64
				for _, v := range m.Histogram.GetBucket() {
					fmt.Println(v)
					fmt.Println(v.GetCumulativeCount(), " ", v.GetUpperBound())
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
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	if err := client.CreateTimeSeries(ctx, &monitoringpb.CreateTimeSeriesRequest{
		Name:       fmt.Sprintf("projects/%s", *config["PROJECT_ID"]),
		TimeSeries: ts,
	}); err != nil {
		log.Fatalf("Failed to write time series data: %v", err)
	}
	if err := client.Close(); err != nil {
		log.Fatalf("Failed to close client: %v", err)
	}
	fmt.Printf("Done writing time series data.\n")
}
