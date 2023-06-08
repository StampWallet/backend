package api

import (
	"encoding/json"
	"log"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	api "github.com/StampWallet/backend/internal/api/models"
	apiUtils "github.com/StampWallet/backend/internal/api/utils"
	"github.com/StampWallet/backend/internal/database"
	accessors "github.com/StampWallet/backend/internal/database/accessors"
	. "github.com/StampWallet/backend/internal/database/accessors/mocks"
	"github.com/StampWallet/backend/internal/managers"
	. "github.com/StampWallet/backend/internal/managers/mocks"
	. "github.com/StampWallet/backend/internal/testutils"
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
		virtualCardManager:            NewMockVirtualCardManager(ctrl),
		transactionManager:            NewMockTransactionManager(ctrl),
		itemDefinitionManager:         NewMockItemDefinitionManager(ctrl),
		userAuthorizedAcessor:         NewMockUserAuthorizedAccessor(ctrl),
		authorizedTransactionAccessor: NewMockAuthorizedTransactionAccessor(ctrl),
		logger:                        log.Default(),
	}
}

func getLocalCardHandlers(ctrl *gomock.Controller) *UserLocalCardHandlers {
	return &UserLocalCardHandlers{
		localCardManager:      NewMockLocalCardManager(ctrl),
		userAuthorizedAcessor: NewMockUserAuthorizedAccessor(ctrl),
		logger:                log.Default(),
	}
}

type pointerMatcherTypes interface {
	string | int | uint
}

// pointerMatcher.Matches(x) checks if *x == pointerMatcher.value
type pointerMatcher[T pointerMatcherTypes] struct {
	value T
}

func (matcher *pointerMatcher[T]) Matches(x interface{}) bool {
	return reflect.ValueOf(x).Elem().Equal(reflect.ValueOf(matcher.value))
}

func (matcher *pointerMatcher[T]) String() string {
	return "pointerMatcher"
}

// TODO
func TestUserHandlersGetUserCardsOk(t *testing.T) {
	testUser := GetDefaultUser()
	testBusinessUser := GetDefaultUser()
	testBusiness := GetDefaultBusiness(testBusinessUser)
	testLocalCard := GetTestLocalCard(nil, testUser)
	testVirtualCard := GetTestVirtualCard(nil, testUser, testBusiness)
	testVirtualCard.Business = testBusiness

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()

	context := NewTestContextBuilder(w).
		SetDefaultUrl().
		SetEndpoint("/user/cards").
		SetUser(testUser).
		SetMethod("GET").
		SetHeader("Content-Type", "application/json").
		SetDefaultToken().
		Context

	respBodyExpected := &api.GetUserCardsResponse{
		LocalCards: []api.LocalCardApiModel{
			{
				PublicId: testLocalCard.PublicId,
				Name:     testLocalCard.Name,
				Type:     testLocalCard.Type,
				Code:     testLocalCard.Code,
			},
		},
		VirtualCards: []api.ShortVirtualCardApiModel{
			{
				BusinessDetails: api.ShortBusinessDetailsApiModel{
					PublicId:       testBusiness.PublicId,
					Name:           testBusiness.Name,
					Description:    testBusiness.Description,
					GpsCoordinates: testBusiness.GPSCoordinates.ToString(),
					BannerImageId:  testBusiness.BannerImageId,
					IconImageId:    testBusiness.IconImageId,
				},
				Points: int32(testVirtualCard.Points),
			},
		},
	}

	// test env prep
	ctrl := gomock.NewController(t)
	handler := getUserHandlers(ctrl)

	handler.userAuthorizedAcessor.(*MockUserAuthorizedAccessor).
		EXPECT().
		GetAll(gomock.Eq(testUser), &database.LocalCard{}, []string{}).
		Return(([]accessors.UserOwnedEntity{testLocalCard}), nil)

	handler.userAuthorizedAcessor.(*MockUserAuthorizedAccessor).
		EXPECT().
		GetAll(gomock.Eq(testUser), &database.VirtualCard{}, []string{"Business"}).
		Return([]accessors.UserOwnedEntity{testVirtualCard}, nil)

	handler.getUserCards(context)

	respBody, respCode, respParseErr := ExtractResponse[api.GetUserCardsResponse](w)

	require.Nilf(t, respParseErr, "Failed to parse JSON response")
	require.Equalf(t, int(200), respCode, "Response returned unexpected status code")
	require.Truef(t, reflect.DeepEqual(*respBodyExpected, *respBody), "Response returned unexpected body contents")
}

