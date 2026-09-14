terraform {
  backend "s3" {
    bucket       = "oskarrosen-terraform"
    key          = "calendar-v2/terraform.tfstate"
    region       = "eu-north-1"
    encrypt      = true
    use_lockfile = true
  }
}
