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
		BaseServices{
			Logger:   log.Default(),
			Database: GetTestDatabase(),
		},
		NewMockFileStorageService(ctrl),
	}
}

func matchBusinessWithDetails(t *testing.T, details BusinessDetails, business Business) {
	require.Equal(t, details.Name, business.Name)
	require.Equal(t, details.Description, business.Description)
	require.Equal(t, details.Address, business.Address)
	require.Equal(t, details.GPSCoordinates, business.GPSCoordinates)
	require.Equal(t, details.NIP, business.NIP)
	require.Equal(t, details.KRS, business.KRS)
	require.Equal(t, details.REGON, business.REGON)
	require.Equal(t, details.OwnerName, business.OwnerName)
}

func TestBusinessManagerCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	manager := GetBusinessManager(ctrl)
	user := GetTestUser(manager.baseServices.Database)
	bannerImage := GetTestFileMetadata(manager.baseServices.Database, user)
	manager.fileStorageService.(*MockFileStorageService).
		EXPECT().
		CreateStub(user).
		Return(bannerImage, nil)
	iconImage := GetTestFileMetadata(manager.baseServices.Database, user)
	manager.fileStorageService.(*MockFileStorageService).
		EXPECT().
		CreateStub(user).
		Return(iconImage, nil)
	details := BusinessDetails{
		Name:           "test business",
		Description:    "Description",
		Address:        "test address",
		GPSCoordinates: FromCoords(27.5916, 086.5640),
		NIP:            "1234567890",
		KRS:            "1234567890",
		REGON:          "1234567890",
		OwnerName:      "test owner",
	}
	business, err := manager.Create(user, &details)
	require.NoErrorf(t, err, "manager create returned an error")
	require.NotNilf(t, business, "business shuold not be nil")

	matchBusinessWithDetails(t, details, *business)

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
		GPSCoordinates: FromCoords(27.5916, 086.5640),
		NIP:            "1234567890",
		KRS:            "1234567890",
		REGON:          "1234567890",
		OwnerName:      "test owner",
	}
	business, err := manager.Create(user, &details)

	require.Nilf(t, business, "business is not nil")
	require.Equalf(t, ErrBusinessAlreadyExists, err, "create err does not equal business already exists")
}

func TestBusinessManagerChangeDetails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	manager := GetBusinessManager(ctrl)
	user := GetTestUser(manager.baseServices.Database)
	business := GetTestBusiness(manager.baseServices.Database, user)
	details := ChangeableBusinessDetails{
		Name:        Ptr("new test business"),
		Description: Ptr("new test description"),
	}
	business, err := manager.ChangeDetails(business, &details)

	require.Nilf(t, err, "BusinessManager.ChangeDetails returned an error")
	if business == nil {
		t.Errorf("business is nil")
		return
	}
	require.Equalf(t, *details.Name, business.Name, "business name does not match")
	require.Equalf(t, *details.Description, business.Description, "business description does not match")

	var dbBusiness Business
	manager.baseServices.Database.Find(&dbBusiness, &Business{Model: gorm.Model{ID: business.ID}})
	require.Equalf(t, *details.Name, dbBusiness.Name, "business name does not match")
	require.Equalf(t, *details.Description, dbBusiness.Description, "business description does not match")
}

func TestBusinessManagerChangeDetailsEmptyName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	manager := GetBusinessManager(ctrl)
	user := GetTestUser(manager.baseServices.Database)
	business := GetTestBusiness(manager.baseServices.Database, user)
	oldName := business.Name
	details := ChangeableBusinessDetails{
		Description: Ptr("new test description"),
	}
	business, err := manager.ChangeDetails(business, &details)

	require.Nilf(t, err, "BusinessManager.ChangeDetails returned an error")
	if business == nil {
		t.Errorf("business is nil")
		return
	}
	require.Equalf(t, oldName, business.Name, "business name does not match")
	require.Equalf(t, *details.Description, business.Description, "business description does not match")

	var dbBusiness Business
	manager.baseServices.Database.Find(&dbBusiness, &Business{Model: gorm.Model{ID: business.ID}})
	require.Equalf(t, oldName, dbBusiness.Name, "business name does not match")
	require.Equalf(t, *details.Description, dbBusiness.Description, "business description does not match")
}

