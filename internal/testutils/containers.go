package testutils

import (
    "context"
    "fmt"
    "testing"
    "time"

    "github.com/stretchr/testify/require"

    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/wait"
)

func TestWithPostgres(t *testing.T) testcontainers.Container {
    ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
    defer cancel()
    req := testcontainers.ContainerRequest{
        Image: "postgres:15",
        ExposedPorts: []string{"5432/tcp"},
        WaitingFor: wait.ForListeningPort("5432/tcp"),
    }
    postgres, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
        ContainerRequest: req,
        Started: true,
    })
    require.NoError(t, err)
    return postgres
}

var PostgresContainerConnectionString string = "postgres://postgres:postgres@127.0.0.1:%d/postgres"

func CloseContainer(cont testcontainers.Container) {
    ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
    defer cancel()
    if err := cont.Terminate(ctx); err != nil {
        panic(fmt.Errorf("failed to terminate container %w", err))
    }
}
