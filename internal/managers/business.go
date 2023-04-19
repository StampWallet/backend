package manager

import (
    . "github.com/StampWallet/backend/internal/database"
)

type BusinessManager interface {
    Create(user *User, businessDetails *BusinessDetails) (*Business, error)
    ChangeDetails(business *Business, businessDetails *BusinessDetails) (*Business, error)
    Search(name string, location string, proximityInMeters uint64) ([]Business, error) //? not a fan
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

type BusinessManagerImpl struct {
    baseServices *BaseServices
    fileStorageService *FileStorageService
}

func (manager *BusinessManagerImpl) Create(user *User, businessDetails *BusinessDetails) (Business, error) {

}

func (manager *BusinessManagerImpl) ChangeDetails(business *Business, businessDetails *BusinessDetails) (Business, error) {

}

//? not a fan
func (manager *BusinessManagerImpl) Search(name string, location string, proximityInMeters uint64) ([]Business, error) {

}
