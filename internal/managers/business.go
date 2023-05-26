package managers

import (
	"errors"
	"fmt"

	. "github.com/StampWallet/backend/internal/database"
	. "github.com/StampWallet/backend/internal/services"
	"github.com/lithammer/shortuuid/v4"
	"gorm.io/gorm"
)

var BusinessAlreadyExists = errors.New("Business already exists")

type BusinessManager interface {
	Create(user *User, businessDetails *BusinessDetails) (*Business, error)
	ChangeDetails(business *Business, businessDetails *ChangeableBusinessDetails) (*Business, error)
	Search(name *string, location *GPSCoordinates, proximityInMeters uint, offset uint, limit uint) ([]Business, error) //? not a fan
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
		r := tx.First(&business, Business{User: user})
		err := r.GetError()
		if err == nil {
			return BusinessAlreadyExists
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
			return fmt.Errorf("tx.Create(business) returned an error: %+v", err)
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
