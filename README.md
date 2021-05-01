# Flake Reporter

> Left uncontrolled, non-deterministic tests can completely destroy the value of an automated regression suite (Martin Fowler)

Flake Reporter is a tool for collating test results so that non-deterministic
(AKA Flaky) tests can be discovered and dealt with.

What Flake Reporter currently does:

- Collect test results in JUnit format
- Report on the flakiness of a given test across the last 100 runs
- Expose the flakiness of tests as Prometheus metrics

What Flake Reporter aims to do, in the future:

- Present results in a more human readable format, so that flaky tests can be
  easily identified
- Automatically open GitHub issues to track flaky tests

What Flake Reporter will never do:

- Prevent you from adding flaky tests
- Fix your flaky tests
- Automatically restart flaky tests

## Design

Flake Reporter uses Google Firestore to track the last 100 runs of a given test
within a test suite for a particular branch.

The success/failure rate for the 100 runs is then converted to a percentage to
represent the flakiness of a test.

## Deployment

### Helm Chart

You will require a GCP project with firestore enabled.

The best way to deploy Flake Reporter is using the Helm chart.

```sh
helm repo add flake-reporter https://playerdata.github.io/flake-reporter
helm repo update
helm install flake-reporter/flake-reporter
```

See [values.yaml](./charts/flake-reporter/values.yaml) for configuration.

### Environment Variables

| Name                   | Value                                            |
| ---------------------- | ------------------------------------------------ |
| `FIRESTORE_PROJECT_ID` | The GCP/Firebase project ID to use for Firestore |

### CI Integration

Once you have Flake Reporter running, you will need to configure your CI to report
test result to Flake Reporter.

The simplest way to do this is using `curl`:

```
curl -X POST -F project=test-project -F branch=main -F file="@path-to-junit.xml" https://flake-reporter-address:9090/recv/junit
```

Note the form fields:

| Name    | Field                                                                                                         |
| ------- | ------------------------------------------------------------------------------------------------------------- |
| project | An arbitrary string used to group test suites together. This could be, for example, your code repository name |
| branch  | The branch that is currently being built                                                                      |

## Development

### Running the Tests

First, start a firestore emulator:

```
gcloud beta emulators firestore start --host-port=localhost
```

Then, run the unit tests

```
FIRESTORE_EMULATOR_HOST=localhost:8080 go test ./...
```

### Running Locally

First, start a firestore emulator:

```
gcloud beta emulators firestore start --host-port=localhost
```

Now, run the app:

```
FIRESTORE_EMULATOR_HOST=localhost:8080 go run main.go
```

Submit a junit report:

```
curl -X POST -F project=test-project -F branch=main -F file="@fixtures/junit-success.xml" http://localhost:9090/recv/junit
```

And look at the results:

```
curl "http://localhost:9090/summary/test-project/Session%20Tags/Session%20Tags%20can%20create%20a%20session%20tag"
```
