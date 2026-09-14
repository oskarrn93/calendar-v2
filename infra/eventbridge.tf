resource "aws_cloudwatch_event_rule" "calendar_schedule" {
  name                = "calendar-v2-calendarv2bucketlambdaschedule08649407-eyhzsTQyr1QR"
  schedule_expression = "cron(0 6 * * ? *)"
  state               = "ENABLED"
}

resource "aws_cloudwatch_event_target" "calendar_schedule" {
  rule      = aws_cloudwatch_event_rule.calendar_schedule.name
  target_id = "Target0"
  arn       = aws_lambda_function.calendar.arn
}

resource "aws_lambda_permission" "allow_eventbridge" {
  statement_id  = "calendar-v2-calendarv2bucketlambdascheduleAllowEventRulecalendarv2calendarv2bucketlambd-VjRCkRcPfUGA"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.calendar.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.calendar_schedule.arn
}
