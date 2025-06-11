package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"go-auth-api/internal/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type DeviceRepository struct {
	db *sqlx.DB
}

func NewDeviceRepository(db *sqlx.DB) *DeviceRepository {
	return &DeviceRepository{db: db}
}

// Device Version methods
func (r *DeviceRepository) CreateDeviceVersion(version *models.DeviceVersion) error {
	query := `
		INSERT INTO device_versions (name, version, description)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRow(query, version.Name, version.Version, version.Description).
		Scan(&version.ID, &version.CreatedAt, &version.UpdatedAt)
}

func (r *DeviceRepository) GetDeviceVersionByID(id uuid.UUID) (*models.DeviceVersion, error) {
	version := &models.DeviceVersion{}
	query := `SELECT id, name, version, description, created_at, updated_at 
			  FROM device_versions WHERE id = $1`

	err := r.db.Get(version, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("device version not found")
		}
		return nil, err
	}
	return version, nil
}

func (r *DeviceRepository) GetDeviceVersions(page, pageSize int) ([]models.DeviceVersion, int, error) {
	offset := (page - 1) * pageSize

	// Get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM device_versions`
	err := r.db.Get(&total, countQuery)
	if err != nil {
		return nil, 0, err
	}

	// Get versions
	query := `SELECT id, name, version, description, created_at, updated_at 
			  FROM device_versions 
			  ORDER BY created_at DESC 
			  LIMIT $1 OFFSET $2`

	var versions []models.DeviceVersion
	err = r.db.Select(&versions, query, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	return versions, total, nil
}

func (r *DeviceRepository) UpdateDeviceVersion(id uuid.UUID, req *models.UpdateDeviceVersionRequest) error {
	setParts := []string{}
	args := []interface{}{}
	argIndex := 1

	if req.Name != nil {
		setParts = append(setParts, fmt.Sprintf("name = $%d", argIndex))
		args = append(args, *req.Name)
		argIndex++
	}

	if req.Version != nil {
		setParts = append(setParts, fmt.Sprintf("version = $%d", argIndex))
		args = append(args, *req.Version)
		argIndex++
	}

	if req.Description != nil {
		setParts = append(setParts, fmt.Sprintf("description = $%d", argIndex))
		args = append(args, *req.Description)
		argIndex++
	}

	if len(setParts) == 0 {
		return fmt.Errorf("no fields to update")
	}

	setParts = append(setParts, fmt.Sprintf("updated_at = CURRENT_TIMESTAMP"))
	args = append(args, id)

	query := fmt.Sprintf("UPDATE device_versions SET %s WHERE id = $%d",
		strings.Join(setParts, ", "), argIndex)

	result, err := r.db.Exec(query, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("device version not found")
	}

	return nil
}

func (r *DeviceRepository) DeleteDeviceVersion(id uuid.UUID) error {
	query := `DELETE FROM device_versions WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("device version not found")
	}

	return nil
}

// Allowed Device methods
func (r *DeviceRepository) CreateAllowedDevice(device *models.AllowedDevice) error {
	query := `
		INSERT INTO allowed_devices (dev_eui, nwk_key, app_key, addr_key, description)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRow(query, device.DevEUI, device.NwkKey, device.AppKey, device.AddrKey, device.Description).
		Scan(&device.ID, &device.CreatedAt, &device.UpdatedAt)
}

func (r *DeviceRepository) GetAllowedDeviceByDevEUI(devEUI string) (*models.AllowedDevice, error) {
	device := &models.AllowedDevice{}
	query := `SELECT id, dev_eui, nwk_key, app_key, addr_key, description, created_at, updated_at 
			  FROM allowed_devices WHERE dev_eui = $1`

	err := r.db.Get(device, query, devEUI)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("allowed device not found")
		}
		return nil, err
	}
	return device, nil
}

func (r *DeviceRepository) GetAllowedDevices(page, pageSize int) ([]models.AllowedDevice, int, error) {
	offset := (page - 1) * pageSize

	// Get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM allowed_devices`
	err := r.db.Get(&total, countQuery)
	if err != nil {
		return nil, 0, err
	}

	// Get devices
	query := `SELECT id, dev_eui, nwk_key, app_key, addr_key, description, created_at, updated_at 
			  FROM allowed_devices 
			  ORDER BY created_at DESC 
			  LIMIT $1 OFFSET $2`

	var devices []models.AllowedDevice
	err = r.db.Select(&devices, query, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	return devices, total, nil
}

func (r *DeviceRepository) UpdateAllowedDevice(devEUI string, req *models.UpdateAllowedDeviceRequest) error {
	setParts := []string{}
	args := []interface{}{}
	argIndex := 1

	if req.NwkKey != nil {
		setParts = append(setParts, fmt.Sprintf("nwk_key = $%d", argIndex))
		args = append(args, *req.NwkKey)
		argIndex++
	}

	if req.AppKey != nil {
		setParts = append(setParts, fmt.Sprintf("app_key = $%d", argIndex))
		args = append(args, *req.AppKey)
		argIndex++
	}

	if req.AddrKey != nil {
		setParts = append(setParts, fmt.Sprintf("addr_key = $%d", argIndex))
		args = append(args, *req.AddrKey)
		argIndex++
	}

	if req.Description != nil {
		setParts = append(setParts, fmt.Sprintf("description = $%d", argIndex))
		args = append(args, *req.Description)
		argIndex++
	}

	if len(setParts) == 0 {
		return fmt.Errorf("no fields to update")
	}

	setParts = append(setParts, fmt.Sprintf("updated_at = CURRENT_TIMESTAMP"))
	args = append(args, devEUI)

	query := fmt.Sprintf("UPDATE allowed_devices SET %s WHERE dev_eui = $%d",
		strings.Join(setParts, ", "), argIndex)

	result, err := r.db.Exec(query, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("allowed device not found")
	}

	return nil
}

func (r *DeviceRepository) DeleteAllowedDevice(devEUI string) error {
	query := `DELETE FROM allowed_devices WHERE dev_eui = $1`
	result, err := r.db.Exec(query, devEUI)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("allowed device not found")
	}

	return nil
}

