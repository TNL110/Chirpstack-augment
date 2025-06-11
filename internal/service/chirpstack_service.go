package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"go-auth-api/internal/config"
	"go-auth-api/internal/models"
	"go-auth-api/internal/repository"
)

type ChirpStackService struct {
	config     *config.Config
	httpClient *http.Client
	userRepo   *repository.UserRepository
}

func NewChirpStackService(cfg *config.Config, userRepo *repository.UserRepository) *ChirpStackService {
	return &ChirpStackService{
		config:   cfg,
		userRepo: userRepo,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (cs *ChirpStackService) IsEnabled() bool {
	return cs.config.ChirpStackEnabled && cs.config.ChirpStackToken != ""
}

func (cs *ChirpStackService) getBaseURL() string {
	return fmt.Sprintf("http://%s:%s/api", cs.config.ChirpStackHost, cs.config.ChirpStackPort)
}

func (cs *ChirpStackService) makeRequest(method, endpoint string, payload interface{}) ([]byte, error) {
	url := cs.getBaseURL() + endpoint

	var body io.Reader
	if payload != nil {
		jsonData, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal payload: %w", err)
		}
		body = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+cs.config.ChirpStackToken)

	resp, err := cs.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("ChirpStack API error (status %d): %s", resp.StatusCode, string(responseBody))
	}

	return responseBody, nil
}

func (cs *ChirpStackService) CreateTenant(username string) (string, error) {
	tenantReq := models.CreateTenantRequest{
		Tenant: models.ChirpStackTenant{
			CanHaveGateways:     true,
			Description:         username,
			MaxDeviceCount:      10000,
			MaxGatewayCount:     10000,
			Name:                username,
			PrivateGatewaysDown: true,
			PrivateGatewaysUp:   true,
			Tags:                make(map[string]string),
		},
	}

	responseBody, err := cs.makeRequest("POST", "/tenants", tenantReq)
	if err != nil {
		return "", fmt.Errorf("failed to create tenant: %w", err)
	}

	var response models.CreateTenantResponse
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return "", fmt.Errorf("failed to parse tenant response: %w", err)
	}

	return response.ID, nil
}

func (cs *ChirpStackService) CreateApplication(tenantID, name string) (string, error) {
	appReq := models.CreateApplicationRequest{
		Application: models.ChirpStackApplication{
			TenantID:    tenantID,
			Name:        name,
			Description: fmt.Sprintf("Application for %s", name),
			Tags:        make(map[string]string),
		},
	}

	responseBody, err := cs.makeRequest("POST", "/applications", appReq)
	if err != nil {
		return "", fmt.Errorf("failed to create application: %w", err)
	}

	var response models.CreateApplicationResponse
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return "", fmt.Errorf("failed to parse application response: %w", err)
	}

	return response.ID, nil
}

