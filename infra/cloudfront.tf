resource "aws_cloudfront_origin_access_control" "calendar" {
  name                              = "calendar-v2"
  description                       = "Allows CloudFront to reach the bucket"
  origin_access_control_origin_type = "s3"
  signing_behavior                  = "always"
  signing_protocol                  = "sigv4"
}

resource "aws_cloudfront_distribution" "calendar" {
  enabled             = true
  is_ipv6_enabled     = true
  default_root_object = "index.html"
  http_version        = "http2"
  price_class         = "PriceClass_100"

  origin {
    origin_id                = "origin1"
    domain_name              = aws_s3_bucket.calendar.bucket_regional_domain_name
    origin_access_control_id = aws_cloudfront_origin_access_control.calendar.id
  }

  default_cache_behavior {
    target_origin_id       = "origin1"
    viewer_protocol_policy = "redirect-to-https"
    allowed_methods        = ["GET", "HEAD"]
    cached_methods         = ["GET", "HEAD"]
    compress               = true
    min_ttl                = 0
    default_ttl            = 86400
    max_ttl                = 31536000

    forwarded_values {
      query_string = false
      cookies {
        forward = "none"
      }
    }
  }

  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }

  viewer_certificate {
    cloudfront_default_certificate = true
    minimum_protocol_version       = "TLSv1"
  }
}
