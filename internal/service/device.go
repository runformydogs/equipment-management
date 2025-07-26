package service

import (
	"equipment-management/internal/dto"
	"equipment-management/internal/models"
	"equipment-management/internal/repository"
	"time"
)

type DeviceService struct {
	repo *repository.DeviceRepository
}

func NewDeviceService(repo *repository.DeviceRepository) *DeviceService {
	return &DeviceService{repo: repo}
}

func (s *DeviceService) CreateDevice(req *dto.CreateDeviceRequest) (*models.Device, error) {
	device := models.Device{
		Type:          req.Type,
		Vendor:        req.Vendor,
		Model:         req.Model,
		Serial:        req.Serial,
		Location:      req.Location,
		NetworkNodeID: req.NetworkNodeID,
		Status:        "active",
	}

	if err := s.repo.Create(&device); err != nil {
		return nil, err
	}

	return &device, nil
}

func (s *DeviceService) GetDevice(id uint) (*models.Device, error) {
	return s.repo.GetByID(id)
}

func (s *DeviceService) UpdateDevice(id uint, req *dto.UpdateDeviceRequest) (*models.Device, error) {
	updateData := models.Device{
		Type:          req.Type,
		Vendor:        req.Vendor,
		Model:         req.Model,
		Serial:        req.Serial,
		Location:      req.Location,
		Status:        req.Status,
		NetworkNodeID: req.NetworkNodeID,
	}

	return s.repo.Update(id, &updateData)
}

func (s *DeviceService) DeleteDevice(id uint) error {
	return s.repo.Delete(id)
}

func (s *DeviceService) GetAllDevices() ([]models.Device, error) {
	return s.repo.GetAll()
}

func (s *DeviceService) ToDeviceResponse(device *models.Device) dto.DeviceResponse {
	return dto.DeviceResponse{
		ID:            device.ID,
		Type:          device.Type,
		Vendor:        device.Vendor,
		Model:         device.Model,
		Serial:        device.Serial,
		Location:      device.Location,
		Status:        device.Status,
		NetworkNodeID: device.NetworkNodeID,
		CreatedAt:     device.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     device.UpdatedAt.Format(time.RFC3339),
	}
}
