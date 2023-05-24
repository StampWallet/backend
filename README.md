# StampWallet

Backend serivce for StampWallet app.
Requires [Go 1.20](https://go.dev/doc/install)

## Regenerate API models

1. Install [openapi-generator](https://openapi-generator.tech/docs/installation)
2. Run scripts/generateApi.sh from the directory where openapi-generator jar is stored

## Build

1. Install [Go 1.20](https://go.dev/doc/install)
2. `make` or `make bin` to just build the binary and skip all tests
3. Run `./stampWalletServer example-config` to generate example config

## Test

Tests require a working Postgres database with PostGIS extensions. Two environment variables are required to configure the tests:

* TEST_DATABASE_URL - Database URL, for example `postgres://postgres@localhost/stampwallet`
* TEST_DATABASE_NAME - Test database name, for example `stampwallet`. NOTE: This database will be dropped and recreated repeatedly. All data from this database *WILL* be lost.

To set up a Postgres database on docker with PostGIS (example): `sudo docker run -e POSTGRES_HOST_AUTH_METHOD=trust -p 5432:5432 -d postgis/postgis`

1. `go test -v ./...` or `test . -run "^TestBusiness.*$"` to only run tests that match a string

## Configuration 

`example-config` subcommand will generate an example configuration file. 

```yaml
BackendDomain: localhost                                        # ?
DatabaseUrl: 'postgresql://user:password@localhost/db'          # Postgres database URL 
EmailVerificationFrontendURL: localhost                         # URL of email verification website
ServerUrl: localhost:8080                                       # IP and port the server will listen on (TODO: change name)
SmtpConfig:
    ServerHostname: smtp.example.com                            # SMTP Server hostname
    ServerPort: 465                                             # SMTP Server port
    Username: test@example.com                                  # SMTP auth username
    Password: 'password'                                        # SMTP auth password
    SenderEmail: test@example.com                               # Email Address to put in "from" field
StoragePath: /tmp/                                              # Where to store uploaded files
```

## Docker image 

To build Docker image, run `sudo docker buildx build --progress=plain .`. [Buildkit](https://docs.docker.com/build/buildkit/) might be required.

Development docker-compose file is also provided. It will set up both the backend, and the database.
* `docker-compose -f docker-compose.yaml build` - build image
* `docker-compose -f docker-compose.yaml start` - start the whole dev stack

To configure the server, copy `docker-compose.dev.example.yaml` to `docker-compose.dev.yaml`, change values in `docker-compose.dev.yml` and add `-f docker-compose.dev.yml` to next docker-compose runs.

`sudo docker-compose -f docker-compose.yaml -f docker-compose.dev.yaml up`

