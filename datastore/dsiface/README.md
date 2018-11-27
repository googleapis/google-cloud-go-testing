Datastore
=========

## Run Tests

### Run unit tests

The unit tests run locally against an example fake and do not require any
configuration.

```
go test --short
```

### Run all tests

The integration tests will reach out and store, read and then delete a value
from a datastore. You have to provide the `DATASTORE_PROJECT_ID` environment
variable and the client has to be able to [find credentials][find-creds].

For example:

```
DATASTORE_PROJECT_ID=<proj-id> GOOGLE_APPLICATION_CREDENTIALS=</path/to/creds.json> go test
```

[find-creds]: https://godoc.org/cloud.google.com/go#hdr-Authentication_and_Authorization
