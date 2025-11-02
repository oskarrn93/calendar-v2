# calendar-v2

This project provides a solution for hosting calendar `.ics` files that allow users to subscribe to NBA and Football game schedules directly in their personal calendars.

A Lambda function runs daily as a scheduled job to generate and upload the `.ics` files to a S3 bucket. These files are then made publicly accessible via a CloudFront distribution, enabling seamless access and subscription to the calendar events.