// Device methods
func (r *DeviceRepository) CreateDevice(device *models.Device) error {
	query := `
		INSERT INTO devices (user_id, version_id, name, dev_eui, description)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRow(query, device.UserID, device.VersionID, device.Name, device.DevEUI, device.Description).
		Scan(&device.ID, &device.CreatedAt, &device.UpdatedAt)
}

func (r *DeviceRepository) GetDeviceByID(id uuid.UUID) (*models.Device, error) {
	device := &models.Device{}
	query := `
		SELECT d.id, d.user_id, d.version_id, d.name, d.dev_eui, d.description,
			   d.chirpstack_device_created, d.chirpstack_device_activated, d.is_active,
			   d.created_at, d.updated_at,
			   dv.id as "version.id", dv.name as "version.name", dv.version as "version.version",
			   dv.description as "version.description", dv.created_at as "version.created_at", dv.updated_at as "version.updated_at"
		FROM devices d
		LEFT JOIN device_versions dv ON d.version_id = dv.id
		WHERE d.id = $1`

	err := r.db.Get(device, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("device not found")
		}
		return nil, err
	}
	return device, nil
}

func (r *DeviceRepository) GetDevicesByUserID(userID uuid.UUID, page, pageSize int) ([]models.Device, int, error) {
	offset := (page - 1) * pageSize

	// Get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM devices WHERE user_id = $1`
	err := r.db.Get(&total, countQuery, userID)
	if err != nil {
		return nil, 0, err
	}

	// Get devices
	query := `
		SELECT d.id, d.user_id, d.version_id, d.name, d.dev_eui, d.description,
			   d.chirpstack_device_created, d.chirpstack_device_activated, d.is_active,
			   d.created_at, d.updated_at,
			   dv.id as "version.id", dv.name as "version.name", dv.version as "version.version",
			   dv.description as "version.description", dv.created_at as "version.created_at", dv.updated_at as "version.updated_at"
		FROM devices d
		LEFT JOIN device_versions dv ON d.version_id = dv.id
		WHERE d.user_id = $1
		ORDER BY d.created_at DESC
		LIMIT $2 OFFSET $3`

	var devices []models.Device
	err = r.db.Select(&devices, query, userID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	return devices, total, nil
}

func (r *DeviceRepository) GetAllDevices(page, pageSize int) ([]models.Device, int, error) {
	offset := (page - 1) * pageSize

	// Get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM devices`
	err := r.db.Get(&total, countQuery)
	if err != nil {
		return nil, 0, err
	}

	// Get devices
	query := `
		SELECT d.id, d.user_id, d.version_id, d.name, d.dev_eui, d.description,
			   d.chirpstack_device_created, d.chirpstack_device_activated, d.is_active,
			   d.created_at, d.updated_at,
			   dv.id as "version.id", dv.name as "version.name", dv.version as "version.version",
			   dv.description as "version.description", dv.created_at as "version.created_at", dv.updated_at as "version.updated_at"
		FROM devices d
		LEFT JOIN device_versions dv ON d.version_id = dv.id
		ORDER BY d.created_at DESC
		LIMIT $1 OFFSET $2`

	var devices []models.Device
	err = r.db.Select(&devices, query, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	return devices, total, nil
}

func (r *DeviceRepository) UpdateDevice(id uuid.UUID, req *models.UpdateDeviceRequest) error {
	setParts := []string{}
	args := []interface{}{}
	argIndex := 1

	if req.Name != nil {
		setParts = append(setParts, fmt.Sprintf("name = $%d", argIndex))
		args = append(args, *req.Name)
		argIndex++
	}

	if req.VersionID != nil {
		setParts = append(setParts, fmt.Sprintf("version_id = $%d", argIndex))
		args = append(args, *req.VersionID)
		argIndex++
	}

	if req.Description != nil {
		setParts = append(setParts, fmt.Sprintf("description = $%d", argIndex))
		args = append(args, *req.Description)
		argIndex++
	}

	if req.IsActive != nil {
		setParts = append(setParts, fmt.Sprintf("is_active = $%d", argIndex))
		args = append(args, *req.IsActive)
		argIndex++
	}

	if len(setParts) == 0 {
		return fmt.Errorf("no fields to update")
	}

	setParts = append(setParts, fmt.Sprintf("updated_at = CURRENT_TIMESTAMP"))
	args = append(args, id)

	query := fmt.Sprintf("UPDATE devices SET %s WHERE id = $%d",
		strings.Join(setParts, ", "), argIndex)

	result, err := r.db.Exec(query, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("device not found")
	}

	return nil
}

func (r *DeviceRepository) UpdateDeviceChirpStackStatus(id uuid.UUID, created, activated bool) error {
	query := `UPDATE devices SET chirpstack_device_created = $1, chirpstack_device_activated = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3`
	result, err := r.db.Exec(query, created, activated, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("device not found")
	}

	return nil
}

func (r *DeviceRepository) DeleteDevice(id uuid.UUID) error {
	query := `DELETE FROM devices WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("device not found")
	}

	return nil
}
