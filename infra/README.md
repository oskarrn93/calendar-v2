# Infrastructure

Terraform configuration for the calendar-v2 AWS resources (S3 bucket, CloudFront distribution, ECR repository, Lambda function, EventBridge schedule, and the GitHub Actions OIDC deployment role).

State is stored remotely in the `oskarrosen-terraform` S3 bucket (`eu-north-1`), with native S3 locking.

Terraform version is pinned via [`.terraform-version`](.terraform-version) for use with [tfenv](https://github.com/tfutils/tfenv).

## Setup

Copy `terraform.tfvars.example` to `terraform.tfvars` and fill in the values (this file is gitignored since it holds the `RAPIDAPI_KEY` secret).

```sh
cp terraform.tfvars.example terraform.tfvars
```

## Useful commands

* `terraform init` initialize the backend and providers
* `terraform plan` preview changes
* `terraform apply` apply changes
