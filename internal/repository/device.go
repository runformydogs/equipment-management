package repository

import (
	"equipment-management/internal/models"
	"gorm.io/gorm"
)

type DeviceRepository struct {
	db *gorm.DB
}

func NewDeviceRepository(db *gorm.DB) *DeviceRepository {
	return &DeviceRepository{db: db}
}

func (r *DeviceRepository) Create(device *models.Device) error {
	return r.db.Create(device).Error
}

func (r *DeviceRepository) GetByID(id uint) (*models.Device, error) {
	var device models.Device
	if err := r.db.First(&device, id).Error; err != nil {
		return nil, err
	}
	return &device, nil
}

func (r *DeviceRepository) Update(id uint, updateData *models.Device) (*models.Device, error) {
	var device models.Device
	if err := r.db.First(&device, id).Error; err != nil {
		return nil, err
	}

	if err := r.db.Model(&device).Updates(updateData).Error; err != nil {
		return nil, err
	}

	return &device, nil
}

func (r *DeviceRepository) Delete(id uint) error {
	return r.db.Delete(&models.Device{}, id).Error
}

func (r *DeviceRepository) GetAll() ([]models.Device, error) {
	var devices []models.Device
	if err := r.db.Find(&devices).Error; err != nil {
		return nil, err
	}
	return devices, nil
}