func TestBusinessManagerChangeDetailsEmptyDescription(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	manager := GetBusinessManager(ctrl)
	user := GetTestUser(manager.baseServices.Database)
	business := GetTestBusiness(manager.baseServices.Database, user)
	oldDescription := business.Description
	details := ChangeableBusinessDetails{
		Name: Ptr("new test name"),
	}
	business, err := manager.ChangeDetails(business, &details)

	require.Nilf(t, err, "BusinessManager.ChangeDetails returned an error")
	if business == nil {
		t.Errorf("business is nil")
		return
	}
	require.Equalf(t, *details.Name, business.Name, "business name does not match")
	require.Equalf(t, oldDescription, business.Description, "business description does not match")

	var dbBusiness Business
	manager.baseServices.Database.Find(&dbBusiness, &Business{Model: gorm.Model{ID: business.ID}})
	require.Equalf(t, *details.Name, dbBusiness.Name, "business name does not match")
	require.Equalf(t, oldDescription, dbBusiness.Description, "business description does not match")
}

func TestBusinessManagerAddMenuImage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	manager := GetBusinessManager(ctrl)
	user := GetTestUser(manager.baseServices.Database)
	business := GetTestBusiness(manager.baseServices.Database, user)

	menuImageMetadata := GetTestFileMetadata(manager.baseServices.Database, user)
	manager.fileStorageService.(*MockFileStorageService).
		EXPECT().
		CreateStub(user).
		Return(menuImageMetadata, nil)

	menuImage, err := manager.AddMenuImage(user, business)
	require.Nilf(t, err, "BusinessManager.AddMenuImage returned an error")
	require.NotNilf(t, menuImage, "BusinessManager.AddMenuImage returned a nil menuImage")
	require.Equalf(t, business.ID, menuImage.BusinessId,
		"BusinessManager.AddMenuImage returned a menuImage with invalid businessId")

	var dbMenuImage MenuImage
	tx := manager.baseServices.Database.First(&dbMenuImage, MenuImage{Model: gorm.Model{ID: menuImage.ID}})
	require.Nilf(t, tx.GetError(), "Database.First returned an error on MenuImage serach")
	require.Equalf(t, business.ID, dbMenuImage.BusinessId,
		"BusinessManager.AddMenuImage created a menuImage with invalid businessId")
}

func TestBusinessManagerAddMenuImageTooManyImages(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	manager := GetBusinessManager(ctrl)
	user := GetTestUser(manager.baseServices.Database)
	business := GetTestBusiness(manager.baseServices.Database, user)

	for i := 0; i != 11; i++ {
		if i != 11 {
			menuImageMetadata := GetTestFileMetadata(manager.baseServices.Database, user)
			manager.fileStorageService.(*MockFileStorageService).
				EXPECT().
				CreateStub(user).
				Return(menuImageMetadata, nil)
		}

		menuImage, err := manager.AddMenuImage(user, business)
		if i == 11 {
			require.ErrorAsf(t, ErrTooManyMenuImages, err,
				"BusinessManager.AddMenuImage did not return TooManyMenuImages")
		} else {
			require.Nilf(t, err, "BusinessManager.AddMenuImage returned an error")
			require.NotNilf(t, menuImage, "BusinessManager.AddMenuImage returned a nil menuImage")
			require.Equalf(t, business.ID, menuImage.BusinessId,
				"BusinessManager.AddMenuImage returned a menuImage with invalid businessId")
		}
	}
}

