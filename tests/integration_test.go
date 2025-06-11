package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"go-auth-api/internal/models"

	"github.com/stretchr/testify/suite"
)

type IntegrationTestSuite struct {
	suite.Suite
	baseURL string
	token   string
	userID  string
}

func (suite *IntegrationTestSuite) SetupSuite() {
	suite.baseURL = "http://localhost:8080/api/v1"

	// Wait for service to be ready
	suite.waitForService()

	// Register and login a test user
	suite.registerAndLogin()
}

func (suite *IntegrationTestSuite) waitForService() {
	maxRetries := 30
	for i := 0; i < maxRetries; i++ {
		resp, err := http.Get("http://localhost:8080/health")
		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			return
		}
		if resp != nil {
			resp.Body.Close()
		}
		time.Sleep(2 * time.Second)
	}
	suite.T().Fatal("Service not ready after 60 seconds")
}

func (suite *IntegrationTestSuite) registerAndLogin() {
	// Register user
	registerReq := models.RegisterRequest{
		Email:    fmt.Sprintf("integration_test_%d@example.com", time.Now().Unix()),
		Password: "password123",
		FullName: "Integration Test User",
	}

	jsonBody, _ := json.Marshal(registerReq)
	resp, err := http.Post(suite.baseURL+"/auth/register", "application/json", bytes.NewBuffer(jsonBody))
	suite.Require().NoError(err)
	suite.Require().Equal(http.StatusCreated, resp.StatusCode)

	var authResp models.AuthResponse
	err = json.NewDecoder(resp.Body).Decode(&authResp)
	suite.Require().NoError(err)
	resp.Body.Close()

	suite.token = authResp.Token
	suite.userID = authResp.User.ID.String()
}

func (suite *IntegrationTestSuite) makeAuthenticatedRequest(method, endpoint string, body interface{}) (*http.Response, error) {
	var reqBody *bytes.Buffer
	if body != nil {
		jsonBody, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(jsonBody)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	req, err := http.NewRequest(method, suite.baseURL+endpoint, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.token)

	client := &http.Client{}
	return client.Do(req)
}

func (suite *IntegrationTestSuite) TestCompleteDeviceWorkflow() {
	// 1. Create device version
	versionReq := models.CreateDeviceVersionRequest{
		Name:        "RAK7200",
		Version:     "v1.0",
		Description: stringPtr("RAK7200 LoRaWAN Tracker v1.0"),
	}

	resp, err := suite.makeAuthenticatedRequest("POST", "/devices/versions", versionReq)
	suite.Require().NoError(err)
	suite.Equal(http.StatusCreated, resp.StatusCode)

	var version models.DeviceVersion
	err = json.NewDecoder(resp.Body).Decode(&version)
	suite.Require().NoError(err)
	resp.Body.Close()

	suite.Equal("RAK7200", version.Name)
	suite.Equal("v1.0", version.Version)

	// 2. Create allowed device
	allowedReq := models.CreateAllowedDeviceRequest{
		DevEUI:      "C5EABC521E8304EE",
		NwkKey:      "C518B15AB390B01762E4A3730E8C5F1C",
		AppKey:      "97784F3B7F2A57EECF19F10E625081E0",
		AddrKey:     "2F972E56",
		Description: stringPtr("Integration test device"),
	}

	resp, err = suite.makeAuthenticatedRequest("POST", "/devices/allowed", allowedReq)
	suite.Require().NoError(err)
	suite.Equal(http.StatusCreated, resp.StatusCode)

	var allowedDevice models.AllowedDevice
	err = json.NewDecoder(resp.Body).Decode(&allowedDevice)
	suite.Require().NoError(err)
	resp.Body.Close()

	suite.Equal("C5EABC521E8304EE", allowedDevice.DevEUI)

	// 3. Create device
	deviceReq := models.CreateDeviceRequest{
		VersionID:   version.ID,
		Name:        "My Integration Test Device",
		DevEUI:      "C5EABC521E8304EE",
		Description: stringPtr("Device created during integration test"),
	}

	resp, err = suite.makeAuthenticatedRequest("POST", "/devices", deviceReq)
	suite.Require().NoError(err)
	suite.Equal(http.StatusCreated, resp.StatusCode)

	var device models.Device
	err = json.NewDecoder(resp.Body).Decode(&device)
	suite.Require().NoError(err)
	resp.Body.Close()

	suite.Equal("My Integration Test Device", device.Name)
	suite.Equal("C5EABC521E8304EE", device.DevEUI)
	suite.True(device.ChirpStackDeviceCreated)
	suite.True(device.ChirpStackDeviceActivated)

	// 4. Get my devices
	resp, err = suite.makeAuthenticatedRequest("GET", "/devices/my", nil)
	suite.Require().NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)

	var deviceList models.DeviceListResponse
	err = json.NewDecoder(resp.Body).Decode(&deviceList)
	suite.Require().NoError(err)
	resp.Body.Close()

	suite.GreaterOrEqual(deviceList.Total, int64(1))
	suite.Len(deviceList.Devices, 1)
	suite.Equal(device.ID, deviceList.Devices[0].ID)

	// 5. Update device
	updateReq := models.UpdateDeviceRequest{
		Name:        stringPtr("Updated Device Name"),
		Description: stringPtr("Updated description"),
	}

	resp, err = suite.makeAuthenticatedRequest("PUT", fmt.Sprintf("/devices/%s", device.ID), updateReq)
	suite.Require().NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	resp.Body.Close()

	// 6. Get updated device
	resp, err = suite.makeAuthenticatedRequest("GET", fmt.Sprintf("/devices/%s", device.ID), nil)
	suite.Require().NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)

	var updatedDevice models.Device
	err = json.NewDecoder(resp.Body).Decode(&updatedDevice)
	suite.Require().NoError(err)
	resp.Body.Close()

	suite.Equal("Updated Device Name", updatedDevice.Name)
	suite.Equal("Updated description", updatedDevice.Description)

	// 7. Delete device
	resp, err = suite.makeAuthenticatedRequest("DELETE", fmt.Sprintf("/devices/%s", device.ID), nil)
	suite.Require().NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	resp.Body.Close()

	// 8. Verify device is deleted
	resp, err = suite.makeAuthenticatedRequest("GET", fmt.Sprintf("/devices/%s", device.ID), nil)
	suite.Require().NoError(err)
	suite.Equal(http.StatusNotFound, resp.StatusCode)
	resp.Body.Close()
}

