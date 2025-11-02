# calendar-v2

This project provides a solution for hosting calendar `.ics` files that allow users to subscribe to NBA and Football game schedules directly in their personal calendars.

A Lambda function runs daily as a scheduled job to generate and upload the `.ics` files to a S3 bucket. These files are then made publicly accessible via a CloudFront distribution, enabling seamless access and subscription to the calendar events.

## Deployment

### Build Image

Images are built when a new tagged version is published (e.g. `v1.2.3`) and is published to AWS ECR.

### Deploy

Manually invoke the [Deploy Github Actions workflow](https://github.com/oskarrn93/calendar-v2/actions/workflows/publish.yaml) and specify the tagged version (e.g. `v1.2.3`) where the image has been published.

### Infrastructure

To deploy infrastructure run the following command from your local machine. This assumes you have setup the AWS Profile and have the expected IAM permissions.

```sh
make deploy
```

## Development

### Tests

Run tests

```sh
make test-ci
```

#### Update snapshots

If the application logic is changed and the test snapshots needs to be update then include `UPDATED_SNAPS=true` when running the tests.

```sh
UPDATE_SNAPS=true make test
```
