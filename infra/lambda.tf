resource "aws_iam_role" "lambda" {
  name = "calendar-v2-calendarv2bucketlambdaServiceRoleA6AC3D-rtvNQ3WMvDJg"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect    = "Allow"
      Principal = { Service = "lambda.amazonaws.com" }
      Action    = "sts:AssumeRole"
    }]
  })
}

resource "aws_iam_role_policy_attachment" "lambda_basic_execution" {
  role       = aws_iam_role.lambda.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_role_policy" "lambda_default" {
  name = "calendarv2bucketlambdaServiceRoleDefaultPolicy12ADA4F0"
  role = aws_iam_role.lambda.name

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "s3:GetObject*",
          "s3:GetBucket*",
          "s3:List*",
          "s3:DeleteObject*",
          "s3:PutObject",
          "s3:PutObjectLegalHold",
          "s3:PutObjectRetention",
          "s3:PutObjectTagging",
          "s3:PutObjectVersionTagging",
          "s3:Abort*",
        ]
        Resource = [aws_s3_bucket.calendar.arn, "${aws_s3_bucket.calendar.arn}/*"]
      },
      {
        Effect   = "Allow"
        Action   = ["ecr:BatchCheckLayerAvailability", "ecr:GetDownloadUrlForLayer", "ecr:BatchGetImage"]
        Resource = aws_ecr_repository.calendar.arn
      },
      {
        Effect   = "Allow"
        Action   = "ecr:GetAuthorizationToken"
        Resource = "*"
      },
    ]
  })
}

resource "aws_lambda_function" "calendar" {
  function_name = "calendar-v2-calendarv2bucketlambda0711F327-EQwreQKMICen"
  role          = aws_iam_role.lambda.arn
  package_type  = "Image"
  image_uri     = "${aws_ecr_repository.calendar.repository_url}:v1.7.0"
  architectures = ["x86_64"]
  timeout       = 30
  memory_size   = 128

  environment {
    variables = {
      RAPIDAPI_KEY   = var.rapidapi_key
      S3_BUCKET_NAME = aws_s3_bucket.calendar.id
    }
  }

  # Deploys update the running image out-of-band via CI (.github/actions/deploy),
  # so Terraform must not fight that with a stale/"latest"-pinned image_uri.
  lifecycle {
    ignore_changes = [image_uri]
  }
}
