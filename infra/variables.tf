variable "rapidapi_key" {
  description = "RapidAPI key used by the Lambda function to fetch schedule data"
  type        = string
  sensitive   = true
}

variable "github_repository" {
  description = "GitHub repository (owner/name) allowed to assume the deployment role via OIDC"
  type        = string
}
