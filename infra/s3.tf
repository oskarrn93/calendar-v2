resource "aws_s3_bucket" "calendar" {
  bucket = "calendar-oskarrosen-io"
}

resource "aws_s3_bucket_public_access_block" "calendar" {
  bucket = aws_s3_bucket.calendar.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

resource "aws_s3_bucket_policy" "calendar" {
  bucket = aws_s3_bucket.calendar.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect    = "Allow"
        Principal = { Service = "cloudfront.amazonaws.com" }
        Action    = "s3:GetObject"
        Resource  = "${aws_s3_bucket.calendar.arn}/*"
        Condition = {
          StringEquals = { "AWS:SourceArn" = aws_cloudfront_distribution.calendar.arn }
        }
      },
    ]
  })
}

# Lets a bad upload be rolled back to the previous calendar.
resource "aws_s3_bucket_versioning" "calendar" {
  bucket = aws_s3_bucket.calendar.id

  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_lifecycle_configuration" "calendar" {
  bucket = aws_s3_bucket.calendar.id

  rule {
    id     = "expire-old-versions"
    status = "Enabled"

    filter {}

    noncurrent_version_expiration {
      noncurrent_days = 30
    }
  }

  depends_on = [aws_s3_bucket_versioning.calendar]
}
