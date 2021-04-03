

```
gcloud beta emulators firestore start --host-port=localhost
```

Set `FIRESTORE_EMULATOR_HOST=localhost:8080`

```
curl -v -X DELETE "http://localhost:8080/emulator/v1/projects/test-project/databases/(default)/documents"
```
