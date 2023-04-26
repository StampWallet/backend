package managers

import (
	"log"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	. "github.com/StampWallet/backend/internal/database"
	//. "github.com/StampWallet/backend/internal/database/mocks"
	. "github.com/StampWallet/backend/internal/services"
	. "github.com/StampWallet/backend/internal/services/mocks"
	. "github.com/StampWallet/backend/internal/testutils"
)

func GetBusinessManager(ctrl *gomock.Controller) *BusinessManagerImpl {
	return &BusinessManagerImpl{
		&BaseServices{
			Logger:   log.Default(),
			Database: GetDatabase(),
		},
		NewMockFileStorageService(ctrl),
	}
}

func TestBusinessManagerCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	manager := GetBusinessManager(ctrl)
	user := GetTestUser(manager.baseServices.Database)
	bannerImage := GetTestFileMetadata(manager.baseServices.Database, user)
	manager.fileStorageService.(*MockFileStorageService).
		EXPECT().
		CreateStub(&user).
		Return(*bannerImage, nil)
	iconImage := GetTestFileMetadata(manager.baseServices.Database, user)
	manager.fileStorageService.(*MockFileStorageService).
		EXPECT().
		CreateStub(&user).
		Return(*iconImage, nil)
	details := BusinessDetails{
		Name:           "test business",
		Description:    "Description",
		Address:        "test address",
		GPSCoordinates: "+27.5916+086.5640+8850CRSWGS_84/",
		NIP:            "1234567890",
		KRS:            "1234567890",
		REGON:          "1234567890",
		OwnerName:      "test owner",
	}
	business, err := manager.Create(user, &details)
	require.NoErrorf(t, err, "manager create returned an error")
	if business == nil {
		t.Errorf("business is nil")
		return
	}
	assert.Truef(t, MatchEntities(details, business), "business details and entity do not match")
	var dbBusiness Business
	manager.baseServices.Database.Find(&dbBusiness, &Business{Model: gorm.Model{ID: business.ID}})
	assert.Truef(t, bannerImage.PublicId == dbBusiness.BannerImageId || bannerImage.PublicId == dbBusiness.IconImageId, "invalid banner image id")
	assert.Truef(t, iconImage.PublicId == dbBusiness.BannerImageId || iconImage.PublicId == dbBusiness.IconImageId, "invalid icon image id")
	assert.Equalf(t, dbBusiness.Name, business.Name, "business name does not match")
}

func TestBusinessManagerCreateAccountAlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	manager := GetBusinessManager(ctrl)
	user := GetTestUser(manager.baseServices.Database)
	business := GetTestBusiness(manager.baseServices.Database, user)
	details := BusinessDetails{
		Name:           "test business",
		Description:    "Description",
		Address:        "test address",
		GPSCoordinates: "+27.5916+086.5640+8850CRSWGS_84/",
		NIP:            "1234567890",
		KRS:            "1234567890",
		REGON:          "1234567890",
		OwnerName:      "test owner",
	}
	business, err := manager.Create(user, &details)

	require.Equalf(t, nil, business, "business is not nil")
	require.Equalf(t, BusinessAlreadyExists, err, "create err does not equal business already exists")
}

func TestBusinessManagerChangeDetails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	manager := GetBusinessManager(ctrl)
	user := GetTestUser(manager.baseServices.Database)
	business := GetTestBusiness(manager.baseServices.Database, user)
	details := ChangeableBusinessDetails{
		Name:        "new test business",
		Description: "new test description",
	}
	business, err := manager.ChangeDetails(business, &details)

	require.Equalf(t, nil, err, "BusinessManager.ChangeDetails returned an error")
	if business == nil {
		t.Errorf("business is nil")
		return
	}
	require.Equalf(t, details.Name, business.Name, "business name does not match")
	require.Equalf(t, details.Description, business.Description, "business name does not match")

	var dbBusiness Business
	manager.baseServices.Database.Find(&dbBusiness, &Business{Model: gorm.Model{ID: business.ID}})
	require.Equalf(t, details.Name, dbBusiness.Name, "business name does not match")
	require.Equalf(t, details.Description, dbBusiness.Description, "business name does not match")
}

func TestBusinessManagerSearchExisting(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	manager := GetBusinessManager(ctrl)
	user := GetTestUser(manager.baseServices.Database)
	business := GetTestBusiness(manager.baseServices.Database, user)
	result, err := manager.Search(business.Name, "", 0, 0, 5)
	require.Equalf(t, nil, err, "BusinessManager.Search returned an error")
	require.Equalf(t, 1, len(result), "BusinessManager.Search returned more or less than one result")
	require.Equalf(t, business.Name, result[0].Name, "BusinessManager.Search returned invalid busines")
	resultNone, errNone := manager.Search("no such business", "", 0, 0, 5)
	require.Equalf(t, nil, errNone, "BusinessManager.Search returned an error")
	require.Equalf(t, 0, len(resultNone), "BusinessManager.Search returned more than one result")
}
