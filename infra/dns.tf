locals {
  domain_name = "calendar.oskarrosen.io"
}

data "aws_route53_zone" "root" {
  name = "oskarrosen.io"
}

resource "aws_acm_certificate" "calendar" {
  # CloudFront only accepts certificates from us-east-1.
  provider          = aws.us_east_1
  domain_name       = local.domain_name
  validation_method = "DNS"

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_route53_record" "calendar_validation" {
  for_each = {
    for dvo in aws_acm_certificate.calendar.domain_validation_options : dvo.domain_name => {
      name   = dvo.resource_record_name
      type   = dvo.resource_record_type
      record = dvo.resource_record_value
    }
  }

  zone_id         = data.aws_route53_zone.root.zone_id
  name            = each.value.name
  type            = each.value.type
  records         = [each.value.record]
  ttl             = 300
  allow_overwrite = true
}

resource "aws_acm_certificate_validation" "calendar" {
  provider                = aws.us_east_1
  certificate_arn         = aws_acm_certificate.calendar.arn
  validation_record_fqdns = [for record in aws_route53_record.calendar_validation : record.fqdn]
}

resource "aws_route53_record" "calendar" {
  for_each = toset(["A", "AAAA"])

  zone_id = data.aws_route53_zone.root.zone_id
  name    = local.domain_name
  type    = each.value

  alias {
    name                   = aws_cloudfront_distribution.calendar.domain_name
    zone_id                = aws_cloudfront_distribution.calendar.hosted_zone_id
    evaluate_target_health = false
  }
}

# Previously pointed at an API Gateway custom domain created outside Terraform.
import {
  to = aws_route53_record.calendar["A"]
  id = "${data.aws_route53_zone.root.zone_id}_${local.domain_name}_A"
}
