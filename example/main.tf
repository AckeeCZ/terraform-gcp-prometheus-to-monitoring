variable "project" {}

variable "region" {
  default = "europe-west3"
}

module "pubsub_to_bq" {
  source     = "../"
  project_id = var.project
  scrape_jobs = {
    testing : {
      schedule : "* * * * *"
      endpoint : "https://example.com/metrics" # change
    }
  }
  region = var.region
}
