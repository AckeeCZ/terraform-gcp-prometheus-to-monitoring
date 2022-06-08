resource "random_string" "random" {
  length  = 16
  special = false
}

data "archive_file" "dotfiles" {
  type        = "zip"
  output_path = "${path.module}/dotfiles.zip"


  dynamic "source" {
    for_each = [
      {
        content  = "${path.module}/code/go.mod"
        filename = "go.mod"
      },
      {
        content  = "${path.module}/code/go.sum"
        filename = "go.sum"
      },
      {
        content  = "${path.module}/code/load_vars.go"
        filename = "load_vars.go"
      },
      {
        content  = "${path.module}/code/main.go"
        filename = "main.go"
      },
      {
        content  = "${path.module}/code/parse_mf.go"
        filename = "parse_mf.go"
      },
      {
        content  = "${path.module}/code/submit_to_monitoring.go"
        filename = "submit_to_monitoring.go"
      },
    ]
    content {
      content  = file(source.value["content"])
      filename = source.value["filename"]
    }
  }
}

resource "google_storage_bucket" "bucket" {
  name = "prometheus_to_monitoring_${lower(random_string.random.result)}_deployment_bucket"

  force_destroy = true
  location      = "EUROPE-WEST3"
}

resource "google_storage_bucket_object" "archive" {
  name   = "dotfiles.zip"
  bucket = google_storage_bucket.bucket.name
  source = data.archive_file.dotfiles.output_path
}

resource "google_cloudfunctions_function" "function" {
  name        = "prometheus_to_monitoring_${lower(random_string.random.result)}"
  description = "Prometheus metrics endpoint scraping to Google Cloud Monitoring"
  runtime     = "go116"
  region      = var.region

  trigger_http          = true
  available_memory_mb   = 128
  source_archive_bucket = google_storage_bucket.bucket.name
  source_archive_object = google_storage_bucket_object.archive.name
  entry_point           = "Consume"
  service_account_email = google_service_account.sa.email

  environment_variables = {
    PROJECT_ID = var.project_id
  }
  depends_on = [google_storage_bucket_object.archive]
}

resource "google_cloud_scheduler_job" "scrape_job" {
  for_each         = var.scrape_jobs
  name             = each.key
  description      = lookup(each.value, "description", "")
  schedule         = each.value["schedule"]
  time_zone        = lookup(each.value, "time_zone", "Europe/Prague")
  attempt_deadline = "60s"
  region           = var.region

  http_target {
    http_method = "POST"
    uri         = google_cloudfunctions_function.function.https_trigger_url
    body        = base64encode("{\"endpoint\":\"${each.value["endpoint"]}\", \"service\": \"${each.key}\"}")

    oidc_token {
      service_account_email = google_service_account.sa.email
    }
  }
}