func TestUserHandlersGetSearchBusinessesOk(t *testing.T) {
	// TODO caly test case do napisania
	testUser := GetDefaultUser()
	testBusinessUser := GetDefaultUser()
	testBusiness := GetDefaultBusiness(testBusinessUser)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()

	context := NewTestContextBuilder(w).
		SetDefaultUrl().
		SetEndpoint("/user/businesses").
		AddQueryParam("text", "example business search").
		SetUser(testUser).
		SetMethod("GET").
		SetHeader("Content-Type", "application/json").
		SetDefaultToken().
		Context

	respBodyExpected := api.GetUserBusinessesSearchResponse{
		Businesses: []api.ShortBusinessDetailsApiModel{
			{
				PublicId:       testBusiness.PublicId,
				Name:           testBusiness.Name,
				Description:    testBusiness.Description,
				GpsCoordinates: testBusiness.GPSCoordinates.ToString(),
				BannerImageId:  testBusiness.BannerImageId,
				IconImageId:    testBusiness.IconImageId,
			},
		},
	}

	// test env prep
	ctrl := gomock.NewController(t)
	handler := getUserHandlers(ctrl)

	// setup mocks
	handler.businessManager.(*MockBusinessManager).
		EXPECT().
		Search(
			&pointerMatcher[string]{"example business search"},
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
		).
		Return([]database.Business{*testBusiness}, nil)

	handler.getSearchBusinesses(context)

	respBody, respCode, respParseErr := ExtractResponse[api.GetUserBusinessesSearchResponse](w)

	require.Nilf(t, respParseErr, "Failed to parse JSON response")
	require.Equalf(t, respCode, int(200), "Response returned unexpected status code")
	require.Truef(t, reflect.DeepEqual(respBodyExpected, *respBody), "Response returned unexpected body contents")
}

func TestUserVirtualCardHandlersPostCardOk(t *testing.T) {
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
		SetParam("businessId", testBusiness.PublicId).
		Context

	respBodyExpected := &api.DefaultResponse{Status: api.CREATED}

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

	handler.postCard(context)

	respBody, respCode, respParseErr := ExtractResponse[api.DefaultResponse](w)

	require.Nilf(t, respParseErr, "Failed to parse JSON response")
	require.Equalf(t, respCode, int(201), "Response returned unexpected status code")
	require.Truef(t, reflect.DeepEqual(respBodyExpected, respBody), "Response returned unexpected body contents")
}

