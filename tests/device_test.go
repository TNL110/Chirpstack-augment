package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-auth-api/internal/handlers"
	"go-auth-api/internal/interfaces"
	"go-auth-api/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock DeviceService
type MockDeviceService struct {
	mock.Mock
}

// Implement DeviceServiceInterface
var _ interfaces.DeviceServiceInterface = (*MockDeviceService)(nil)

func (m *MockDeviceService) CreateDeviceVersion(req *models.CreateDeviceVersionRequest) (*models.DeviceVersion, error) {
	args := m.Called(req)
	return args.Get(0).(*models.DeviceVersion), args.Error(1)
}

func (m *MockDeviceService) GetDeviceVersionByID(id uuid.UUID) (*models.DeviceVersion, error) {
	args := m.Called(id)
	return args.Get(0).(*models.DeviceVersion), args.Error(1)
}

func (m *MockDeviceService) GetDeviceVersions(page, pageSize int) (*models.DeviceVersionListResponse, error) {
	args := m.Called(page, pageSize)
	return args.Get(0).(*models.DeviceVersionListResponse), args.Error(1)
}

func (m *MockDeviceService) UpdateDeviceVersion(id uuid.UUID, req *models.UpdateDeviceVersionRequest) error {
	args := m.Called(id, req)
	return args.Error(0)
}

func (m *MockDeviceService) DeleteDeviceVersion(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockDeviceService) CreateAllowedDevice(req *models.CreateAllowedDeviceRequest) (*models.AllowedDevice, error) {
	args := m.Called(req)
	return args.Get(0).(*models.AllowedDevice), args.Error(1)
}

func (m *MockDeviceService) GetAllowedDeviceByDevEUI(devEUI string) (*models.AllowedDevice, error) {
	args := m.Called(devEUI)
	return args.Get(0).(*models.AllowedDevice), args.Error(1)
}

func (m *MockDeviceService) GetAllowedDevices(page, pageSize int) (*models.AllowedDeviceListResponse, error) {
	args := m.Called(page, pageSize)
	return args.Get(0).(*models.AllowedDeviceListResponse), args.Error(1)
}

func (m *MockDeviceService) UpdateAllowedDevice(devEUI string, req *models.UpdateAllowedDeviceRequest) error {
	args := m.Called(devEUI, req)
	return args.Error(0)
}

func (m *MockDeviceService) DeleteAllowedDevice(devEUI string) error {
	args := m.Called(devEUI)
	return args.Error(0)
}

func (m *MockDeviceService) CreateDevice(userID uuid.UUID, req *models.CreateDeviceRequest) (*models.Device, error) {
	args := m.Called(userID, req)
	return args.Get(0).(*models.Device), args.Error(1)
}

func (m *MockDeviceService) GetDeviceByID(id uuid.UUID) (*models.Device, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Device), args.Error(1)
}

func (m *MockDeviceService) GetDevicesByUserID(userID uuid.UUID, page, pageSize int) (*models.DeviceListResponse, error) {
	args := m.Called(userID, page, pageSize)
	return args.Get(0).(*models.DeviceListResponse), args.Error(1)
}

func (m *MockDeviceService) GetAllDevices(page, pageSize int) (*models.DeviceListResponse, error) {
	args := m.Called(page, pageSize)
	return args.Get(0).(*models.DeviceListResponse), args.Error(1)
}

func (m *MockDeviceService) UpdateDevice(id uuid.UUID, req *models.UpdateDeviceRequest) error {
	args := m.Called(id, req)
	return args.Error(0)
}

