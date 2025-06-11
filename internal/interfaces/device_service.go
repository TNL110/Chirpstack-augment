package interfaces

import (
	"go-auth-api/internal/models"
	"github.com/google/uuid"
)

type DeviceServiceInterface interface {
	// Device Version methods
	CreateDeviceVersion(req *models.CreateDeviceVersionRequest) (*models.DeviceVersion, error)
	GetDeviceVersionByID(id uuid.UUID) (*models.DeviceVersion, error)
	GetDeviceVersions(page, pageSize int) (*models.DeviceVersionListResponse, error)
	UpdateDeviceVersion(id uuid.UUID, req *models.UpdateDeviceVersionRequest) error
	DeleteDeviceVersion(id uuid.UUID) error

	// Allowed Device methods
	CreateAllowedDevice(req *models.CreateAllowedDeviceRequest) (*models.AllowedDevice, error)
	GetAllowedDeviceByDevEUI(devEUI string) (*models.AllowedDevice, error)
	GetAllowedDevices(page, pageSize int) (*models.AllowedDeviceListResponse, error)
	UpdateAllowedDevice(devEUI string, req *models.UpdateAllowedDeviceRequest) error
	DeleteAllowedDevice(devEUI string) error

	// Device methods
	CreateDevice(userID uuid.UUID, req *models.CreateDeviceRequest) (*models.Device, error)
	GetDeviceByID(id uuid.UUID) (*models.Device, error)
	GetDevicesByUserID(userID uuid.UUID, page, pageSize int) (*models.DeviceListResponse, error)
	GetAllDevices(page, pageSize int) (*models.DeviceListResponse, error)
	UpdateDevice(id uuid.UUID, req *models.UpdateDeviceRequest) error
	DeleteDevice(id uuid.UUID) error
}
