package api

import (
	"encoding/json"
	"log"
	"net/http/httptest"
	"testing"

	api "github.com/StampWallet/backend/internal/api/models"
	"github.com/StampWallet/backend/internal/database"
	. "github.com/StampWallet/backend/internal/database/accessors/mocks"
	"github.com/StampWallet/backend/internal/managers"
	. "github.com/StampWallet/backend/internal/managers/mocks"
	. "github.com/StampWallet/backend/internal/testutils"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func getUserHandlers(ctrl *gomock.Controller) *UserHandlers {
	commonTransactionManager := NewMockTransactionManager(ctrl)
	commonUserAuthorizedAccessor := NewMockUserAuthorizedAccessor(ctrl)
	commonVirtualCardManager := NewMockVirtualCardManager(ctrl)
	commonLocalCardManager := NewMockLocalCardManager(ctrl)
	return &UserHandlers{
		virtualCardManager: commonVirtualCardManager,
		localCardManager:   commonLocalCardManager,
		businessManager:    NewMockBusinessManager(ctrl),
		transactionManager: commonTransactionManager,
		virtualCardHandlers: &UserVirtualCardHandlers{
			virtualCardManager:    commonVirtualCardManager,
			transactionManager:    commonTransactionManager,
			itemDefinitionManager: NewMockItemDefinitionManager(ctrl),
			userAuthorizedAcessor: commonUserAuthorizedAccessor,
			logger:                log.Default(),
		},
		localCardHandlers: &UserLocalCardHandlers{
			localCardManager:      commonLocalCardManager,
			userAuthorizedAcessor: commonUserAuthorizedAccessor,
			logger:                log.Default(),
		},
		userAuthorizedAcessor: commonUserAuthorizedAccessor,
		logger:                log.Default(),
	}

}

func getVirtualCardHandlers(ctrl *gomock.Controller) *UserVirtualCardHandlers {
	return &UserVirtualCardHandlers{
		virtualCardManager:    NewMockVirtualCardManager(ctrl),
		transactionManager:    NewMockTransactionManager(ctrl),
		itemDefinitionManager: NewMockItemDefinitionManager(ctrl),
		userAuthorizedAcessor: NewMockUserAuthorizedAccessor(ctrl),
		logger:                log.Default(),
	}
}

func getLocalCardHandlers(ctrl *gomock.Controller) *UserLocalCardHandlers {
	return &UserLocalCardHandlers{
		localCardManager:      NewMockLocalCardManager(ctrl),
		userAuthorizedAcessor: NewMockUserAuthorizedAccessor(ctrl),
		logger:                log.Default(),
	}
}

func setupBusinessHandlersPostAccount() (
	w *httptest.ResponseRecorder,
	context *gin.Context,
	testBusinessUser *database.User,
	testBusiness *database.Business,
	testBusinessDetails *managers.BusinessDetails,
	respBodyExpected *api.PostBusinessAccountResponse,
) {
	testBusinessUser = GetDefaultUser()
	testBusiness = GetDefaultBusiness(testBusinessUser)
	testBusinessDetails = &managers.BusinessDetails{
		Name:        testBusiness.Name,
		Description: testBusiness.Description,
		Address:     testBusiness.Address,
		// GPSCoordinates: testBusiness.GPSCoordinates, TODO GPSCoordinates to string
		NIP:       testBusiness.NIP,
		KRS:       testBusiness.KRS,
		REGON:     testBusiness.REGON,
		OwnerName: testBusiness.OwnerName,
	}

	payload := api.PostBusinessAccountRequest{
		Name:    testBusiness.Name,
		Address: testBusiness.Address,
		// GpsCoordinates: testBusiness.GPSCoordinates, TODO GPSCoordinates to string
		Nip:       testBusiness.NIP,
		Krs:       testBusiness.KRS,
		Regon:     testBusiness.REGON,
		OwnerName: testBusiness.OwnerName,
	}
	payloadJson, _ := json.Marshal(payload)

	gin.SetMode(gin.TestMode)
	w = httptest.NewRecorder()

	context = NewTestContextBuilder(w).
		SetDefaultUrl().
		SetEndpoint("/business/account").
		SetUser(testBusinessUser).
		SetMethod("POST").
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetDefaultToken().
		SetBody(payloadJson).
		Context

	respBodyExpected = &api.PostBusinessAccountResponse{
		PublicId:      testBusiness.PublicId,
		BannerImageId: testBusiness.BannerImageId,
		IconImageId:   testBusiness.IconImageId,
	}

	return w, context, testBusinessUser, testBusiness, testBusinessDetails, respBodyExpected
}

func TestUserLocalCardHandlersPostCardOk(t *testing.T) {
	testUser := GetDefaultUser()
	testCard := GetTestLocalCard(nil, testUser)
	testCardDetails := managers.LocalCardDetails{
		Type: testCard.Type,
		Code: testCard.Code,
		Name: testCard.Name,
	}

	payload := api.PostUserLocalCardsRequest{
		Name: testCard.Name,
		Type: testCard.Type,
		Code: testCard.Code,
	}
	payloadJson, _ := json.Marshal(payload)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()

	context := NewTestContextBuilder(w).
		SetDefaultUrl().
		SetEndpoint("/user/cards/local").
		SetUser(testUser).
		SetMethod("POST").
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetDefaultToken().
		SetBody(payloadJson).
		Context

	respBodyExpected := &api.DefaultResponse{Status: api.OK}

	// test env prep
	ctrl := gomock.NewController(t)
	handler := getLocalCardHandlers(ctrl)

	handler.localCardManager.(*MockLocalCardManager).
		EXPECT().
		Create(
			gomock.Eq(testUser),
			gomock.Eq(testCardDetails),
		).
		Return(
			testCard,
			nil,
		)

	handler.postCard(context)

	respBody, respCode, respParseErr := ExtractResponse[api.DefaultResponse](w)

	require.Nilf(t, respParseErr, "Failed to parse JSON response")
	require.Equalf(t, respCode, int(201), "Response returned unexpected status code")
	require.Truef(t, MatchEntities(respBodyExpected, respBody), "Response returned unexpected body contents")
	// TODO: test MatchEntities and gomock.Eq
}

func TestUserLocalCardHandlersDeleteCardOk(t *testing.T) {
	testUser := GetDefaultUser()
	testCard := GetTestLocalCard(nil, testUser)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()

	context := NewTestContextBuilder(w).
		SetDefaultUrl().
		SetEndpoint("/user/cards/local/"+testCard.PublicId).
		SetUser(testUser).
		SetMethod("DELETE").
		SetHeader("Content-Type", "application/json").
		SetDefaultToken().
		Context

	respBodyExpected := &api.DefaultResponse{Status: api.OK}

	// test env prep
	ctrl := gomock.NewController(t)
	handler := getLocalCardHandlers(ctrl)

	handler.userAuthorizedAcessor.(*MockUserAuthorizedAccessor).
		EXPECT().
		Get(
			gomock.Eq(testUser),
			gomock.Eq(database.LocalCard{PublicId: testCard.PublicId}),
		).
		Return(
			testCard,
			nil,
		)

	handler.localCardManager.(*MockLocalCardManager).
		EXPECT().
		Remove(gomock.Eq(testCard)).
		Return(nil)

	handler.deleteCard(context)

	respBody, respCode, respParseErr := ExtractResponse[api.DefaultResponse](w)

	require.Nilf(t, respParseErr, "Failed to parse JSON response")
	require.Equalf(t, respCode, int(200), "Response returned unexpected status code")
	require.Truef(t, MatchEntities(respBodyExpected, respBody), "Response returned unexpected body contents")
	// TODO: test MatchEntities and gomock.Eq
}
