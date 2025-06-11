package service

import (
	"fmt"
	"math"

	"go-auth-api/internal/models"
	"go-auth-api/internal/repository"

	"github.com/google/uuid"
)

type DeviceService struct {
	deviceRepo        *repository.DeviceRepository
	userRepo          *repository.UserRepository
	chirpStackService *ChirpStackService
}

func NewDeviceService(deviceRepo *repository.DeviceRepository, userRepo *repository.UserRepository, chirpStackService *ChirpStackService) *DeviceService {
	return &DeviceService{
		deviceRepo:        deviceRepo,
		userRepo:          userRepo,
		chirpStackService: chirpStackService,
	}
}

// Device Version methods
func (s *DeviceService) CreateDeviceVersion(req *models.CreateDeviceVersionRequest) (*models.DeviceVersion, error) {
	version := &models.DeviceVersion{
		Name:        req.Name,
		Version:     req.Version,
		Description: req.Description,
	}

	err := s.deviceRepo.CreateDeviceVersion(version)
	if err != nil {
		return nil, fmt.Errorf("failed to create device version: %w", err)
	}

	return version, nil
}

func (s *DeviceService) GetDeviceVersionByID(id uuid.UUID) (*models.DeviceVersion, error) {
	return s.deviceRepo.GetDeviceVersionByID(id)
}

