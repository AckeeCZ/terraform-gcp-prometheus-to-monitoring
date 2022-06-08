variable "region" {
  description = "GCP region"
  type        = string
}

variable "project_id" {
  description = "Project ID"
  type        = string
}

variable "scrape_jobs" {
  description = "Metrics scraping setup, each item needs key schedule (e.g. * * * * *) and endpoint (https://example.com/metrics)"
  type        = map(map(string))
}
