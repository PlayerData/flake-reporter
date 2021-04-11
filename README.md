# Flake Reporter

A very early tool for tracking test flakiness.

## Configuration

### Environment Variables

| Name                 | Value                                            |
| -------------------- | ------------------------------------------------ |
| `FIRESTORE_PROJECT_ID` | The GCP/Firebase project ID to use for Firestore |

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
