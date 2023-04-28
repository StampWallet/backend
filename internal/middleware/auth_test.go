package middleware

import (
	"log"
	"net/http/httptest"
	"testing"

	. "github.com/StampWallet/backend/internal/database"
	. "github.com/StampWallet/backend/internal/services/mocks"
	. "github.com/StampWallet/backend/internal/testutils"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func getAuthMiddleware(ctrl *gomock.Controller) *AuthMiddleware {
	return &AuthMiddleware{
		logger:       log.Default(),
		tokenService: NewMockTokenService(ctrl),
	}
}

func TestHandleOk(t *testing.T) {
	// data prep
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()

	context := NewTestContextBuilder(w).
		SetDefaultUrl().
		SetEndpoint("/user/cards").
		SetDefaultUser().
		SetMethod("GET").
		SetHeader("Authorization", "Bearer ZWVnaDhhZWg4bGVpbDJhaXBlaW5nZWViNWFpU2hlaGUK").
		SetHeader("Accept", "application/json").
		Context

	// test env prep
	ctrl := gomock.NewController(t)
	authMiddleware := getAuthMiddleware(ctrl)

	// TODO: mock accessor to token

	expectedTokenId := "0"

	authMiddleware.tokenService.(*MockTokenService).
		EXPECT().
		Check(
			gomock.Eq(expectedTokenId),
			gomock.Eq("ZWVnaDhhZWg4bGVpbDJhaXBlaW5nZWViNWFpU2hlaGUK"),
		).
		Return(
			GetDefaultUser(),
			&Token{
				OwnerId:   0,
				TokenId:   expectedTokenId,
				TokenHash: "ZWVnaDhhZWg4bGVpbDJhaXBlaW5nZWViNWFpU2hlaGUK",
				// Expires: _, TODO: needed?
				TokenPurpose: TokenPurposeSession,
				Used:         false,
				Recalled:     false,
			},
			nil,
		)

	authMiddleware.Handle(context)

	user, _ := context.Get("user")
	require.Truef(t, MatchEntities(GetDefaultUser(), user), "User entity provided by authenticator failed validation")
}
