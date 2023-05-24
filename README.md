# StampWallet

Backend serivce for StampWallet app.
Requires [Go 1.20](https://go.dev/doc/install)

## Regenerate API models

1. Install [openapi-generator](https://openapi-generator.tech/docs/installation)
2. Run scripts/generateApi.sh from the directory where openapi-generator jar is stored

## Build

1. Install [Go 1.20](https://go.dev/doc/install)
2. `make` or `make bin` to just build the binary and skip all tests

## Test

Tests require a working Postgres database with PostGIS extensions. Two environment variables are required to configure the tests:

* TEST_DATABASE_URL - Database URL, for example `postgres://postgres@localhost/stampwallet`
* TEST_DATABASE_NAME - Test database name, for example `stampwallet`. NOTE: This database will be dropped and recreated repeatedly. All data from this database *WILL* be lost.

To set up a Postgres database on docker with PostGIS (example): `sudo docker run -e POSTGRES_HOST_AUTH_METHOD=trust -p 5432:5432 -d postgis/postgis`

1. `go test -v ./...` or `test . -run "^TestBusiness.*$"` to only run tests that match a string

