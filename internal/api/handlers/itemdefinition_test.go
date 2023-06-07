package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http/httptest"
	"testing"
	"time"

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

func getItemHandlers(ctrl *gomock.Controller) *ItemDefinitionHandlers {
	return &ItemDefinitionHandlers{
		itemDefinitionManager:      NewMockItemDefinitionManager(ctrl),
		userAuthorizedAcessor:      NewMockUserAuthorizedAccessor(ctrl),
		businessAuthorizedAccessor: NewMockBusinessAuthorizedAccessor(ctrl),
		logger:                     log.Default(),
	}
}

func TestItemDefinitionHandlersGetItemDefinitionOk(t *testing.T) {
	// testBusinessUser := GetDefaultUser()
	// testBusiness := GetDefaultBusiness(testBusinessUser)
	// testItemDef := GetDefaultItem(testBusiness)

	// // test env prep
	// ctrl := gomock.NewController(t)
	// handler := getItemHandlers(ctrl)

	// handler.businessManager.(*MockBusinessManager).
	// 	EXPECT().
	// 	Create(
	// 		gomock.Eq(testBusinessUser),
	// 		gomock.Eq(testBusinessDetails),
	// 	).
	// 	Return(
	// 		testBusiness,
	// 		nil,
	// 	)

	// handler.postAccount(context)

	// respBody, respCode, respParseErr := ExtractResponse[api.PostBusinessAccountResponse](w)

	// require.Nilf(t, respParseErr, "Failed to parse JSON response")
	// require.Equalf(t, respCode, int(200), "Response returned unexpected status code")
	// require.Truef(t, MatchEntities(respBodyExpected, respBody), "Response returned unexpected body contents")
	// // TODO: test MatchEntities and gomock.Eq
	// TODO: This endpoint needs a corrected openapi spec definition
}

func setupItemDefinitionHandlersPostItemDefinition() (
	w *httptest.ResponseRecorder,
	context *gin.Context,
	testBusinessUser *database.User,
	testBusiness *database.Business,
	testItemDetails *managers.ItemDetails,
	testItemDef *database.ItemDefinition,
) {
	testBusinessUser = GetDefaultUser()
	testBusiness = GetDefaultBusiness(testBusinessUser)
	testItemDef = GetDefaultItem(testBusiness)

	testItemDetails = &managers.ItemDetails{
		Name:        testItemDef.Name,
		Price:       &testItemDef.Price,
		Description: testItemDef.Description,
		StartDate:   &testItemDef.StartDate.Time,
		EndDate:     &testItemDef.EndDate.Time,
		MaxAmount:   &testItemDef.MaxAmount,
		Available:   &testItemDef.Available,
	}

	payload := api.PostBusinessItemDefinitionRequest{
		PublicId:    testItemDef.PublicId,
		Name:        testItemDef.Name,
		Price:       int32(testItemDef.Price),
		Description: testItemDef.Description,
		ImageId:     testItemDef.ImageId,
		StartDate:   &testItemDef.StartDate.Time,
		EndDate:     &testItemDef.EndDate.Time,
		MaxAmount:   int32(testItemDef.MaxAmount),
		Available:   testItemDef.Available,
	}
	payloadJson, _ := json.Marshal(payload)

	gin.SetMode(gin.TestMode)
	w = httptest.NewRecorder()
	context = NewTestContextBuilder(w).
		SetDefaultUrl().
		SetEndpoint("/business/itemDefinitions").
		SetUser(testBusinessUser).
		SetMethod("POST").
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetDefaultToken().
		SetBody(payloadJson).
		Context

	return w, context, testBusinessUser, testBusiness, testItemDetails, testItemDef
}

