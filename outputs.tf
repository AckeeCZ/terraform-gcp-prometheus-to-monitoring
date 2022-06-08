output "sa" {
  description = "Service account used for Cloud Function runtime email"
  value       = google_service_account.sa.email
}
