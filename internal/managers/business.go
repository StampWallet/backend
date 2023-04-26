package managers

import (
    "errors"

    . "github.com/StampWallet/backend/internal/database"
    . "github.com/StampWallet/backend/internal/services"
)

var BusinessAlreadyExists = errors.New("business already exists")

type BusinessManager interface {
    Create(user *User, businessDetails *BusinessDetails) (*Business, error)
    ChangeDetails(business *Business, businessDetails *BusinessDetails) (*Business, error)
    Search(name string, location string, proximityInMeters uint, offset uint, limit uint) ([]Business, error) //? not a fan
}

type BusinessDetails struct {
    Name string
    Description string
    Address string
    GPSCoordinates string
    NIP string
    KRS string
    REGON string
    OwnerName string
}

type ChangeableBusinessDetails struct {
    Name string
    Description string
}

type BusinessManagerImpl struct {
    baseServices *BaseServices
    fileStorageService FileStorageService
}

func (manager *BusinessManagerImpl) Create(user *User, businessDetails *BusinessDetails) (*Business, error) {
    return nil, nil
}

func (manager *BusinessManagerImpl) ChangeDetails(business *Business, businessDetails *ChangeableBusinessDetails ) (*Business, error) {
    return nil, nil
}

//? not a fan
func (manager *BusinessManagerImpl) Search(name string, location string, proximityInMeters uint, offset uint, limit uint) ([]Business, error) {
    return nil, nil
}
