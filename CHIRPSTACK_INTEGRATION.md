# ChirpStack Integration

This document describes the ChirpStack integration feature that automatically creates ChirpStack resources when users register.

## Overview

When a new user registers in the system, the application automatically creates the following ChirpStack resources:

1. **Tenant** - A dedicated tenant for the user
2. **Application** - An "Lnode" application within the tenant
3. **Device Profile** - A "RAK_ABP" device profile with custom payload codec

## Configuration

### Environment Variables

Add the following environment variables to enable ChirpStack integration:

```bash
# ChirpStack Configuration
CHIRPSTACK_ENABLED=true
CHIRPSTACK_HOST=192.168.0.21
CHIRPSTACK_PORT=8090
CHIRPSTACK_TOKEN=your-chirpstack-api-token
```

### Docker Compose

The environment variables are already configured in `docker-compose.yml`:

```yaml
environment:
  - CHIRPSTACK_ENABLED=true
  - CHIRPSTACK_HOST=192.168.0.21
  - CHIRPSTACK_PORT=8090
  - CHIRPSTACK_TOKEN=eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9...
```

## Features

### Automatic Resource Creation

When a user registers via `/api/v1/auth/register`, the system automatically:

1. **Creates a Tenant** with the user's email as the name
2. **Creates an Application** named "Lnode" within the tenant
3. **Creates a Device Profile** named "RAK_ABP" with:
   - LoRaWAN 1.0.3 compatibility
   - AS923_2 region configuration
   - ABP activation mode
   - Custom JavaScript payload decoder
   - Predefined measurements for IoT sensors

### Device Profile Specifications

The automatically created device profile includes:

- **Name**: RAK_ABP
- **Region**: AS923_2
- **MAC Version**: LoRaWAN 1.0.3
- **Activation**: ABP (Activation By Personalization)
- **Class**: Class C support
- **Payload Codec**: JavaScript with custom decoder

### Payload Decoder

The device profile includes a comprehensive JavaScript payload decoder that handles:

- **Header Device Types**:
  - Type 1: Sensor data (dimming, voltage, current, power, energy, PF, tilt, lamp status)
  - Type 2: GPS coordinates (latitude, longitude, altitude)
  - Type 3: Timestamp data
  - Type 4: Status codes

- **Measurements**:
  - Dimming, Energy, PF, Power, Status_lamp, Tilt
  - Altitude, Current, Voltage
  - GPS coordinates (lat, lng)
  - Header device type, Status code, Timestamp

## API Integration

### ChirpStack Service

The `ChirpStackService` provides the following methods:

```go
// Check if ChirpStack integration is enabled
func (cs *ChirpStackService) IsEnabled() bool

// Create a tenant for a user
func (cs *ChirpStackService) CreateTenant(name string) (string, error)

// Create an application within a tenant
func (cs *ChirpStackService) CreateApplication(tenantID, name string) (string, error)

// Create a device profile within a tenant
func (cs *ChirpStackService) CreateDeviceProfile(tenantID string) (string, error)

// Create all resources for a user (tenant, application, device profile)
func (cs *ChirpStackService) CreateUserResources(username string) (*models.ChirpStackUserData, error)
```

### User Registration Flow

1. User submits registration request
2. User account is created in database
3. JWT token is generated
4. **ChirpStack resources are created** (if enabled)
5. Success response is returned

The ChirpStack integration is non-blocking - if it fails, the user registration still succeeds, but a warning is logged.

## Testing

### Manual Testing

Use the provided test script to verify the integration:

```bash
./test_chirpstack_integration.sh
```

This script will:
1. Test ChirpStack API connectivity
2. Register a new user
3. Verify that tenant, application, and device profile were created

### Expected Output

```
=== Testing ChirpStack Integration ===

1. Testing ChirpStack API connectivity...
✓ ChirpStack API is accessible

2. Testing User Registration with ChirpStack Integration...
✓ User registration successful

3. Checking ChirpStack Resources Creation...

3.1. Checking Tenants...
✓ Successfully retrieved tenants
✓ Tenant for user found
Tenant ID: 31230386-cd2a-4b6c-955c-afabc2d21079

3.2. Checking Applications...
✓ Successfully retrieved applications
✓ Lnode application found
Application ID: 2a998c23-739f-4df2-9b0d-af420fb6c745

3.3. Checking Device Profiles...
✓ Successfully retrieved device profiles
✓ RAK_ABP device profile found
Device Profile ID: cc925331-76cd-4446-b5e4-1f3ae4f0eca0
```

## Troubleshooting

### Common Issues

1. **ChirpStack API not accessible**
   - Check CHIRPSTACK_HOST and CHIRPSTACK_PORT
   - Verify network connectivity
   - Ensure ChirpStack server is running

2. **Authentication failed**
   - Verify CHIRPSTACK_TOKEN is valid
   - Check token permissions in ChirpStack

3. **Resources not created**
   - Check application logs for error messages
   - Verify ChirpStack API permissions
   - Ensure tenant limits are not exceeded

### Logs

The application logs ChirpStack operations:

```
Successfully created ChirpStack resources for user test@example.com: 
TenantID=31230386-cd2a-4b6c-955c-afabc2d21079, 
ApplicationID=2a998c23-739f-4df2-9b0d-af420fb6c745, 
DeviceProfileID=cc925331-76cd-4446-b5e4-1f3ae4f0eca0
```

### Disabling Integration

To disable ChirpStack integration, set:

```bash
CHIRPSTACK_ENABLED=false
```

Or remove the environment variable entirely.
