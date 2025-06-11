package models

import (
	"time"

	"github.com/google/uuid"
)

// DeviceVersion represents a device version/model
type DeviceVersion struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Version     string    `json:"version" db:"version"`
	Description *string   `json:"description,omitempty" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// AllowedDevice represents pre-configured device keys
type AllowedDevice struct {
	ID          uuid.UUID `json:"id" db:"id"`
	DevEUI      string    `json:"dev_eui" db:"dev_eui"`
	NwkKey      string    `json:"nwk_key" db:"nwk_key"`
	AppKey      string    `json:"app_key" db:"app_key"`
	AddrKey     string    `json:"addr_key" db:"addr_key"`
	Description *string   `json:"description,omitempty" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Device represents a user's IoT device
type Device struct {
	ID                        uuid.UUID      `json:"id" db:"id"`
	UserID                    uuid.UUID      `json:"user_id" db:"user_id"`
	VersionID                 uuid.UUID      `json:"version_id" db:"version_id"`
	Name                      string         `json:"name" db:"name"`
	DevEUI                    string         `json:"dev_eui" db:"dev_eui"`
	Description               *string        `json:"description,omitempty" db:"description"`
	ChirpStackDeviceCreated   bool           `json:"chirpstack_device_created" db:"chirpstack_device_created"`
	ChirpStackDeviceActivated bool           `json:"chirpstack_device_activated" db:"chirpstack_device_activated"`
	IsActive                  bool           `json:"is_active" db:"is_active"`
	CreatedAt                 time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt                 time.Time      `json:"updated_at" db:"updated_at"`
	
	// Joined fields
	Version *DeviceVersion `json:"version,omitempty"`
	User    *User          `json:"user,omitempty"`
}

// Request/Response models
type CreateDeviceVersionRequest struct {
	Name        string  `json:"name" binding:"required"`
	Version     string  `json:"version" binding:"required"`
	Description *string `json:"description"`
}

type UpdateDeviceVersionRequest struct {
	Name        *string `json:"name"`
	Version     *string `json:"version"`
	Description *string `json:"description"`
}

type CreateAllowedDeviceRequest struct {
	DevEUI      string  `json:"dev_eui" binding:"required,len=16"`
	NwkKey      string  `json:"nwk_key" binding:"required,len=32"`
	AppKey      string  `json:"app_key" binding:"required,len=32"`
	AddrKey     string  `json:"addr_key" binding:"required,len=8"`
	Description *string `json:"description"`
}

type UpdateAllowedDeviceRequest struct {
	NwkKey      *string `json:"nwk_key"`
	AppKey      *string `json:"app_key"`
	AddrKey     *string `json:"addr_key"`
	Description *string `json:"description"`
}

type CreateDeviceRequest struct {
	Name        string    `json:"name" binding:"required"`
	VersionID   uuid.UUID `json:"version_id" binding:"required"`
	DevEUI      string    `json:"dev_eui" binding:"required,len=16"`
	Description *string   `json:"description"`
}

type UpdateDeviceRequest struct {
	Name        *string    `json:"name"`
	VersionID   *uuid.UUID `json:"version_id"`
	Description *string    `json:"description"`
	IsActive    *bool      `json:"is_active"`
}

type DeviceListResponse struct {
	Devices     []Device `json:"devices"`
	Total       int      `json:"total"`
	Page        int      `json:"page"`
	PageSize    int      `json:"page_size"`
	TotalPages  int      `json:"total_pages"`
}

type DeviceVersionListResponse struct {
	Versions    []DeviceVersion `json:"versions"`
	Total       int             `json:"total"`
	Page        int             `json:"page"`
	PageSize    int             `json:"page_size"`
	TotalPages  int             `json:"total_pages"`
}

type AllowedDeviceListResponse struct {
	Devices     []AllowedDevice `json:"devices"`
	Total       int             `json:"total"`
	Page        int             `json:"page"`
	PageSize    int             `json:"page_size"`
	TotalPages  int             `json:"total_pages"`
}

// ChirpStack API models
type ChirpStackCreateDeviceRequest struct {
	Device ChirpStackDeviceInfo `json:"device"`
}

type ChirpStackDeviceInfo struct {
	ApplicationID     string            `json:"applicationId"`
	Description       string            `json:"description"`
	DevEUI            string            `json:"devEui"`
	DeviceProfileID   string            `json:"deviceProfileId"`
	IsDisabled        bool              `json:"isDisabled"`
	JoinEUI           string            `json:"joinEui"`
	Name              string            `json:"name"`
	SkipFcntCheck     bool              `json:"skipFcntCheck"`
	Tags              map[string]string `json:"tags"`
	Variables         map[string]string `json:"variables"`
}

type ChirpStackActivateDeviceRequest struct {
	DeviceActivation ChirpStackDeviceActivation `json:"deviceActivation"`
}

type ChirpStackDeviceActivation struct {
	AFCntDown    int    `json:"aFCntDown"`
	AppSKey      string `json:"appSKey"`
	DevAddr      string `json:"devAddr"`
	FCntUp       int    `json:"fCntUp"`
	FNwkSIntKey  string `json:"fNwkSIntKey"`
	NFCntDown    int    `json:"nFCntDown"`
	NwkSEncKey   string `json:"nwkSEncKey"`
	SNwkSIntKey  string `json:"sNwkSIntKey"`
}
