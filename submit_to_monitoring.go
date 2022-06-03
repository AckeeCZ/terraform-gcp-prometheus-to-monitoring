package main

import (
	"context"
	"fmt"
	"log"
	"time"

	monitoring "cloud.google.com/go/monitoring/apiv3"
	googlepb "github.com/golang/protobuf/ptypes/timestamp"
	dto "github.com/prometheus/client_model/go"
	metricpb "google.golang.org/genproto/googleapis/api/metric"
	monitoredrespb "google.golang.org/genproto/googleapis/api/monitoredres"
	monitoringpb "google.golang.org/genproto/googleapis/monitoring/v3"
)

func sendToGCPMonitoring(mf map[string]*dto.MetricFamily, config map[string]*string) {
	for metric, v := range mf {
		fmt.Println("KEY: ", metric)
		fmt.Println("VAL: ", v)

		fmt.Println(v.GetType())

		if v.GetType() == dto.MetricType_GAUGE {
			dataPoints := []*monitoringpb.Point{}
			for _, m := range v.GetMetric() {
				if m == nil {
					continue
				}
				labels := make(map[string]string)
				for _, l := range m.GetLabel() {
					labels[*l.Name] = *l.Value
				}
				dataPoints = append(dataPoints)
				ctx := context.Background()

				client, err := monitoring.NewMetricClient(ctx)
				if err != nil {
					log.Fatalf("Failed to create client: %v", err)
				}
				if err := client.CreateTimeSeries(ctx, &monitoringpb.CreateTimeSeriesRequest{
					Name: fmt.Sprintf("projects/%s", *config["PROJECT_ID"]),
					TimeSeries: []*monitoringpb.TimeSeries{
						{
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
											Seconds: time.Now().Unix(),
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
						},
					},
				}); err != nil {
					log.Fatalf("Failed to write time series data: %v", err)
				}
				if err := client.Close(); err != nil {
					log.Fatalf("Failed to close client: %v", err)
				}
			}
		}

		// if v.GetType() == dto.MetricType_COUNTER {
		// 	dataPoints := []*monitoringpb.Point{}
		// 	labels := make(map[string]string)
		// 	for _, m := range v.GetMetric() {
		// 		if m == nil {
		// 			continue
		// 		}
		// 		for _, l := range m.GetLabel() {
		// 			labels[*l.Name] = *l.Value
		// 		}
		// 		dataPoints = append(dataPoints, &monitoringpb.Point{
		// 			Interval: &monitoringpb.TimeInterval{
		// 				EndTime: &googlepb.Timestamp{
		// 					Seconds: time.Now().Unix(),
		// 				},
		// 			},
		// 			Value: &monitoringpb.TypedValue{
		// 				Value: &monitoringpb.TypedValue_DoubleValue{
		// 					DoubleValue: m.Counter.GetValue(),
		// 				},
		// 			},
		// 		})
		// 	}
		// 	ctx := context.Background()

		// 	client, err := monitoring.NewMetricClient(ctx)
		// 	if err != nil {
		// 		log.Fatalf("Failed to create client: %v", err)
		// 	}
		// 	if err := client.CreateTimeSeries(ctx, &monitoringpb.CreateTimeSeriesRequest{
		// 		Name: fmt.Sprintf("projects/%s", *config["PROJECT_ID"]),
		// 		TimeSeries: []*monitoringpb.TimeSeries{
		// 			{
		// 				Metric: &metricpb.Metric{
		// 					Type:   fmt.Sprintf("custom.googleapis.com/%s/%s", *config["SERVICE"], *v.Name),
		// 					Labels: labels,
		// 				},
		// 				Resource: &monitoredrespb.MonitoredResource{
		// 					Type:   "global",
		// 					Labels: labels,
		// 				},
		// 				Points:     dataPoints,
		// 				MetricKind: metricpb.MetricDescriptor_CUMULATIVE,
		// 			},
		// 		},
		// 	}); err != nil {
		// 		log.Fatalf("Failed to write time series data: %v", err)
		// 	}
		// 	if err := client.Close(); err != nil {
		// 		log.Fatalf("Failed to close client: %v", err)
		// 	}
		// }

		// if v.GetType() == dto.MetricType_COUNTER {
		// 	for _, m := range v.GetMetric() {
		// 		if m == nil {
		// 			continue
		// 		}
		// 		labels := make(map[string]string)
		// 		config["project_id"] = config["PROJECT_ID"]
		// 		for _, l := range m.GetLabel() {
		// 			config[*l.Name] = l.Value
		// 		}
		// 		ctx := context.Background()

		// 		client, err := monitoring.NewMetricClient(ctx)
		// 		if err != nil {
		// 			log.Fatalf("Failed to create client: %v", err)
		// 		}

		// 		fmt.Println(m.Counter.GetValue())
		// 		dataPoint := &monitoringpb.Point{
		// 			Interval: &monitoringpb.TimeInterval{
		// 				EndTime: &googlepb.Timestamp{
		// 					Seconds: time.Now().Unix(),
		// 				},
		// 			},
		// 			Value: &monitoringpb.TypedValue{
		// 				Value: &monitoringpb.TypedValue_DoubleValue{
		// 					DoubleValue: m.Counter.GetValue(),
		// 				},
		// 			},
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
		// 						Type:   "global",
		// 						Labels: labels,
		// 					},
		// 					Points: []*monitoringpb.Point{
		// 						dataPoint,
		// 					},
		// 					MetricKind: metricpb.MetricDescriptor_CUMULATIVE,
		// 				},
		// 			},
		// 		}); err != nil {
		// 			log.Fatalf("Failed to write time series data: %v", err)
		// 		}
		// 		if err := client.Close(); err != nil {
		// 			log.Fatalf("Failed to close client: %v", err)
		// 		}
		// 	}
		// }
		fmt.Printf("Done writing time series data.\n")
	}
}
