/*
StampWallet API Server

Testing TransactionsApiService

*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech);

package openapi

import (
	"context"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_openapi_TransactionsApiService(t *testing.T) {

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)

	t.Run("Test TransactionsApiService FinishTransaction", func(t *testing.T) {

		t.Skip("skip test") // remove to run test

		var transactionCode string

		resp, httpRes, err := apiClient.TransactionsApi.FinishTransaction(context.Background(), transactionCode).Execute()

		require.Nil(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, 200, httpRes.StatusCode)

	})

	t.Run("Test TransactionsApiService GetTransactionDetails", func(t *testing.T) {

		t.Skip("skip test") // remove to run test

		var transactionCode string

		resp, httpRes, err := apiClient.TransactionsApi.GetTransactionDetails(context.Background(), transactionCode).Execute()

		require.Nil(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, 200, httpRes.StatusCode)

	})

}