func (m *MockDeviceService) DeleteDevice(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func setupDeviceTestRouter() (*gin.Engine, *MockDeviceService) {
	gin.SetMode(gin.TestMode)

	mockDeviceService := &MockDeviceService{}
	deviceHandler := handlers.NewDeviceHandler(mockDeviceService)

	router := gin.New()

	// Add middleware to set user_id for testing
	router.Use(func(c *gin.Context) {
		c.Set("user_id", uuid.New())
		c.Next()
	})

	api := router.Group("/api/v1/devices")
	{
		// Device version routes
		api.POST("/versions", deviceHandler.CreateDeviceVersion)
		api.GET("/versions", deviceHandler.GetDeviceVersions)
		api.GET("/versions/:id", deviceHandler.GetDeviceVersionByID)
		api.PUT("/versions/:id", deviceHandler.UpdateDeviceVersion)
		api.DELETE("/versions/:id", deviceHandler.DeleteDeviceVersion)

		// Allowed device routes
		api.POST("/allowed", deviceHandler.CreateAllowedDevice)
		api.GET("/allowed", deviceHandler.GetAllowedDevices)
		api.GET("/allowed/:devEUI", deviceHandler.GetAllowedDeviceByDevEUI)
		api.PUT("/allowed/:devEUI", deviceHandler.UpdateAllowedDevice)
		api.DELETE("/allowed/:devEUI", deviceHandler.DeleteAllowedDevice)

		// Device routes
		api.POST("", deviceHandler.CreateDevice)
		api.GET("/my", deviceHandler.GetMyDevices)
		api.GET("/all", deviceHandler.GetAllDevices)
		api.GET("/:id", deviceHandler.GetDeviceByID)
		api.PUT("/:id", deviceHandler.UpdateDevice)
		api.DELETE("/:id", deviceHandler.DeleteDevice)
	}

	return router, mockDeviceService
}

func TestCreateDeviceVersion(t *testing.T) {
	router, mockService := setupDeviceTestRouter()

	t.Run("Successful Creation", func(t *testing.T) {
		versionID := uuid.New()
		expectedVersion := &models.DeviceVersion{
			ID:          versionID,
			Name:        "RAK7200",
			Version:     "v1.0",
			Description: stringPtr("Test device version"),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		mockService.On("CreateDeviceVersion", mock.AnythingOfType("*models.CreateDeviceVersionRequest")).Return(expectedVersion, nil)

		reqBody := models.CreateDeviceVersionRequest{
			Name:        "RAK7200",
			Version:     "v1.0",
			Description: stringPtr("Test device version"),
		}

		jsonBody, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/api/v1/devices/versions", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response models.DeviceVersion
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "RAK7200", response.Name)
		assert.Equal(t, "v1.0", response.Version)

		mockService.AssertExpectations(t)
	})
}

func TestGetDeviceVersions(t *testing.T) {
	router, mockService := setupDeviceTestRouter()

	t.Run("Successful Get", func(t *testing.T) {
		expectedResponse := &models.DeviceVersionListResponse{
			Versions: []models.DeviceVersion{
				{
					ID:      uuid.New(),
					Name:    "RAK7200",
					Version: "v1.0",
				},
			},
			Total:      1,
			Page:       1,
			PageSize:   10,
			TotalPages: 1,
		}

		mockService.On("GetDeviceVersions", 1, 10).Return(expectedResponse, nil)

		req, _ := http.NewRequest("GET", "/api/v1/devices/versions", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.DeviceVersionListResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 1, response.Total)
		assert.Len(t, response.Versions, 1)

		mockService.AssertExpectations(t)
	})
}

func TestCreateAllowedDevice(t *testing.T) {
	router, mockService := setupDeviceTestRouter()

	t.Run("Successful Creation", func(t *testing.T) {
		deviceID := uuid.New()
		expectedDevice := &models.AllowedDevice{
			ID:          deviceID,
			DevEUI:      "C5EABC521E8304EE",
			NwkKey:      "C518B15AB390B01762E4A3730E8C5F1C",
			AppKey:      "97784F3B7F2A57EECF19F10E625081E0",
			AddrKey:     "2F972E56",
			Description: stringPtr("Test device"),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		mockService.On("CreateAllowedDevice", mock.AnythingOfType("*models.CreateAllowedDeviceRequest")).Return(expectedDevice, nil)

		reqBody := models.CreateAllowedDeviceRequest{
			DevEUI:      "C5EABC521E8304EE",
			NwkKey:      "C518B15AB390B01762E4A3730E8C5F1C",
			AppKey:      "97784F3B7F2A57EECF19F10E625081E0",
			AddrKey:     "2F972E56",
			Description: stringPtr("Test device"),
		}

		jsonBody, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/api/v1/devices/allowed", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response models.AllowedDevice
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "C5EABC521E8304EE", response.DevEUI)

		mockService.AssertExpectations(t)
	})
}

func TestCreateDevice(t *testing.T) {
	router, mockService := setupDeviceTestRouter()

	t.Run("Successful Creation", func(t *testing.T) {
		deviceID := uuid.New()
		userID := uuid.New()
		versionID := uuid.New()

		expectedDevice := &models.Device{
			ID:                        deviceID,
			UserID:                    userID,
			VersionID:                 versionID,
			Name:                      "My Device",
			DevEUI:                    "C5EABC521E8304EE",
			Description:               stringPtr("Test device"),
			ChirpStackDeviceCreated:   true,
			ChirpStackDeviceActivated: true,
			IsActive:                  true,
			CreatedAt:                 time.Now(),
			UpdatedAt:                 time.Now(),
		}

		mockService.On("CreateDevice", mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("*models.CreateDeviceRequest")).Return(expectedDevice, nil)

		reqBody := models.CreateDeviceRequest{
			VersionID:   versionID,
			Name:        "My Device",
			DevEUI:      "C5EABC521E8304EE",
			Description: stringPtr("Test device"),
		}

		jsonBody, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/api/v1/devices", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response models.Device
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "My Device", response.Name)
		assert.True(t, response.ChirpStackDeviceCreated)
		assert.True(t, response.ChirpStackDeviceActivated)

		mockService.AssertExpectations(t)
	})
}

func TestDeleteDeviceWithChirpStack(t *testing.T) {
	router, mockService := setupDeviceTestRouter()

	t.Run("Successful Deletion with ChirpStack", func(t *testing.T) {
		deviceID := uuid.New()

		mockService.On("DeleteDevice", deviceID).Return(nil)

		req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/v1/devices/%s", deviceID), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Device deleted successfully", response["message"])

		mockService.AssertExpectations(t)
	})
}