func (cs *ChirpStackService) CreateDeviceProfile(tenantID string) (string, error) {
	measurements := map[string]models.DeviceProfileMeasurement{
		"Dimming":       {Name: "", Kind: "UNKNOWN"},
		"Energy":        {Name: "", Kind: "UNKNOWN"},
		"PF":            {Name: "", Kind: "UNKNOWN"},
		"Power":         {Name: "", Kind: "UNKNOWN"},
		"Status_lamp":   {Name: "", Kind: "UNKNOWN"},
		"Tilt":          {Name: "", Kind: "UNKNOWN"},
		"alt":           {Name: "", Kind: "UNKNOWN"},
		"current":       {Name: "", Kind: "UNKNOWN"},
		"header_device": {Name: "", Kind: "UNKNOWN"},
		"lat":           {Name: "", Kind: "UNKNOWN"},
		"lng":           {Name: "", Kind: "UNKNOWN"},
		"status_code":   {Name: "", Kind: "UNKNOWN"},
		"timestamp":     {Name: "", Kind: "UNKNOWN"},
		"voltage":       {Name: "", Kind: "UNKNOWN"},
	}

	payloadCodecScript := `function decodeUplink(input) {
  let header_device = input.bytes[0];

  if (header_device == 4) {
    let status_code = (input.bytes[2] << 8) | input.bytes[1];
    if(status_code == 50||status_code == 51||status_code == 52||status_code == 53)
    {
      let ID = (input.bytes[6] << 24) |(input.bytes[5] << 16) |(input.bytes[4] << 8) | input.bytes[3];
      return {
        data: {
          header_device: header_device,
          status_code: status_code,
          ID: ID,
        },
        warnings: [],
        errors: []
      };
    }
    else
    {
      return {
        data: {
          header_device: header_device,
          status_code: status_code,
        },
        warnings: [],
        errors: []
      };
	}
  }

  else if (header_device == 1)
  {
    let Dimming = input.bytes[1];
    let Status_lamp = input.bytes[2];
    let Energy_raw = (input.bytes[6] << 8) |(input.bytes[5] << 8) |(input.bytes[4] << 8) | input.bytes[3];
    let voltage_raw = (input.bytes[8] << 8) | input.bytes[7];
    let current_raw = (input.bytes[10] << 8) | input.bytes[9];
    let PF_raw = (input.bytes[12] << 8) | input.bytes[11];
    let Power_raw = (input.bytes[14] << 8) | input.bytes[13];
    let Tilt_raw = (input.bytes[16] << 8) | input.bytes[15];

    let Energy = Energy_raw/100;
    let voltage = voltage_raw/100;
    let current = current_raw/100;
    let PF = PF_raw/100;
    let Power = Power_raw/100;
    let Tilt = Tilt_raw / 100;

    return {
      data: {
        header_device: header_device,
        voltage: voltage,
        current: current,
        Power: Power,
        Energy: Energy,
        PF: PF,
        Tilt: Tilt,
        Status_lamp: Status_lamp,
        Dimming: Dimming,
      },
      warnings: [],
      errors: []
    }
  }

  else if (header_device == 2)
  {
    let lat_raw = (input.bytes[4] << 24) |(input.bytes[3] << 16) |(input.bytes[2] << 8) | input.bytes[1];
    let lng_raw = (input.bytes[8] << 24) |(input.bytes[7] << 16) |(input.bytes[6] << 8) | input.bytes[5];
    let alt_raw = input.bytes[9];

    let lat = lat_raw/1000000;
    let lng = lng_raw / 1000000;

    return {
      data: {
        header_device: header_device,
        lat: lat,
        lng: lng,
        alt: alt_raw,
      },
      warnings: [],
      errors: []
    }
  }
  else if (header_device == 3)
  {
    let timestamp = (input.bytes[4] << 24) |(input.bytes[3] << 16) |(input.bytes[2] << 8) | input.bytes[1];

    return {
      data: {
        header_device: header_device,
        timestamp: timestamp,
      },
      warnings: [],
      errors: []
    }
  }
  else {
    return {
      data: { header_device: header_device },
      warnings: ["Gói tin không thuộc thiết bị được mong đợi"],
      errors: []
    };
  }
}


function encodeNumber(number) {
  const binaryString = number.toString(2);
  const paddedBinaryString = '0'.repeat(8 - binaryString.length) + binaryString;
  const part = paddedBinaryString.substring(0, 8);
  return part;
}`

	deviceProfileReq := models.CreateDeviceProfileRequest{
		DeviceProfile: models.ChirpStackDeviceProfile{
			TenantID:                         tenantID,
			Name:                             "RAK_ABP",
			Description:                      "",
			Region:                           "AS923_2",
			MacVersion:                       "LORAWAN_1_0_3",
			RegParamsRevision:                "A",
			AdrAlgorithmID:                   "default",
			PayloadCodecRuntime:              "JS",
			PayloadCodecScript:               payloadCodecScript,
			FlushQueueOnActivate:             true,
			UplinkInterval:                   3600,
			DeviceStatusReqInterval:          1,
			SupportsOtaa:                     false,
			SupportsClassB:                   false,
			SupportsClassC:                   true,
			ClassBTimeout:                    0,
			ClassBPingSlotNbK:                0,
			ClassBPingSlotDr:                 0,
			ClassBPingSlotFreq:               0,
			ClassCTimeout:                    0,
			AbpRx1Delay:                      1,
			AbpRx1DrOffset:                   0,
			AbpRx2Dr:                         2,
			AbpRx2Freq:                       921400000,
			Tags:                             make(map[string]string),
			Measurements:                     measurements,
			AutoDetectMeasurements:           true,
			RegionConfigID:                   "as923_2",
			IsRelay:                          false,
			IsRelayEd:                        false,
			RelayEdRelayOnly:                 false,
			RelayEnabled:                     false,
			RelayCadPeriodicity:              "SEC_1",
			RelayDefaultChannelIndex:         0,
			RelaySecondChannelFreq:           0,
			RelaySecondChannelDr:             0,
			RelaySecondChannelAckOffset:      "KHZ_0",
			RelayEdActivationMode:            "DISABLE_RELAY_MODE",
			RelayEdSmartEnableLevel:          0,
			RelayEdBackOff:                   0,
			RelayEdUplinkLimitBucketSize:     0,
			RelayEdUplinkLimitReloadRate:     0,
			RelayJoinReqLimitReloadRate:      0,
			RelayNotifyLimitReloadRate:       0,
			RelayGlobalUplinkLimitReloadRate: 0,
			RelayOverallLimitReloadRate:      0,
			RelayJoinReqLimitBucketSize:      0,
			RelayNotifyLimitBucketSize:       0,
			RelayGlobalUplinkLimitBucketSize: 0,
			RelayOverallLimitBucketSize:      0,
			AllowRoaming:                     false,
			Rx1Delay:                         0,
		},
	}

	responseBody, err := cs.makeRequest("POST", "/device-profiles", deviceProfileReq)
	if err != nil {
		return "", fmt.Errorf("failed to create device profile: %w", err)
	}

	var response models.CreateDeviceProfileResponse
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return "", fmt.Errorf("failed to parse device profile response: %w", err)
	}

	return response.ID, nil
}

