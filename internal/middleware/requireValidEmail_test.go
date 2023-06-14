package middleware

import (
	"log"
	"net/http/httptest"
	"testing"

	api "github.com/StampWallet/backend/internal/api/models"
	. "github.com/StampWallet/backend/internal/testutils"
	"github.com/stretchr/testify/require"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

func createRequireValidEmailMiddleware(ctrl *gomock.Controller) *RequireValidEmailMiddleware {
	return &RequireValidEmailMiddleware{
		logger: log.Default(),
	}
}

func TestRequireValidEmailMiddlewareOk(t *testing.T) {
	testUser := GetDefaultUser()
	testUser.EmailVerified = true

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	context := NewTestContextBuilder(w).
		SetDefaultUrl().
		SetEndpoint("/user/cards").
		SetMethod("GET").
		SetUser(testUser).
		SetDefaultToken().
		SetHeader("Accept", "application/json").
		Context

	ctrl := gomock.NewController(t)
	handler := createRequireValidEmailMiddleware(ctrl)

	handler.Handle(context)
	require.Equalf(t, 0, w.Body.Len(), "Response is not nil")
	require.Falsef(t, context.IsAborted(), "Context was aborted despite user being verified")
	// Q: Is this the best way to test context was no-op'ed?
}

func TestRequireValidEmailMiddlewareNok(t *testing.T) {
	testUser := GetDefaultUser()
	testUser.EmailVerified = false

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	context := NewTestContextBuilder(w).
		SetDefaultUrl().
		SetEndpoint("/user/cards").
		SetMethod("GET").
		SetUser(testUser).
		SetDefaultToken().
		SetHeader("Accept", "application/json").
		Context

	ctrl := gomock.NewController(t)
	handler := createRequireValidEmailMiddleware(ctrl)

	handler.Handle(context)

	respBodyExpected := api.DefaultResponse{Status: api.FORBIDDEN, Message: "EMAIL_NOT_VERIFIED"}
	respBody, respCode, respParseErr := ExtractResponse[api.DefaultResponse](w)

	require.Nilf(t, respParseErr, "Failed to parse JSON response")
	require.Equalf(t, respCode, int(403), "Response returned unexpected status code")
	require.Truef(t, MatchEntities(respBodyExpected, *respBody), "Response returned unexpected body contents")
}
