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

func TestUserVirtualCardHandlersPostCard(t *testing.T) {
	testUser := GetDefaultUser()
	testBusinessUser := GetDefaultUser()
	testBusiness := GetDefaultBusiness(testBusinessUser)
	testCard := GetTestVirtualCard(nil, testUser, testBusiness)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()

	context := NewTestContextBuilder(w).
		SetDefaultUrl().
		SetEndpoint("/user/cards/virtual/"+testBusiness.PublicId).
		SetUser(testUser).
		SetMethod("POST").
		SetHeader("Content-Type", "application/json").
		SetDefaultToken().
		Context

	respBodyExpected := &api.DefaultResponse{Status: api.OK}

	// test env prep
	ctrl := gomock.NewController(t)
	handler := getVirtualCardHandlers(ctrl)

	handler.virtualCardManager.(*MockVirtualCardManager).
		EXPECT().
		Create(
			gomock.Eq(testUser),
			gomock.Eq(testBusiness),
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

func TestUserVirtualCardHandlersDeleteCard(t *testing.T) {
	testUser := GetDefaultUser()
	testBusinessUser := GetDefaultUser()
	testBusiness := GetDefaultBusiness(testBusinessUser)
	testCard := GetTestVirtualCard(nil, testUser, testBusiness)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()

	context := NewTestContextBuilder(w).
		SetDefaultUrl().
		SetEndpoint("/user/cards/virtual/"+testBusiness.PublicId).
		SetUser(testUser).
		SetMethod("DELETE").
		SetHeader("Content-Type", "application/json").
		SetDefaultToken().
		Context

	respBodyExpected := &api.DefaultResponse{Status: api.OK}

	// test env prep
	ctrl := gomock.NewController(t)
	handler := getVirtualCardHandlers(ctrl)

	handler.userAuthorizedAcessor.(*MockUserAuthorizedAccessor).
		EXPECT().
		Get(
			gomock.Eq(testUser),
			gomock.Eq(database.VirtualCard{PublicId: testCard.PublicId}),
		).
		Return(
			testCard,
			nil,
		)

	handler.virtualCardManager.(*MockVirtualCardManager).
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

func TestUserVirtualCardHandlersGetCard(t *testing.T) {
	testUser := GetDefaultUser()
	testBusinessUser := GetDefaultUser()
	testBusiness := GetDefaultBusiness(testBusinessUser)
	testCard := GetTestVirtualCard(nil, testUser, testBusiness)
	testItemDef := GetDefaultItem(testBusiness)
	testOwnedItem := GetDefaultOwnedItem(testItemDef, testCard)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()

	context := NewTestContextBuilder(w).
		SetDefaultUrl().
		SetEndpoint("/user/cards/virtual/"+testBusiness.PublicId).
		SetUser(testUser).
		SetMethod("GET").
		SetHeader("Content-Type", "application/json").
		SetDefaultToken().
		Context

	respBodyExpected := &api.GetUserVirtualCardResponse{
		Points: int32(testCard.Points),
		OwnedItems: []api.OwnedItemApiModel{
			{
				PublicId:     testOwnedItem.PublicId,
				DefinitionId: testItemDef.PublicId,
			},
		},
		BusinessDetails: api.PublicBusinessDetailsApiModel{
			PublicId: testBusiness.PublicId,
			Name:     testBusiness.Name,
			Address:  testBusiness.Address,
			// GpsCoordinates: testBusiness.GPSCoordinates, TODO GPSCoordinates to string
			BannerImageId: testBusiness.BannerImageId,
			IconImageId:   testBusiness.IconImageId,
			MenuImageIds: []string{
				"bXU5YWltMm1haUdpCg",
			},
			ItemDefinitions: []api.ItemDefinitionApiModel{
				{
					PublicId:    testItemDef.PublicId,
					Name:        testItemDef.Name,
					Price:       int32(testItemDef.Price),
					Description: testItemDef.Description,
					ImageId:     testItemDef.ImageId,
					StartDate:   &testItemDef.StartDate.Time,
					EndDate:     &testItemDef.EndDate.Time,
					MaxAmount:   int32(testItemDef.MaxAmount),
					Available:   testItemDef.Available,
				},
			},
		},
	}

	// test env prep
	ctrl := gomock.NewController(t)
	handler := getVirtualCardHandlers(ctrl)

	handler.virtualCardManager.(*MockVirtualCardManager).
		EXPECT().
		Create(
			gomock.Eq(testUser),
			gomock.Eq(testBusiness.PublicId),
		).
		Return(
			testCard,
			nil,
		)

	handler.getCard(context)

	respBody, respCode, respParseErr := ExtractResponse[api.GetUserVirtualCardResponse](w)

	require.Nilf(t, respParseErr, "Failed to parse JSON response")
	require.Equalf(t, respCode, int(200), "Response returned unexpected status code")
	require.Truef(t, MatchEntities(respBodyExpected, respBody), "Response returned unexpected body contents")
	// TODO: test MatchEntities and gomock.Eq
}

func TestUserVirtualCardHandlersPostItemOk(t *testing.T) {
	testUser := GetDefaultUser()
	testBusinessUser := GetDefaultUser()
	testBusiness := GetDefaultBusiness(testBusinessUser)
	testCard := GetTestVirtualCard(nil, testUser, testBusiness)
	testItemDef := GetDefaultItem(testBusiness)
	testOwnedItem := GetDefaultOwnedItem(testItemDef, testCard)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()

	context := NewTestContextBuilder(w).
		SetDefaultUrl().
		SetEndpoint("/user/cards/virtual/"+testBusiness.PublicId+"/items/"+testItemDef.PublicId).
		SetUser(testUser).
		SetMethod("POST").
		SetHeader("Content-Type", "application/json").
		SetDefaultToken().
		Context

	respBodyExpected := &api.PostUserVirtualCardItemResponse{
		ItemId: testOwnedItem.PublicId,
	}

	// test env prep
	ctrl := gomock.NewController(t)
	handler := getVirtualCardHandlers(ctrl)

	handler.userAuthorizedAcessor.(*MockUserAuthorizedAccessor).
		EXPECT().
		Get(
			gomock.Eq(testUser),
			gomock.Eq(database.VirtualCard{PublicId: testBusiness.PublicId}), // Q: equivalence of vcard id and business id?
		).
		Return(
			testCard,
			nil,
		)

	handler.userAuthorizedAcessor.(*MockUserAuthorizedAccessor).
		EXPECT().
		Get(
			gomock.Eq(testUser),
			gomock.Eq(database.ItemDefinition{PublicId: testItemDef.PublicId}),
		).
		Return(
			testItemDef,
			nil,
		)

	handler.virtualCardManager.(*MockVirtualCardManager).
		EXPECT().
		BuyItem(
			gomock.Eq(testCard),
			gomock.Eq(testItemDef),
		).
		Return(
			testOwnedItem,
			nil,
		)

	handler.postItem(context)

	respBody, respCode, respParseErr := ExtractResponse[api.PostUserVirtualCardItemResponse](w)

	require.Nilf(t, respParseErr, "Failed to parse JSON response")
	require.Equalf(t, respCode, int(201), "Response returned unexpected status code")
	require.Truef(t, MatchEntities(respBodyExpected, respBody), "Response returned unexpected body contents")
	// TODO: test MatchEntities and gomock.Eq
}

func TestUserVirtualCardHandlersDeleteItemOk(t *testing.T) {
	testUser := GetDefaultUser()
	testBusinessUser := GetDefaultUser()
	testBusiness := GetDefaultBusiness(testBusinessUser)
	testCard := GetTestVirtualCard(nil, testUser, testBusiness)
	testItemDef := GetDefaultItem(testBusiness)
	testOwnedItem := GetDefaultOwnedItem(testItemDef, testCard)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()

	context := NewTestContextBuilder(w).
		SetDefaultUrl().
		SetEndpoint("/user/cards/virtual/"+testBusiness.PublicId+"/items/"+testOwnedItem.PublicId). // Q: differs from post endpoint: owned vs definition id
		SetUser(testUser).
		SetMethod("DELETE").
		SetHeader("Content-Type", "application/json").
		SetDefaultToken().
		Context

	respBodyExpected := &api.DefaultResponse{Status: api.OK}

	// test env prep
	ctrl := gomock.NewController(t)
	handler := getVirtualCardHandlers(ctrl)

	handler.userAuthorizedAcessor.(*MockUserAuthorizedAccessor).
		EXPECT().
		Get(
			gomock.Eq(testUser),
			gomock.Eq(database.OwnedItem{PublicId: testOwnedItem.PublicId}),
		).
		Return(
			testOwnedItem,
			nil,
		)

	handler.virtualCardManager.(*MockVirtualCardManager).
		EXPECT().
		ReturnItem(gomock.Eq(testOwnedItem)).
		Return(nil)

	handler.deleteItem(context)

	respBody, respCode, respParseErr := ExtractResponse[api.DefaultResponse](w)

	require.Nilf(t, respParseErr, "Failed to parse JSON response")
	require.Equalf(t, respCode, int(201), "Response returned unexpected status code")
	require.Truef(t, MatchEntities(respBodyExpected, respBody), "Response returned unexpected body contents")
	// TODO: test MatchEntities and gomock.Eq
}

func TestUserVirtualCardHandlersPostTransactionOk(t *testing.T) {
	testUser := GetDefaultUser()
	testBusinessUser := GetDefaultUser()
	testBusiness := GetDefaultBusiness(testBusinessUser)
	testCard := GetTestVirtualCard(nil, testUser, testBusiness)
	testItemDef := GetDefaultItem(testBusiness)
	testOwnedItem := GetDefaultOwnedItem(testItemDef, testCard)
	testTransaction, _ := GetTestTransaction(
		nil,
		testCard,
		[]database.OwnedItem{*testOwnedItem},
	)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()

	payload := api.PostUserVirtualCardTransactionRequest{
		ItemIds: []string{
			testOwnedItem.PublicId,
		},
	}
	payloadJson, _ := json.Marshal(payload)

	context := NewTestContextBuilder(w).
		SetDefaultUrl().
		SetEndpoint("/user/cards/virtual/"+testBusiness.PublicId+"/transaction").
		SetUser(testUser).
		SetMethod("POST").
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetDefaultToken().
		SetBody(payloadJson).
		Context

	respBodyExpected := &api.PostUserVirtualCardTransactionResponse{
		PublicId: testTransaction.PublicId,
		Code:     testTransaction.Code,
	}

	// test env prep
	ctrl := gomock.NewController(t)
	handler := getVirtualCardHandlers(ctrl)

	handler.userAuthorizedAcessor.(*MockUserAuthorizedAccessor).
		EXPECT().
		Get(
			gomock.Eq(testUser),
			gomock.Eq(database.VirtualCard{PublicId: testBusiness.PublicId}), // Q: vcard business id equivalence
		).
		Return(
			testCard,
			nil,
		)

	handler.transactionManager.(*MockTransactionManager).
		EXPECT().
		Start(
			gomock.Eq(testCard),
			gomock.Eq([]database.OwnedItem{*testOwnedItem}),
		).
		Return(
			testTransaction,
			nil,
		)

	handler.postTransaction(context)

	respBody, respCode, respParseErr := ExtractResponse[api.PostUserVirtualCardTransactionResponse](w)

	require.Nilf(t, respParseErr, "Failed to parse JSON response")
	require.Equalf(t, respCode, int(200), "Response returned unexpected status code")
	require.Truef(t, MatchEntities(respBodyExpected, respBody), "Response returned unexpected body contents")
	// TODO: test MatchEntities and gomock.Eq
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
