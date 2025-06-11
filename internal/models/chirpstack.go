package models

import "time"

// ChirpStack Tenant models
type ChirpStackTenant struct {
	CanHaveGateways      bool              `json:"canHaveGateways"`
	Description          string            `json:"description"`
	MaxDeviceCount       int               `json:"maxDeviceCount"`
	MaxGatewayCount      int               `json:"maxGatewayCount"`
	Name                 string            `json:"name"`
	PrivateGatewaysDown  bool              `json:"privateGatewaysDown"`
	PrivateGatewaysUp    bool              `json:"privateGatewaysUp"`
	Tags                 map[string]string `json:"tags"`
}

type CreateTenantRequest struct {
	Tenant ChirpStackTenant `json:"tenant"`
}

type CreateTenantResponse struct {
	ID string `json:"id"`
}

// ChirpStack Application models
type ChirpStackApplication struct {
	TenantID    string            `json:"tenantId"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Tags        map[string]string `json:"tags"`
}

type CreateApplicationRequest struct {
	Application ChirpStackApplication `json:"application"`
}

type CreateApplicationResponse struct {
	ID string `json:"id"`
}

// ChirpStack Device Profile models
type DeviceProfileMeasurement struct {
	Name string `json:"name"`
	Kind string `json:"kind"`
}

type ChirpStackDeviceProfile struct {
	TenantID                           string                              `json:"tenantId"`
	Name                               string                              `json:"name"`
	Description                        string                              `json:"description"`
	Region                             string                              `json:"region"`
	MacVersion                         string                              `json:"macVersion"`
	RegParamsRevision                  string                              `json:"regParamsRevision"`
	AdrAlgorithmID                     string                              `json:"adrAlgorithmId"`
	PayloadCodecRuntime                string                              `json:"payloadCodecRuntime"`
	PayloadCodecScript                 string                              `json:"payloadCodecScript"`
	FlushQueueOnActivate               bool                                `json:"flushQueueOnActivate"`
	UplinkInterval                     int                                 `json:"uplinkInterval"`
	DeviceStatusReqInterval            int                                 `json:"deviceStatusReqInterval"`
	SupportsOtaa                       bool                                `json:"supportsOtaa"`
	SupportsClassB                     bool                                `json:"supportsClassB"`
	SupportsClassC                     bool                                `json:"supportsClassC"`
	ClassBTimeout                      int                                 `json:"classBTimeout"`
	ClassBPingSlotNbK                  int                                 `json:"classBPingSlotNbK"`
	ClassBPingSlotDr                   int                                 `json:"classBPingSlotDr"`
	ClassBPingSlotFreq                 int                                 `json:"classBPingSlotFreq"`
	ClassCTimeout                      int                                 `json:"classCTimeout"`
	AbpRx1Delay                        int                                 `json:"abpRx1Delay"`
	AbpRx1DrOffset                     int                                 `json:"abpRx1DrOffset"`
	AbpRx2Dr                           int                                 `json:"abpRx2Dr"`
	AbpRx2Freq                         int                                 `json:"abpRx2Freq"`
	Tags                               map[string]string                   `json:"tags"`
	Measurements                       map[string]DeviceProfileMeasurement `json:"measurements"`
	AutoDetectMeasurements             bool                                `json:"autoDetectMeasurements"`
	RegionConfigID                     string                              `json:"regionConfigId"`
	IsRelay                            bool                                `json:"isRelay"`
	IsRelayEd                          bool                                `json:"isRelayEd"`
	RelayEdRelayOnly                   bool                                `json:"relayEdRelayOnly"`
	RelayEnabled                       bool                                `json:"relayEnabled"`
	RelayCadPeriodicity                string                              `json:"relayCadPeriodicity"`
	RelayDefaultChannelIndex           int                                 `json:"relayDefaultChannelIndex"`
	RelaySecondChannelFreq             int                                 `json:"relaySecondChannelFreq"`
	RelaySecondChannelDr               int                                 `json:"relaySecondChannelDr"`
	RelaySecondChannelAckOffset        string                              `json:"relaySecondChannelAckOffset"`
	RelayEdActivationMode              string                              `json:"relayEdActivationMode"`
	RelayEdSmartEnableLevel            int                                 `json:"relayEdSmartEnableLevel"`
	RelayEdBackOff                     int                                 `json:"relayEdBackOff"`
	RelayEdUplinkLimitBucketSize       int                                 `json:"relayEdUplinkLimitBucketSize"`
	RelayEdUplinkLimitReloadRate       int                                 `json:"relayEdUplinkLimitReloadRate"`
	RelayJoinReqLimitReloadRate        int                                 `json:"relayJoinReqLimitReloadRate"`
	RelayNotifyLimitReloadRate         int                                 `json:"relayNotifyLimitReloadRate"`
	RelayGlobalUplinkLimitReloadRate   int                                 `json:"relayGlobalUplinkLimitReloadRate"`
	RelayOverallLimitReloadRate        int                                 `json:"relayOverallLimitReloadRate"`
	RelayJoinReqLimitBucketSize        int                                 `json:"relayJoinReqLimitBucketSize"`
	RelayNotifyLimitBucketSize         int                                 `json:"relayNotifyLimitBucketSize"`
	RelayGlobalUplinkLimitBucketSize   int                                 `json:"relayGlobalUplinkLimitBucketSize"`
	RelayOverallLimitBucketSize        int                                 `json:"relayOverallLimitBucketSize"`
	AllowRoaming                       bool                                `json:"allowRoaming"`
	Rx1Delay                           int                                 `json:"rx1Delay"`
}

type CreateDeviceProfileRequest struct {
	DeviceProfile ChirpStackDeviceProfile `json:"deviceProfile"`
}

type CreateDeviceProfileResponse struct {
	ID string `json:"id"`
}

// ChirpStack User Integration Response
type ChirpStackUserData struct {
	TenantID          string    `json:"tenant_id"`
	ApplicationID     string    `json:"application_id"`
	DeviceProfileID   string    `json:"device_profile_id"`
	CreatedAt         time.Time `json:"created_at"`
}