func (suite *IntegrationTestSuite) TestDeviceVersionManagement() {
	// Create multiple versions
	versions := []models.CreateDeviceVersionRequest{
		{Name: "RAK7200", Version: "v1.0", Description: stringPtr("Version 1.0")},
		{Name: "RAK7200", Version: "v1.1", Description: stringPtr("Version 1.1")},
		{Name: "RAK7200", Version: "v2.0", Description: stringPtr("Version 2.0")},
	}

	createdVersions := make([]models.DeviceVersion, 0, len(versions))

	for _, versionReq := range versions {
		resp, err := suite.makeAuthenticatedRequest("POST", "/devices/versions", versionReq)
		suite.Require().NoError(err)
		suite.Equal(http.StatusCreated, resp.StatusCode)

		var version models.DeviceVersion
		err = json.NewDecoder(resp.Body).Decode(&version)
		suite.Require().NoError(err)
		resp.Body.Close()

		createdVersions = append(createdVersions, version)
	}

	// Get all versions
	resp, err := suite.makeAuthenticatedRequest("GET", "/devices/versions", nil)
	suite.Require().NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)

	var versionList models.DeviceVersionListResponse
	err = json.NewDecoder(resp.Body).Decode(&versionList)
	suite.Require().NoError(err)
	resp.Body.Close()

	suite.GreaterOrEqual(versionList.Total, int64(3))

	// Update a version
	updateReq := models.UpdateDeviceVersionRequest{
		Description: stringPtr("Updated description"),
	}

	resp, err = suite.makeAuthenticatedRequest("PUT", fmt.Sprintf("/devices/versions/%s", createdVersions[0].ID), updateReq)
	suite.Require().NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	resp.Body.Close()

	// Delete a version
	resp, err = suite.makeAuthenticatedRequest("DELETE", fmt.Sprintf("/devices/versions/%s", createdVersions[2].ID), nil)
	suite.Require().NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	resp.Body.Close()
}

func (suite *IntegrationTestSuite) TestUserProfile() {
	// Get user profile
	resp, err := suite.makeAuthenticatedRequest("GET", "/user/profile", nil)
	suite.Require().NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)

	var user models.User
	err = json.NewDecoder(resp.Body).Decode(&user)
	suite.Require().NoError(err)
	resp.Body.Close()

	suite.Equal(suite.userID, user.ID.String())
	suite.NotNil(user.TenantID)
	suite.NotNil(user.ApplicationID)
	suite.NotNil(user.DeviceProfileID)
}

func stringPtr(s string) *string {
	return &s
}

func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