func TestItemDefinitionHandlersPostItemDefinitionOk(t *testing.T) {
	w, context, _, testBusiness, testItemDetails, testItemDef := setupItemDefinitionHandlersPostItemDefinition()

	// test env prep
	ctrl := gomock.NewController(t)
	handler := getItemHandlers(ctrl)

	handler.itemDefinitionManager.(*MockItemDefinitionManager).
		EXPECT().
		AddItem(
			gomock.Eq(testBusiness),
			gomock.Eq(testBusiness),
			gomock.Eq(testItemDetails),
		).
		Return(
			testItemDef,
			nil,
		)

	handler.postItemDefinition(context)

	respBodyExpected := api.DefaultResponse{Status: api.OK}
	respBody, respCode, respParseErr := ExtractResponse[api.DefaultResponse](w)

	require.Nilf(t, respParseErr, "Failed to parse JSON response")
	require.Equalf(t, respCode, int(201), "Response returned unexpected status code")
	require.Truef(t, MatchEntities(respBodyExpected, respBody), "Response returned unexpected body contents")
}

func TestItemDefinitionHandlersPostItemDefinitionNok_BadDef(t *testing.T) {
	w, context, testBusinessUser, testBusiness, testItemDetails, _ := setupItemDefinitionHandlersPostItemDefinition()

	// test env prep
	ctrl := gomock.NewController(t)
	handler := getItemHandlers(ctrl)

	handler.itemDefinitionManager.(*MockItemDefinitionManager).
		EXPECT().
		AddItem(
			gomock.Eq(testBusinessUser),
			gomock.Eq(testBusiness),
			gomock.Eq(testItemDetails),
		).
		Return(
			nil,
			managers.ErrInvalidItemDetails,
		)

	handler.postItemDefinition(context)

	respBodyExpected := api.DefaultResponse{Status: api.INVALID_REQUEST}
	respBody, respCode, respParseErr := ExtractResponse[api.DefaultResponse](w)

	require.Nilf(t, respParseErr, "Failed to parse JSON response")
	require.Equalf(t, respCode, int(400), "Response returned unexpected status code")
	require.Truef(t, MatchEntities(respBodyExpected, respBody), "Response returned unexpected body contents")
}

func setupItemDefinitionHandlersPatchItemDefinition() (
	w *httptest.ResponseRecorder,
	context *gin.Context,
	testBusiness *database.Business,
	testItemDef *database.ItemDefinition,
	newItemDef *database.ItemDefinition,
	newItemDetails *managers.ItemDetails,
) {
	testBusinessUser := GetDefaultUser()
	testBusiness = GetDefaultBusiness(testBusinessUser)
	testItemDef = GetDefaultItem(testBusiness)

	newItemDetails = &managers.ItemDetails{
		Name:        "new item name",
		Price:       Ptr(uint(2137)), // sorry
		Description: "new item description",
		StartDate:   Ptr(time.Now()),
		EndDate:     Ptr(time.Now().Add(time.Hour * 24)),
		MaxAmount:   Ptr(uint(20)),
		Available:   Ptr(true),
	}

	newItemDef = &database.ItemDefinition{
		PublicId:    testItemDef.PublicId,
		BusinessId:  testBusiness.ID,
		Name:        newItemDetails.Name,
		Price:       *newItemDetails.Price,
		Description: newItemDetails.Description,
		ImageId:     testItemDef.ImageId,
		StartDate:   sql.NullTime{*newItemDetails.StartDate, true},
		EndDate:     sql.NullTime{*newItemDetails.EndDate, true},
		MaxAmount:   *newItemDetails.MaxAmount,
		Available:   *newItemDetails.Available,
		Withdrawn:   testItemDef.Withdrawn,
	}

	payload := api.PostBusinessItemDefinitionRequest{
		PublicId:    newItemDef.PublicId,
		Name:        newItemDef.Name,
		Price:       int32(newItemDef.Price),
		Description: newItemDef.Description,
		ImageId:     newItemDef.ImageId,
		StartDate:   &newItemDef.StartDate.Time,
		EndDate:     &newItemDef.EndDate.Time,
		MaxAmount:   int32(newItemDef.MaxAmount),
		Available:   newItemDef.Available,
	}
	payloadJson, _ := json.Marshal(payload)

	gin.SetMode(gin.TestMode)
	w = httptest.NewRecorder()
	context = NewTestContextBuilder(w).
		SetDefaultUrl().
		SetEndpoint("/business/itemDefinitions"+testItemDef.PublicId).
		SetUser(testBusinessUser).
		SetMethod("PATCH").
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetDefaultToken().
		SetBody(payloadJson).
		Context

	return w, context, testBusiness, testItemDef, newItemDef, newItemDetails
}

