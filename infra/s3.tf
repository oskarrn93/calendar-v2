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
        Effect = "Allow"
        Principal = {
          AWS = "arn:aws:iam::cloudfront:user/CloudFront Origin Access Identity ${aws_cloudfront_origin_access_identity.calendar.id}"
        }
        Action   = ["s3:GetObject*", "s3:GetBucket*", "s3:List*"]
        Resource = [aws_s3_bucket.calendar.arn, "${aws_s3_bucket.calendar.arn}/*"]
      },
    ]
  })
}
