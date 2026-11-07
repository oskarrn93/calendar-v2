resource "aws_sns_topic" "alerts" {
  name = "calendar-v2-alerts"
}

# AWS emails a confirmation link that must be accepted before alerts arrive.
resource "aws_sns_topic_subscription" "alerts_email" {
  topic_arn = aws_sns_topic.alerts.arn
  protocol  = "email"
  endpoint  = var.alert_email
}

# Errors also counts timeouts, so this fires whenever a daily run fails.
resource "aws_cloudwatch_metric_alarm" "lambda_errors" {
  alarm_name          = "calendar-v2-lambda-errors"
  alarm_description   = "The calendar Lambda failed; calendars may be stale."
  namespace           = "AWS/Lambda"
  metric_name         = "Errors"
  dimensions          = { FunctionName = aws_lambda_function.calendar.function_name }
  statistic           = "Sum"
  period              = 3600
  evaluation_periods  = 1
  threshold           = 1
  comparison_operator = "GreaterThanOrEqualToThreshold"
  treat_missing_data  = "notBreaching"
  alarm_actions       = [aws_sns_topic.alerts.arn]
}