func TestBusinessManagerRemoveMenuImage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	manager := GetBusinessManager(ctrl)
	user := GetTestUser(manager.baseServices.Database)
	business := GetTestBusiness(manager.baseServices.Database, user)

	menuImageMetadata := GetTestFileMetadata(manager.baseServices.Database, user)
	manager.fileStorageService.(*MockFileStorageService).
		EXPECT().
		CreateStub(user).
		Return(menuImageMetadata, nil)

	menuImage, err := manager.AddMenuImage(user, business)
	require.Nilf(t, err, "BusinessManager.AddMenuImage returned an error")
	require.NotNilf(t, menuImage, "BusinessManager.AddMenuImage returned a nil menuImage")
	require.Equalf(t, business.ID, menuImage.BusinessId,
		"BusinessManager.AddMenuImage returned a menuImage with invalid businessId")

	var dbFileMetadata FileMetadata
	tx := manager.baseServices.Database.First(&dbFileMetadata, &FileMetadata{PublicId: menuImageMetadata.PublicId})
	require.Nilf(t, tx.GetError(), "manager.baseServices.Database.First(&dbFileMetadata returned an error")

	manager.fileStorageService.(*MockFileStorageService).
		EXPECT().
		RemoveMetadata(dbFileMetadata).
		Return(nil)

	err = manager.RemoveMenuImage(menuImage)
	require.Nilf(t, err, "BusinessManager.RemoveMenuImage returned an error")

	var dbMenuImage MenuImage
	tx = manager.baseServices.Database.First(&dbMenuImage, MenuImage{Model: gorm.Model{ID: menuImage.ID}})
	require.ErrorIsf(t, tx.GetError(), gorm.ErrRecordNotFound,
		"Database.First did not return ErrRecordNotFound on MenuImage serach")
}

func TestBusinessManagerSearchExistingByName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	manager := GetBusinessManager(ctrl)
	user := GetTestUser(manager.baseServices.Database)
	business := GetTestBusiness(manager.baseServices.Database, user)

	result, err := manager.Search(&business.Name, nil, 0, 0, 5)
	require.Nilf(t, err, "BusinessManager.Search returned an error")
	require.Equalf(t, 1, len(result), "BusinessManager.Search returned more or less than one result")
	require.Equalf(t, business.Name, result[0].Name, "BusinessManager.Search returned invalid busines")

	resultNone, errNone := manager.Search(Ptr("no such business"), nil, 0, 0, 5)
	require.Nilf(t, errNone, "BusinessManager.Search returned an error")
	require.Equalf(t, 0, len(resultNone), "BusinessManager.Search returned more than one result")
}

func TestBusinessManagerSearchExistingByLocation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	manager := GetBusinessManager(ctrl)
	user := GetTestUser(manager.baseServices.Database)
	business := GetTestBusiness(manager.baseServices.Database, user)

	result, err := manager.Search(nil, Ptr(FromCoords(27.59161, 086.56401)), 100, 0, 5)
	require.Nilf(t, err, "BusinessManager.Search returned an error")
	require.Equalf(t, 1, len(result), "BusinessManager.Search returned more or less than one result")
	require.Equalf(t, business.Name, result[0].Name, "BusinessManager.Search returned invalid busines")

	resultNone, errNone := manager.Search(nil, Ptr(FromCoords(27.69161, 086.16401)), 100, 0, 5)
	require.Nilf(t, errNone, "BusinessManager.Search returned an error")
	require.Equalf(t, 0, len(resultNone), "BusinessManager.Search returned more than one result")
}

func TestBusinessManagerSearchExistingByNameAndLocation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	manager := GetBusinessManager(ctrl)
	user := GetTestUser(manager.baseServices.Database)
	business := GetTestBusiness(manager.baseServices.Database, user)

	result, err := manager.Search(&business.Name, Ptr(FromCoords(27.59161, 086.56401)), 100, 0, 5)
	require.Nilf(t, err, "BusinessManager.Search returned an error")
	require.Equalf(t, 1, len(result), "BusinessManager.Search returned more or less than one result")
	require.Equalf(t, business.Name, result[0].Name, "BusinessManager.Search returned invalid busines")

	resultNone, errNone := manager.Search(&business.Name, Ptr(FromCoords(27.19161, 086.86401)), 100, 0, 5)
	require.Nilf(t, errNone, "BusinessManager.Search returned an error")
	require.Equalf(t, 0, len(resultNone), "BusinessManager.Search returned more than one result")

	resultNone, errNone = manager.Search(Ptr("invalid name"), Ptr(FromCoords(27.59161, 086.56401)), 100, 0, 5)
	require.Nilf(t, errNone, "BusinessManager.Search returned an error")
	require.Equalf(t, 0, len(resultNone), "BusinessManager.Search returned more than one result")
}
