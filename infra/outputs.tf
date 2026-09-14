output "cloudfront_domain_name" {
  value = aws_cloudfront_distribution.calendar.domain_name
}

output "ecr_repository_url" {
  value = aws_ecr_repository.calendar.repository_url
}

output "lambda_function_name" {
  value = aws_lambda_function.calendar.function_name
}

output "github_actions_role_arn" {
  value = aws_iam_role.github_actions.arn
}