func (cs *ChirpStackService) CreateUserResources(userID, username string) (*models.ChirpStackUserData, error) {
	if !cs.IsEnabled() {
		return nil, fmt.Errorf("ChirpStack integration is disabled")
	}

	// 1. Create Tenant
	tenantID, err := cs.CreateTenant(username)
	if err != nil {
		return nil, fmt.Errorf("failed to create tenant: %w", err)
	}

	// 2. Create Application
	applicationID, err := cs.CreateApplication(tenantID, "Lnode")
	if err != nil {
		return nil, fmt.Errorf("failed to create application: %w", err)
	}

	// 3. Create Device Profile
	deviceProfileID, err := cs.CreateDeviceProfile(tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to create device profile: %w", err)
	}

	// 4. Update user with ChirpStack data
	if cs.userRepo != nil {
		err = cs.userRepo.UpdateUserChirpStackData(userID, tenantID, applicationID, deviceProfileID)
		if err != nil {
			fmt.Printf("Warning: Failed to update user ChirpStack data: %v\n", err)
		}
	}

	return &models.ChirpStackUserData{
		TenantID:        tenantID,
		ApplicationID:   applicationID,
		DeviceProfileID: deviceProfileID,
		CreatedAt:       time.Now(),
	}, nil
}

// DeleteDevice deletes a device from ChirpStack
func (cs *ChirpStackService) DeleteDevice(devEUI string) error {
	if !cs.IsEnabled() {
		return fmt.Errorf("ChirpStack integration is disabled")
	}

	deleteURL := fmt.Sprintf("/devices/%s", devEUI)
	_, err := cs.makeRequest("DELETE", deleteURL, nil)
	if err != nil {
		return fmt.Errorf("failed to delete ChirpStack device: %w", err)
	}

	fmt.Printf("ChirpStack device deleted: %s\n", devEUI)
	return nil
}
