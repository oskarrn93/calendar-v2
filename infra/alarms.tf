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

# Catches a schedule that stopped firing, which the errors alarm can't see. 25
# hourly periods always contain the previous daily run.
resource "aws_cloudwatch_metric_alarm" "lambda_not_invoked" {
  alarm_name          = "calendar-v2-lambda-not-invoked"
  alarm_description   = "The calendar Lambda has not run in over a day; calendars are going stale."
  namespace           = "AWS/Lambda"
  metric_name         = "Invocations"
  dimensions          = { FunctionName = aws_lambda_function.calendar.function_name }
  statistic           = "Sum"
  period              = 3600
  evaluation_periods  = 25
  datapoints_to_alarm = 25
  threshold           = 1
  comparison_operator = "LessThanThreshold"
  treat_missing_data  = "breaching"
  alarm_actions       = [aws_sns_topic.alerts.arn]
}

resource "aws_cloudwatch_log_metric_filter" "no_events" {
  name           = "calendar-v2-no-events"
  log_group_name = "/aws/lambda/${aws_lambda_function.calendar.function_name}"
  pattern        = "{ $.event = \"no_events\" }"

  metric_transformation {
    name       = "NoEvents"
    namespace  = "calendar-v2"
    value      = "1"
    dimensions = { job = "$.job" }
  }
}

# Empty runs don't fail the Lambda, so a broken source would otherwise go
# unnoticed. Hourly periods over 7 days hold exactly 7 daily runs; requiring 6
# stops the alarm flapping while the newest run is still being reported. Expect
# one alert per off-season for seasonal jobs.
resource "aws_cloudwatch_metric_alarm" "no_events" {
  for_each = toset(["basketball", "esport", "f1", "football", "nba"])

  alarm_name          = "calendar-v2-no-events-${each.key}"
  alarm_description   = "The ${each.key} job has found no events for about a week. Expected in the off-season; otherwise the API or its schema may have changed."
  namespace           = aws_cloudwatch_log_metric_filter.no_events.metric_transformation[0].namespace
  metric_name         = aws_cloudwatch_log_metric_filter.no_events.metric_transformation[0].name
  dimensions          = { job = each.key }
  statistic           = "Sum"
  period              = 3600
  evaluation_periods  = 168
  datapoints_to_alarm = 6
  threshold           = 1
  comparison_operator = "GreaterThanOrEqualToThreshold"
  treat_missing_data  = "notBreaching"
  alarm_actions       = [aws_sns_topic.alerts.arn]
}
