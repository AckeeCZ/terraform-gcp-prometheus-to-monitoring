from google.cloud import monitoring_v3

descriptor_name="projects/flash-news-development/metricDescriptors/custom.googleapis.com/testing/flask_exporter_info"

client = monitoring_v3.MetricServiceClient()
client.delete_metric_descriptor(name=descriptor_name)
print("Deleted metric descriptor {}.".format(descriptor_name))

# flask_http_request_created
# flask_exporter_info