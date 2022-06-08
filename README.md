# Terraform module for pushing Prometheus metrics to Google Cloud Monitoring 

Module is based on Cloud Function called by Cloud Scheduler. Each call says which endpoint should be pushed to
Google Cloud Monitoring custom metrics.

For further setup, use variable `scrape_jobs`, e.g.:

```hcl
module "metrics_push" {
  source     = "../"  # use correct module path
  project_id = var.project
  scrape_jobs = {
    testing : {
      schedule : "* * * * *"
      endpoint : "https://example.com/metrics"
    }
  }
  region = var.region
}
```

Each key in `scrape_jobs` can contain following keys:

 * `schedule` (required) - cronline
 * `endpoint` (required) - url to metrics endpoint
 * `description` - used for scheduler description
 * `time_zone` - time zone, default "Europe/Prague"

Kuddos for https://github.com/google/go-metrics-stackdriver/blob/main/stackdriver.go on how to handle histograms.

Currently supported Prometheus metrics are:

 * Gauge,
 * Counter,
 * Histogram.

Metrics are kept in Google Cloud Monitoring under following path: `custom.googleapis.com/<key>/<metric>`, e.g. `custom.googleapis.com/testing/app_info`
where `testing` was given as a key in `scrape_jobs` variable.

## Before you do anything in this module

Install pre-commit hooks by running following commands:

```shell script
brew install pre-commit terraform-docs
pre-commit install
```

<!-- BEGINNING OF PRE-COMMIT-TERRAFORM DOCS HOOK -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | >= 0.13 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_archive"></a> [archive](#provider\_archive) | n/a |
| <a name="provider_google"></a> [google](#provider\_google) | n/a |
| <a name="provider_random"></a> [random](#provider\_random) | n/a |

## Modules

No modules.

## Resources

| Name | Type |
|------|------|
| [google_cloud_scheduler_job.scrape_job](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/cloud_scheduler_job) | resource |
| [google_cloudfunctions_function.function](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/cloudfunctions_function) | resource |
| [google_cloudfunctions_function_iam_member.invoker](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/cloudfunctions_function_iam_member) | resource |
| [google_project_iam_member.monitoring](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/project_iam_member) | resource |
| [google_service_account.sa](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/service_account) | resource |
| [google_storage_bucket.bucket](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/storage_bucket) | resource |
| [google_storage_bucket_object.archive](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/storage_bucket_object) | resource |
| [random_string.random](https://registry.terraform.io/providers/hashicorp/random/latest/docs/resources/string) | resource |
| [archive_file.dotfiles](https://registry.terraform.io/providers/hashicorp/archive/latest/docs/data-sources/file) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_project_id"></a> [project\_id](#input\_project\_id) | Project ID | `string` | n/a | yes |
| <a name="input_region"></a> [region](#input\_region) | GCP region | `string` | n/a | yes |
| <a name="input_scrape_jobs"></a> [scrape\_jobs](#input\_scrape\_jobs) | Metrics scraping setup, each item needs key schedule (e.g. * * * * *) and endpoint (https://example.com/metrics) | `map(map(string))` | n/a | yes |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_sa"></a> [sa](#output\_sa) | Service account used for Cloud Function runtime email |
<!-- END OF PRE-COMMIT-TERRAFORM DOCS HOOK -->