func TestUserVirtualCardHandlersDeleteCardOk(t *testing.T) {
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
		SetParam("businessId", testBusiness.PublicId).
		Context

	respBodyExpected := &api.DefaultResponse{Status: api.OK}

	// test env prep
	ctrl := gomock.NewController(t)
	handler := getVirtualCardHandlers(ctrl)

	handler.userAuthorizedAcessor.(*MockUserAuthorizedAccessor).
		EXPECT().
		Get(
			gomock.Eq(testUser),
			gomock.Eq(&database.VirtualCard{Business: &database.Business{PublicId: testBusiness.PublicId}}),
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
	require.Truef(t, reflect.DeepEqual(respBodyExpected, respBody), "Response returned unexpected body contents")
}

func TestUserVirtualCardHandlersGetCardOk(t *testing.T) {
	testUser := GetDefaultUser()
	testBusinessUser := GetDefaultUser()
	testBusiness := GetDefaultBusiness(testBusinessUser)
	testCard := GetTestVirtualCard(nil, testUser, testBusiness)
	testItemDef := GetDefaultItem(testBusiness)
	testOwnedItem := GetDefaultOwnedItem(testItemDef, testCard)
	testMenuImage := GetTestMenuImage(nil, testBusiness)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()

	context := NewTestContextBuilder(w).
		SetDefaultUrl().
		SetEndpoint("/user/cards/virtual/"+testBusiness.PublicId).
		SetUser(testUser).
		SetMethod("GET").
		SetHeader("Content-Type", "application/json").
		SetDefaultToken().
		SetParam("businessId", testBusiness.PublicId).
		Context

	sd := testItemDef.StartDate.Time.Truncate(time.Duration(0))
	ed := testItemDef.EndDate.Time.Truncate(time.Duration(0))

	respBodyExpected := &api.GetUserVirtualCardResponse{
		Points: int32(testCard.Points),
		OwnedItems: []api.OwnedItemApiModel{
			{
				PublicId:     testOwnedItem.PublicId,
				DefinitionId: testItemDef.PublicId,
			},
		},
		BusinessDetails: api.PublicBusinessDetailsApiModel{
			PublicId:       testBusiness.PublicId,
			Name:           testBusiness.Name,
			Address:        testBusiness.Address,
			Description:    testBusiness.Description,
			GpsCoordinates: testBusiness.GPSCoordinates.ToString(),
			BannerImageId:  testBusiness.BannerImageId,
			IconImageId:    testBusiness.IconImageId,
			MenuImageIds: []string{
				testMenuImage.FileId,
			},
			ItemDefinitions: []api.ItemDefinitionApiModel{
				{
					PublicId:    testItemDef.PublicId,
					Name:        testItemDef.Name,
					Price:       int32(testItemDef.Price),
					Description: testItemDef.Description,
					ImageId:     testItemDef.ImageId,
					StartDate:   &sd,
					EndDate:     &ed,
					MaxAmount:   int32(testItemDef.MaxAmount),
					Available:   testItemDef.Available,
				},
			},
		},
	}

	// test env prep
	ctrl := gomock.NewController(t)
	handler := getVirtualCardHandlers(ctrl)

	testCard.Business = testBusiness
	testCard.OwnedItems = []database.OwnedItem{*testOwnedItem}
	testBusiness.ItemDefinitions = []database.ItemDefinition{*testItemDef}
	testBusiness.MenuImages = []database.MenuImage{*testMenuImage}

	handler.userAuthorizedAcessor.(*MockUserAuthorizedAccessor).
		EXPECT().
		GetAll(
			gomock.Eq(testUser),
			gomock.Eq(&database.VirtualCard{Business: &database.Business{PublicId: testBusiness.PublicId}}),
			gomock.Eq([]string{"Business", "Business.OwnedItems", "Business.MenuImages.FileId",
				"OwnedItems", "OwnedItems.ItemDefinition.PublicId"}),
		).
		Return(
			[]accessors.UserOwnedEntity{testCard},
			nil,
		)

	handler.getCard(context)

	respBody, respCode, respParseErr := ExtractResponse[api.GetUserVirtualCardResponse](w)

	require.Nilf(t, respParseErr, "Failed to parse JSON response")
	require.Equalf(t, respCode, int(200), "Response returned unexpected status code")
	require.Truef(t, reflect.DeepEqual(respBodyExpected, respBody), "Response returned unexpected body contents")
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
		SetParam("businessId", testBusiness.PublicId).
		SetParam("itemDefinitionId", testItemDef.PublicId).
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
			gomock.Eq(&database.VirtualCard{Business: &database.Business{PublicId: testBusiness.PublicId}}), // Q: equivalence of vcard id and business id?
		).
		Return(
			testCard,
			nil,
		)

	handler.virtualCardManager.(*MockVirtualCardManager).
		EXPECT().
		BuyItem(
			gomock.Eq(testCard),
			gomock.Eq(testItemDef.PublicId),
		).
		Return(
			testOwnedItem,
			nil,
		)

	handler.postItemDefinition(context)

	respBody, respCode, respParseErr := ExtractResponse[api.PostUserVirtualCardItemResponse](w)

	require.Nilf(t, respParseErr, "Failed to parse JSON response")
	require.Equalf(t, respCode, int(201), "Response returned unexpected status code")
	require.Truef(t, reflect.DeepEqual(respBodyExpected, respBody), "Response returned unexpected body contents")
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
		SetParam("businessId", testBusiness.PublicId).
		SetParam("itemId", testOwnedItem.PublicId).
		Context

	respBodyExpected := &api.DefaultResponse{Status: api.OK}

	// test env prep
	ctrl := gomock.NewController(t)
	handler := getVirtualCardHandlers(ctrl)

	handler.userAuthorizedAcessor.(*MockUserAuthorizedAccessor).
		EXPECT().
		Get(
			gomock.Eq(testUser),
			gomock.Eq(&database.VirtualCard{Business: &database.Business{PublicId: testBusiness.PublicId}}),
		).
		Return(
			testCard,
			nil,
		)

	handler.userAuthorizedAcessor.(*MockUserAuthorizedAccessor).
		EXPECT().
		Get(
			gomock.Eq(testUser),
			gomock.Eq(&database.OwnedItem{PublicId: testOwnedItem.PublicId, VirtualCardId: testCard.ID}),
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
	require.Equalf(t, respCode, int(200), "Response returned unexpected status code")
	require.Truef(t, reflect.DeepEqual(respBodyExpected, respBody), "Response returned unexpected body contents")
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
		SetParam("businessId", testBusiness.PublicId).
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
			gomock.Eq(&database.VirtualCard{Business: &database.Business{PublicId: testBusiness.PublicId}}),
			// Q: vcard business id equivalence
		).
		Return(
			testCard,
			nil,
		)

	handler.virtualCardManager.(*MockVirtualCardManager).
		EXPECT().
		FilterOwnedItems(
			gomock.Eq(testCard),
			gomock.Eq([]string{testOwnedItem.PublicId}),
		).
		Return(
			[]database.OwnedItem{*testOwnedItem},
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
	require.Equalf(t, respCode, int(201), "Response returned unexpected status code")
	require.Truef(t, reflect.DeepEqual(respBodyExpected, respBody), "Response returned unexpected body contents")
}

func TestUserVirtualCardHandlersGetTransactionOk(t *testing.T) {
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

	context := NewTestContextBuilder(w).
		SetDefaultUrl().
		SetEndpoint("/user/cards/virtual/"+testBusiness.PublicId+"/transaction/"+testTransaction.Code).
		SetUser(testUser).
		SetMethod("POST").
		SetHeader("Accept", "application/json").
		SetDefaultToken().
		SetParam("businessId", testBusiness.PublicId).
		SetParam("transactionCode", testTransaction.Code).
		Context

	respBodyExpected := &api.GetUserVirtualCardTransactionResponse{
		PublicId:    testTransaction.PublicId,
		State:       apiUtils.ConvertDbTransactionState(testTransaction.State),
		AddedPoints: int32(testTransaction.AddedPoints),
		ItemActions: []api.ItemActionApiModel{
			{
				ItemId: testOwnedItem.PublicId,
				Action: apiUtils.ConvertDbItemAction(database.NoActionType),
			},
		},
	}

	// test env prep
	ctrl := gomock.NewController(t)
	handler := getVirtualCardHandlers(ctrl)

	handler.authorizedTransactionAccessor.(*MockAuthorizedTransactionAccessor).
		EXPECT().
		GetForUser(
			gomock.Eq(testUser),
			gomock.Eq(testTransaction.Code),
		).
		Return(
			testTransaction,
			nil,
		)

	handler.getTransaction(context)

	respBody, respCode, respParseErr := ExtractResponse[api.GetUserVirtualCardTransactionResponse](w)

	require.Nilf(t, respParseErr, "Failed to parse JSON response")
	require.Equalf(t, respCode, int(200), "Response returned unexpected status code")
	require.Truef(t, reflect.DeepEqual(respBodyExpected, respBody), "Response returned unexpected body contents")
}

func testUserVirtualCardHandlersGetTransactionNotOk_tmpl(t *testing.T, err error) {
	testUser := GetDefaultUser()
	testBusinessUser := GetDefaultUser()
	testBusiness := GetDefaultBusiness(testBusinessUser)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()

	context := NewTestContextBuilder(w).
		SetDefaultUrl().
		SetEndpoint("/user/cards/virtual/"+testBusiness.PublicId+"/transaction/"+"nothing").
		SetUser(testUser).
		SetMethod("POST").
		SetHeader("Accept", "application/json").
		SetDefaultToken().
		SetParam("businessId", testBusiness.PublicId).
		SetParam("transactionCode", "nothing").
		Context

	respBodyExpected := &api.DefaultResponse{Status: api.NOT_FOUND}

	// test env prep
	ctrl := gomock.NewController(t)
	handler := getVirtualCardHandlers(ctrl)

	handler.authorizedTransactionAccessor.(*MockAuthorizedTransactionAccessor).
		EXPECT().
		GetForUser(
			gomock.Eq(testUser),
			gomock.Eq("nothing"),
		).
		Return(
			nil,
			err,
		)

	handler.getTransaction(context)

	respBody, respCode, respParseErr := ExtractResponse[api.DefaultResponse](w)

	require.Nilf(t, respParseErr, "Failed to parse JSON response")
	require.Equalf(t, respCode, int(404), "Response returned unexpected status code")
	require.Truef(t, reflect.DeepEqual(respBodyExpected, respBody), "Response returned unexpected body contents")
}

func TestUserVirtualCardHandlersGetTransactionNotOk(t *testing.T) {
	testUserVirtualCardHandlersGetTransactionNotOk_tmpl(t, accessors.ErrNoAccess)
	testUserVirtualCardHandlersGetTransactionNotOk_tmpl(t, accessors.ErrNotFound)
}

// UserLocalCardHandlers

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

	respBodyExpected := &api.PostUserLocalCardsResponse{PublicId: testCard.PublicId}

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

	respBody, respCode, respParseErr := ExtractResponse[api.PostUserLocalCardsResponse](w)

	require.Nilf(t, respParseErr, "Failed to parse JSON response")
	require.Equalf(t, respCode, int(201), "Response returned unexpected status code")
	require.Truef(t, reflect.DeepEqual(respBodyExpected, respBody), "Response returned unexpected body contents")
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
		SetParam("cardId", testCard.PublicId).
		Context

	respBodyExpected := &api.DefaultResponse{Status: api.OK}

	// test env prep
	ctrl := gomock.NewController(t)
	handler := getLocalCardHandlers(ctrl)

	handler.userAuthorizedAcessor.(*MockUserAuthorizedAccessor).
		EXPECT().
		Get(
			gomock.Eq(testUser),
			gomock.Eq(&database.LocalCard{PublicId: testCard.PublicId}),
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
	require.Truef(t, reflect.DeepEqual(respBodyExpected, respBody), "Response returned unexpected body contents")
}
