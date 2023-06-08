package managers

import (
	"errors"
	"fmt"

	. "github.com/StampWallet/backend/internal/database"
	. "github.com/StampWallet/backend/internal/services"
	"github.com/lithammer/shortuuid/v4"
	"gorm.io/gorm"
)

var ErrBusinessAlreadyExists = errors.New("Business already exists")
var ErrTooManyMenuImages = errors.New("Too many menu images")

type BusinessManager interface {
	Create(user *User, businessDetails *BusinessDetails) (*Business, error)
	ChangeDetails(business *Business, businessDetails *ChangeableBusinessDetails) (*Business, error)
	AddMenuImage(business *Business) (*MenuImage, error)
	RemoveMenuImage(menuImage *MenuImage) error

	//? not a fan
	Search(name *string, location *GPSCoordinates, proximityInMeters uint, offset uint, limit uint) ([]Business, error)
	GetById(businessId string, preloadDetails bool) (*Business, error)
}

type BusinessDetails struct {
	Name           string
	Description    string
	Address        string
	GPSCoordinates GPSCoordinates
	NIP            string
	KRS            string
	REGON          string
	OwnerName      string
}

type ChangeableBusinessDetails struct {
	Name        *string
	Description *string
}

type BusinessManagerImpl struct {
	baseServices       BaseServices
	fileStorageService FileStorageService
}

func CreateBusinessManagerImpl(baseServices BaseServices, fileStorageService FileStorageService) *BusinessManagerImpl {
	return &BusinessManagerImpl{
		baseServices:       baseServices,
		fileStorageService: fileStorageService,
	}
}

func (manager *BusinessManagerImpl) Create(user *User, businessDetails *BusinessDetails) (*Business, error) {
	db := manager.baseServices.Database
	var business Business
	err := db.Transaction(func(tx GormDB) error {
		// somehow returns the first business found. does not have to belong to the user. how
		//r := tx.First(&business, &Business{User: user})
		r := tx.First(&business, &Business{OwnerId: user.ID})
		err := r.GetError()
		if err == nil {
			return ErrBusinessAlreadyExists
		} else if err != gorm.ErrRecordNotFound {
			return fmt.Errorf("tx.First returned an error: %+v", err)
		}

		bannerImageStub, err := manager.fileStorageService.CreateStub(user)
		if err != nil {
			return fmt.Errorf("fileStorageService.CreateStub for bannerImageStub returned an error: %+v", err)
		}
		iconImageStub, err := manager.fileStorageService.CreateStub(user)
		if err != nil {
			return fmt.Errorf("fileStorageService.CreateStub for iconImageStub returned an error: %+v", err)
		}

		business = Business{
			PublicId:       shortuuid.New(),
			Name:           businessDetails.Name,
			Description:    businessDetails.Description,
			Address:        businessDetails.Address,
			GPSCoordinates: businessDetails.GPSCoordinates,
			NIP:            businessDetails.NIP,
			KRS:            businessDetails.KRS,
			REGON:          businessDetails.REGON,
			OwnerName:      businessDetails.OwnerName,
			BannerImageId:  bannerImageStub.PublicId,
			IconImageId:    iconImageStub.PublicId,
			OwnerId:        user.ID,
		}

		r = tx.Create(&business)
		if err := r.GetError(); err != nil {
			if err == gorm.ErrDuplicatedKey {
				return ErrBusinessAlreadyExists
			} else {
				return fmt.Errorf("tx.Create(business) returned an error: %+v", err)
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return &business, nil
}

func (manager *BusinessManagerImpl) ChangeDetails(business *Business, businessDetails *ChangeableBusinessDetails) (*Business, error) {
	if businessDetails.Name != nil {
		business.Name = *businessDetails.Name
	}
	if businessDetails.Description != nil {
		business.Description = *businessDetails.Description
	}

	tx := manager.baseServices.Database.Save(business)
	if err := tx.GetError(); err != nil {
		return nil, err
	}

	return business, nil
}

func (manager *BusinessManagerImpl) AddMenuImage(business *Business) (*MenuImage, error) {
	var menuImage *MenuImage
	err := manager.baseServices.Database.Transaction(func(db GormDB) error {
		var images []MenuImage
		tx := db.Find(&images, MenuImage{BusinessId: business.ID})
		if err := tx.GetError(); err != nil {
			return fmt.Errorf("database error when searching menuImages: %w", err)
		}
		if len(images) > 10 {
			return ErrTooManyMenuImages
		}

		metadata, err := manager.fileStorageService.CreateStub(business.User)
		if err != nil {
			return fmt.Errorf("failed to create image stub: %w", err)
		}

		menuImage = &MenuImage{
			BusinessId: business.ID,
			FileId:     metadata.PublicId,
		}
		tx = db.Save(menuImage)
		if err := tx.GetError(); err != nil {
			return fmt.Errorf("database error when saving menuImage: %w", err)
		}

		return nil
	})

	return menuImage, err
}

func (manager *BusinessManagerImpl) RemoveMenuImage(menuImage *MenuImage) error {
	tx := manager.baseServices.Database.Delete(menuImage)
	if err := tx.GetError(); err != nil {
		return fmt.Errorf("database error when removing menuImage: %w", err)
	}

	return nil
}

// NOTE limit offset is not a very good pagination method
// https://www.citusdata.com/blog/2016/03/30/five-ways-to-paginate/
func (manager *BusinessManagerImpl) Search(name *string, location *GPSCoordinates, proximityInMeters uint, offset uint, limit uint) ([]Business, error) {
	var args []interface{}
	//not a fan of constructing sql queries, but it should be safe this time. query is constructed only
	//from constants, args are passed in parameters
	query := "SELECT * FROM businesses WHERE "
	if name != nil {
		query += `to_tsvector('simple', f_concat_ws(' ', name, description, address))
		@@ plainto_tsquery('simple', ?)`
		if location != nil {
			query += " AND "
		}
		args = append(args, name)
	}
	if location != nil {
		query += `ST_DWithin(gps_coordinates, ?, ?)`
		args = append(args, location, proximityInMeters)
	}
	query += ` LIMIT ? OFFSET ?`
	args = append(args, limit, offset)

	var businesses []Business
	result := manager.baseServices.Database.Raw(query, args...).Scan(&businesses)
	if err := result.GetError(); err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return businesses, nil
}

func (manager *BusinessManagerImpl) GetById(businessId string, preloadDetails bool) (*Business, error) {
	var business Business
	db := manager.baseServices.Database

	if preloadDetails {
		db = db.Preload("ItemDefinitions").Preload("MenuImages")
	}

	r := db.First(&business, &Business{PublicId: businessId})
	err := r.GetError()
	if err == nil {
		return nil, ErrNoSuchBusiness
	} else if err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("tx.First returned an error: %+v", err)
	}

	return &business, nil
}