func TestItemDefinitionHandlersPatchItemDefinitionOk(t *testing.T) {
	w, context, testBusiness, testItemDef, newItemDef, newItemDetails := setupItemDefinitionHandlersPatchItemDefinition()

	// test env prep
	ctrl := gomock.NewController(t)
	handler := getItemHandlers(ctrl)

	handler.businessAuthorizedAccessor.(*MockBusinessAuthorizedAccessor).
		EXPECT().
		Get(
			gomock.Eq(testBusiness),
			gomock.Eq(&database.ItemDefinition{PublicId: testItemDef.PublicId}),
		).
		Return(
			testItemDef,
			nil,
		)

	handler.itemDefinitionManager.(*MockItemDefinitionManager).
		EXPECT().
		ChangeItemDetails(
			gomock.Eq(testItemDef),
			gomock.Eq(newItemDetails),
		).
		Return(
			newItemDef,
			nil,
		)

	handler.patchItemDefinition(context)

	respBodyExpected := api.DefaultResponse{Status: api.OK}
	respBody, respCode, respParseErr := ExtractResponse[api.DefaultResponse](w)

	require.Nilf(t, respParseErr, "Failed to parse JSON response")
	require.Equalf(t, respCode, int(200), "Response returned unexpected status code")
	require.Truef(t, MatchEntities(respBodyExpected, respBody), "Response returned unexpected body contents")
}

func TestItemDefinitionHandlersPatchItemDefinitionNok_InvDef(t *testing.T) {
	w, context, testBusiness, testItemDef, _, _ := setupItemDefinitionHandlersPatchItemDefinition()

	// test env prep
	ctrl := gomock.NewController(t)
	handler := getItemHandlers(ctrl)

	handler.businessAuthorizedAccessor.(*MockBusinessAuthorizedAccessor).
		EXPECT().
		Get(
			gomock.Eq(testBusiness),
			gomock.Eq(&database.ItemDefinition{PublicId: testItemDef.PublicId}),
		).
		Return(
			nil,
			managers.ErrInvalidItemDetails,
		)

	handler.patchItemDefinition(context)

	respBodyExpected := api.DefaultResponse{Status: api.NOT_FOUND}
	respBody, respCode, respParseErr := ExtractResponse[api.DefaultResponse](w)

	require.Nilf(t, respParseErr, "Failed to parse JSON response")
	require.Equalf(t, respCode, int(404), "Response returned unexpected status code")
	require.Truef(t, MatchEntities(respBodyExpected, respBody), "Response returned unexpected body contents")
}

func TestItemDefinitionHandlersPatchItemDefinitionNok_BadDef(t *testing.T) {
	w, context, testBusiness, testItemDef, _, newItemDetails := setupItemDefinitionHandlersPatchItemDefinition()

	// test env prep
	ctrl := gomock.NewController(t)
	handler := getItemHandlers(ctrl)

	handler.businessAuthorizedAccessor.(*MockBusinessAuthorizedAccessor).
		EXPECT().
		Get(
			gomock.Eq(testBusiness),
			gomock.Eq(&database.ItemDefinition{PublicId: testItemDef.PublicId}),
		).
		Return(
			testItemDef,
			nil,
		)

	handler.itemDefinitionManager.(*MockItemDefinitionManager).
		EXPECT().
		ChangeItemDetails(
			gomock.Eq(testItemDef),
			gomock.Eq(newItemDetails),
		).
		Return(
			nil,
			managers.ErrInvalidItemDetails,
		)

	handler.patchItemDefinition(context)

	respBodyExpected := api.DefaultResponse{Status: api.INVALID_REQUEST}
	respBody, respCode, respParseErr := ExtractResponse[api.DefaultResponse](w)

	require.Nilf(t, respParseErr, "Failed to parse JSON response")
	require.Equalf(t, respCode, int(400), "Response returned unexpected status code")
	require.Truef(t, MatchEntities(respBodyExpected, respBody), "Response returned unexpected body contents")
}