func (s *DeviceService) GetDeviceVersions(page, pageSize int) (*models.DeviceVersionListResponse, error) {
	versions, total, err := s.deviceRepo.GetDeviceVersions(page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("failed to get device versions: %w", err)
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return &models.DeviceVersionListResponse{
		Versions:   versions,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *DeviceService) UpdateDeviceVersion(id uuid.UUID, req *models.UpdateDeviceVersionRequest) error {
	return s.deviceRepo.UpdateDeviceVersion(id, req)
}

func (s *DeviceService) DeleteDeviceVersion(id uuid.UUID) error {
	return s.deviceRepo.DeleteDeviceVersion(id)
}

// Allowed Device methods
func (s *DeviceService) CreateAllowedDevice(req *models.CreateAllowedDeviceRequest) (*models.AllowedDevice, error) {
	device := &models.AllowedDevice{
		DevEUI:      req.DevEUI,
		NwkKey:      req.NwkKey,
		AppKey:      req.AppKey,
		AddrKey:     req.AddrKey,
		Description: req.Description,
	}

	err := s.deviceRepo.CreateAllowedDevice(device)
	if err != nil {
		return nil, fmt.Errorf("failed to create allowed device: %w", err)
	}

	return device, nil
}

func (s *DeviceService) GetAllowedDeviceByDevEUI(devEUI string) (*models.AllowedDevice, error) {
	return s.deviceRepo.GetAllowedDeviceByDevEUI(devEUI)
}

func (s *DeviceService) GetAllowedDevices(page, pageSize int) (*models.AllowedDeviceListResponse, error) {
	devices, total, err := s.deviceRepo.GetAllowedDevices(page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("failed to get allowed devices: %w", err)
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return &models.AllowedDeviceListResponse{
		Devices:    devices,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *DeviceService) UpdateAllowedDevice(devEUI string, req *models.UpdateAllowedDeviceRequest) error {
	return s.deviceRepo.UpdateAllowedDevice(devEUI, req)
}

func (s *DeviceService) DeleteAllowedDevice(devEUI string) error {
	return s.deviceRepo.DeleteAllowedDevice(devEUI)
}

// Device methods
func (s *DeviceService) CreateDevice(userID uuid.UUID, req *models.CreateDeviceRequest) (*models.Device, error) {
	// Check if user exists and has ChirpStack data
	user, err := s.userRepo.GetUserByID(userID.String())
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	if user.ApplicationID == nil || user.DeviceProfileID == nil {
		return nil, fmt.Errorf("user does not have ChirpStack application or device profile")
	}

	// Check if devEUI exists in allowed devices
	allowedDevice, err := s.deviceRepo.GetAllowedDeviceByDevEUI(req.DevEUI)
	if err != nil {
		return nil, fmt.Errorf("device with devEUI %s not found in allowed devices", req.DevEUI)
	}

	// Check if version exists
	_, err = s.deviceRepo.GetDeviceVersionByID(req.VersionID)
	if err != nil {
		return nil, fmt.Errorf("device version not found: %w", err)
	}

	// Create device in database
	device := &models.Device{
		UserID:      userID,
		VersionID:   req.VersionID,
		Name:        req.Name,
		DevEUI:      req.DevEUI,
		Description: req.Description,
	}

	err = s.deviceRepo.CreateDevice(device)
	if err != nil {
		return nil, fmt.Errorf("failed to create device: %w", err)
	}

	// Create device in ChirpStack if service is enabled
	if s.chirpStackService != nil && s.chirpStackService.IsEnabled() {
		err = s.createChirpStackDevice(device, user, allowedDevice)
		if err != nil {
			fmt.Printf("Warning: Failed to create ChirpStack device: %v\n", err)
		}
	}

	// Get the created device with version info
	return s.deviceRepo.GetDeviceByID(device.ID)
}

func (s *DeviceService) createChirpStackDevice(device *models.Device, user *models.User, allowedDevice *models.AllowedDevice) error {
	// Create device in ChirpStack
	createReq := models.ChirpStackCreateDeviceRequest{
		Device: models.ChirpStackDeviceInfo{
			ApplicationID:   *user.ApplicationID,
			Description:     device.Name,
			DevEUI:          device.DevEUI,
			DeviceProfileID: *user.DeviceProfileID,
			IsDisabled:      false,
			JoinEUI:         "0000000000000000",
			Name:            device.Name,
			SkipFcntCheck:   true,
			Tags:            make(map[string]string),
			Variables:       make(map[string]string),
		},
	}

	responseBody, err := s.chirpStackService.makeRequest("POST", "/devices", createReq)
	if err != nil {
		return fmt.Errorf("failed to create ChirpStack device: %w", err)
	}

	fmt.Printf("ChirpStack device created: %s\n", string(responseBody))

	// Activate device in ChirpStack
	activateReq := models.ChirpStackActivateDeviceRequest{
		DeviceActivation: models.ChirpStackDeviceActivation{
			AFCntDown:   0,
			AppSKey:     allowedDevice.AppKey,
			DevAddr:     allowedDevice.AddrKey,
			FCntUp:      0,
			FNwkSIntKey: allowedDevice.NwkKey,
			NFCntDown:   0,
			NwkSEncKey:  allowedDevice.NwkKey,
			SNwkSIntKey: allowedDevice.NwkKey,
		},
	}

	activateURL := fmt.Sprintf("/devices/%s/activate", device.DevEUI)
	responseBody, err = s.chirpStackService.makeRequest("POST", activateURL, activateReq)
	if err != nil {
		return fmt.Errorf("failed to activate ChirpStack device: %w", err)
	}

	fmt.Printf("ChirpStack device activated: %s\n", string(responseBody))

	// Update device status in database
	err = s.deviceRepo.UpdateDeviceChirpStackStatus(device.ID, true, true)
	if err != nil {
		fmt.Printf("Warning: Failed to update device ChirpStack status: %v\n", err)
	}

	return nil
}

func (s *DeviceService) GetDeviceByID(id uuid.UUID) (*models.Device, error) {
	return s.deviceRepo.GetDeviceByID(id)
}

func (s *DeviceService) GetDevicesByUserID(userID uuid.UUID, page, pageSize int) (*models.DeviceListResponse, error) {
	devices, total, err := s.deviceRepo.GetDevicesByUserID(userID, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("failed to get user devices: %w", err)
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return &models.DeviceListResponse{
		Devices:    devices,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *DeviceService) GetAllDevices(page, pageSize int) (*models.DeviceListResponse, error) {
	devices, total, err := s.deviceRepo.GetAllDevices(page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("failed to get all devices: %w", err)
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return &models.DeviceListResponse{
		Devices:    devices,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *DeviceService) UpdateDevice(id uuid.UUID, req *models.UpdateDeviceRequest) error {
	// Check if version exists if version_id is being updated
	if req.VersionID != nil {
		_, err := s.deviceRepo.GetDeviceVersionByID(*req.VersionID)
		if err != nil {
			return fmt.Errorf("device version not found: %w", err)
		}
	}

	return s.deviceRepo.UpdateDevice(id, req)
}

func (s *DeviceService) DeleteDevice(id uuid.UUID) error {
	// Get device info before deleting
	device, err := s.deviceRepo.GetDeviceByID(id)
	if err != nil {
		return fmt.Errorf("device not found: %w", err)
	}

	// Delete device from ChirpStack if service is enabled and device was created in ChirpStack
	if s.chirpStackService != nil && s.chirpStackService.IsEnabled() && device.ChirpStackDeviceCreated {
		err = s.chirpStackService.DeleteDevice(device.DevEUI)
		if err != nil {
			fmt.Printf("Warning: Failed to delete ChirpStack device: %v\n", err)
			// Continue with database deletion even if ChirpStack deletion fails
		}
	}

	// Delete device from database
	return s.deviceRepo.DeleteDevice(id)
}
