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

	// payload := api.GetBusinessItemDefinitionsReq

	// // test env prep
	// ctrl := gomock.NewController(t)
	// handler := getItemHandlers(ctrl)

	// handler.userAuthorizedAcessor.(*MockUserAuthorizedAccessor).
	// 	EXPECT().
	// 	Get(
	// 		gomock.Eq(testBusinessUser),
	// 		gomock.Eq(&database.Business{PublicId: testBusiness.PublicId}),
	// 	).
	// 	Return(
	// 		testBusiness,
	// 		nil,
	// 	)

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

func TestItemDefinitionHandlersPostItemDefinitionOk(t *testing.T) {
	testBusinessUser := GetDefaultUser()
	testBusiness := GetDefaultBusiness(testBusinessUser)
	testItemDef := GetDefaultItem(testBusiness)

	testItemDetails := managers.ItemDetails{
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
	w := httptest.NewRecorder()
	context := NewTestContextBuilder(w).
		SetDefaultUrl().
		SetEndpoint("/business/itemDefinitions").
		SetUser(testBusinessUser).
		SetMethod("POST").
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetDefaultToken().
		SetBody(payloadJson).
		Context

	// test env prep
	ctrl := gomock.NewController(t)
	handler := getItemHandlers(ctrl)

	handler.userAuthorizedAcessor.(*MockUserAuthorizedAccessor).
		EXPECT().
		Get(
			gomock.Eq(testBusinessUser),
			gomock.Eq(&database.Business{PublicId: testBusiness.PublicId}),
		).
		Return(
			testBusiness,
			nil,
		)

	handler.itemDefinitionManager.(*MockItemDefinitionManager).
		EXPECT().
		AddItem(
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

func TestItemDefinitionHandlersPatchItemDefinitionOk(t *testing.T) {
	testBusinessUser := GetDefaultUser()
	testBusiness := GetDefaultBusiness(testBusinessUser)
	testItemDef := GetDefaultItem(testBusiness)

	newItemDetails := managers.ItemDetails{
		Name:        "new item name",
		Price:       Ptr(uint(2137)), // sorry
		Description: "new item description",
		StartDate:   Ptr(time.Now()),
		EndDate:     Ptr(time.Now().Add(time.Hour * 24)),
		MaxAmount:   Ptr(uint(20)),
		Available:   Ptr(true),
	}

	newItemDef := &database.ItemDefinition{
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
	w := httptest.NewRecorder()
	context := NewTestContextBuilder(w).
		SetDefaultUrl().
		SetEndpoint("/business/itemDefinitions"+testItemDef.PublicId).
		SetUser(testBusinessUser).
		SetMethod("PATCH").
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetDefaultToken().
		SetBody(payloadJson).
		Context

	// test env prep
	ctrl := gomock.NewController(t)
	handler := getItemHandlers(ctrl)

	handler.userAuthorizedAcessor.(*MockUserAuthorizedAccessor).
		EXPECT().
		Get(
			gomock.Eq(testBusinessUser),
			gomock.Eq(&database.Business{PublicId: testBusiness.PublicId}),
		).
		Return(
			testBusiness,
			nil,
		)

	handler.businessAuthorizedAccessor.(*MockBusinessAuthorizedAccessor).
		EXPECT().
		Get(
			gomock.Eq(testBusiness),
			gomock.Eq(&database.ItemDefinition{PublicId: payload.PublicId}),
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

	handler.userAuthorizedAcessor.(*MockUserAuthorizedAccessor).
		EXPECT().
		Get(
			gomock.Eq(testBusinessUser),
			gomock.Eq(&database.Business{PublicId: testBusiness.PublicId}),
		).
		Return(
			testBusiness,
			nil,
		)

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

// FUTURE: Rainy day scenarios tests (sequence diagrams)