func TestItemDefinitionHandlersDeleteItemDefinitionOk(t *testing.T) {
	testBusinessUser := GetDefaultUser()
	testBusiness := GetDefaultBusiness(testBusinessUser)
	testItemDef := GetDefaultItem(testBusiness)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	context := NewTestContextBuilder(w).
		SetDefaultUrl().
		SetEndpoint("/business/itemDefinitions"+testItemDef.PublicId).
		SetUser(testBusinessUser).
		SetMethod("DELETE").
		SetHeader("Content-Type", "application/json").
		SetDefaultToken().
		Context

	// test env prep
	ctrl := gomock.NewController(t)
	handler := getItemHandlers(ctrl)

	handler.businessAuthorizedAccessor.(*MockBusinessAuthorizedAccessor).
		EXPECT().
		Get(
			gomock.Eq(testBusiness),
			gomock.Eq(&database.ItemDefinition{
				PublicId: testBusiness.ItemDefinitions[0].PublicId,
			}),
		).
		Return(
			testItemDef,
			nil,
		)

	handler.itemDefinitionManager.(*MockItemDefinitionManager).
		EXPECT().
		WithdrawItem(gomock.Eq(testItemDef)).
		Return(nil)

	handler.deleteItemDefinition(context)

	respBodyExpected := api.DefaultResponse{Status: api.OK}
	respBody, respCode, respParseErr := ExtractResponse[api.DefaultResponse](w)

	require.Nilf(t, respParseErr, "Failed to parse JSON response")
	require.Equalf(t, respCode, int(200), "Response returned unexpected status code")
	require.Truef(t, MatchEntities(respBodyExpected, respBody), "Response returned unexpected body contents")
}

func TestItemDefinitionHandlersDeleteItemDefinitionNok_InvDef(t *testing.T) {
	testBusinessUser1 := GetDefaultUser()
	testBusinessUser2 := GetDefaultUser()
	testBusiness1 := GetDefaultBusiness(testBusinessUser1)
	testBusiness2 := GetDefaultBusiness(testBusinessUser2)
	testItemDef := GetDefaultItem(testBusiness2)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	context := NewTestContextBuilder(w).
		SetDefaultUrl().
		SetEndpoint("/business/itemDefinitions"+testItemDef.PublicId).
		SetUser(testBusinessUser1).
		SetMethod("DELETE").
		SetHeader("Content-Type", "application/json").
		SetDefaultToken().
		Context

	// test env prep
	ctrl := gomock.NewController(t)
	handler := getItemHandlers(ctrl)

	handler.businessAuthorizedAccessor.(*MockBusinessAuthorizedAccessor).
		EXPECT().
		Get(
			gomock.Eq(testBusiness1),
			gomock.Eq(&database.ItemDefinition{
				PublicId: testItemDef.PublicId,
			}),
		).
		Return(
			nil,
			managers.ErrUnknownItem,
		)

	handler.deleteItemDefinition(context)

	respBodyExpected := api.DefaultResponse{Status: api.NOT_FOUND}
	respBody, respCode, respParseErr := ExtractResponse[api.DefaultResponse](w)

	require.Nilf(t, respParseErr, "Failed to parse JSON response")
	require.Equalf(t, respCode, int(404), "Response returned unexpected status code")
	require.Truef(t, MatchEntities(respBodyExpected, respBody), "Response returned unexpected body contents")
}
