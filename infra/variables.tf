variable "rapidapi_key" {
  description = "RapidAPI key used by the Lambda function to fetch schedule data"
  type        = string
  sensitive   = true
  ephemeral   = true
}

variable "rapidapi_key_version" {
  description = "Increment to store a new rapidapi_key value in SSM"
  type        = number
  default     = 1
}

variable "github_repository" {
  description = "GitHub repository (owner/name) allowed to assume the deployment role via OIDC"
  type        = string
}

variable "alert_email" {
  description = "Email address that receives an alert when the Lambda function fails"
  type        = string
}
