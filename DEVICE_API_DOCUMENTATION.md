# Device Management API Documentation

## Overview

This API provides comprehensive device management functionality including:
- Device version management
- Allowed device management (device keys)
- User device management with ChirpStack integration

## Base URL
```
http://localhost:8080/api/v1
```

## Authentication
All endpoints require JWT authentication via Bearer token:
```
Authorization: Bearer <your-jwt-token>
```

---

## Device Versions

### Create Device Version
**POST** `/devices/versions`

Create a new device version.

**Request Body:**
```json
{
  "name": "RAK7200",
  "version": "v1.2",
  "description": "RAK7200 LoRaWAN Tracker v1.2"
}
```

**Response:**
```json
{
  "id": "uuid",
  "name": "RAK7200",
  "version": "v1.2",
  "description": "RAK7200 LoRaWAN Tracker v1.2",
  "created_at": "2025-06-10T16:32:18Z",
  "updated_at": "2025-06-10T16:32:18Z"
}
```

### Get Device Versions
**GET** `/devices/versions?page=1&page_size=10`

Get paginated list of device versions.

**Response:**
```json
{
  "versions": [
    {
      "id": "uuid",
      "name": "RAK7200",
      "version": "v1.0",
      "description": "RAK7200 LoRaWAN Tracker v1.0",
      "created_at": "2025-06-10T16:32:18Z",
      "updated_at": "2025-06-10T16:32:18Z"
    }
  ],
  "total": 4,
  "page": 1,
  "page_size": 10,
  "total_pages": 1
}
```

### Get Device Version by ID
**GET** `/devices/versions/{id}`

### Update Device Version
**PUT** `/devices/versions/{id}`

### Delete Device Version
**DELETE** `/devices/versions/{id}`

---

## Allowed Devices

### Create Allowed Device
**POST** `/devices/allowed`

Add a device to the allowed devices list with its keys.

**Request Body:**
```json
{
  "dev_eui": "C5EABC521E8304EE",
  "nwk_key": "C518B15AB390B01762E4A3730E8C5F1C",
  "app_key": "97784F3B7F2A57EECF19F10E625081E0",
  "addr_key": "2F972E56",
  "description": "Test device 1"
}
```

**Response:**
```json
{
  "id": "uuid",
  "dev_eui": "C5EABC521E8304EE",
  "nwk_key": "C518B15AB390B01762E4A3730E8C5F1C",
  "app_key": "97784F3B7F2A57EECF19F10E625081E0",
  "addr_key": "2F972E56",
  "description": "Test device 1",
  "created_at": "2025-06-10T16:32:18Z",
  "updated_at": "2025-06-10T16:32:18Z"
}
```

### Get Allowed Devices
**GET** `/devices/allowed?page=1&page_size=10`

### Get Allowed Device by DevEUI
**GET** `/devices/allowed/{devEUI}`

### Update Allowed Device
**PUT** `/devices/allowed/{devEUI}`

### Delete Allowed Device
**DELETE** `/devices/allowed/{devEUI}`

---

## User Devices

### Create Device
**POST** `/devices`

Create a device for the authenticated user. Automatically creates and activates device in ChirpStack.

**Request Body:**
```json
{
  "version_id": "9c521c6f-6e94-4668-ac90-d5f077f79c6f",
  "name": "My RAK7200 Device",
  "dev_eui": "C5EABC521E8304EE",
  "description": "My first IoT device"
}
```

**Response:**
```json
{
  "id": "uuid",
  "user_id": "uuid",
  "version_id": "uuid",
  "name": "My RAK7200 Device",
  "dev_eui": "C5EABC521E8304EE",
  "description": "My first IoT device",
  "chirpstack_device_created": true,
  "chirpstack_device_activated": true,
  "is_active": true,
  "created_at": "2025-06-10T16:38:33Z",
  "updated_at": "2025-06-10T16:38:33Z",
  "version": {
    "id": "uuid",
    "name": "RAK7200",
    "version": "v1.0",
    "description": "RAK7200 LoRaWAN Tracker v1.0",
    "created_at": "2025-06-10T16:32:18Z",
    "updated_at": "2025-06-10T16:32:18Z"
  }
}
```

### Get My Devices
**GET** `/devices/my?page=1&page_size=10`

Get devices for the authenticated user.

### Get All Devices (Admin)
**GET** `/devices/all?page=1&page_size=10`

Get all devices in the system.

### Get Device by ID
**GET** `/devices/{id}`

### Update Device
**PUT** `/devices/{id}`

**Request Body:**
```json
{
  "name": "Updated Device Name",
  "description": "Updated description",
  "version_id": "new-version-uuid"
}
```

### Delete Device
**DELETE** `/devices/{id}`

Deletes a device from the database and automatically removes it from ChirpStack if it was created there.

**Response:**
```json
{
  "message": "Device deleted successfully"
}
```

**Note:** This operation will:
1. Retrieve device information from database
2. Delete device from ChirpStack (if ChirpStack integration is enabled and device exists there)
3. Delete device from database
4. Continue with database deletion even if ChirpStack deletion fails (with warning log)

---

## Error Responses

All endpoints return appropriate HTTP status codes and error messages:

**400 Bad Request:**
```json
{
  "error": "Invalid request body"
}
```

**401 Unauthorized:**
```json
{
  "error": "User not authenticated"
}
```

**404 Not Found:**
```json
{
  "error": "Device not found"
}
```

**500 Internal Server Error:**
```json
{
  "error": "Internal server error"
}
```

---

## ChirpStack Integration

When creating a device:
1. System checks if user has ChirpStack tenant, application, and device profile
2. Validates that the DevEUI exists in allowed devices
3. Creates device in ChirpStack with the provided name and DevEUI
4. Activates device in ChirpStack using keys from allowed devices table
5. Updates device status in database

**ChirpStack Requirements:**
- User must be registered (automatic ChirpStack resources creation)
- DevEUI must exist in allowed_devices table
- ChirpStack service must be enabled and configured

---

## Sample Workflow

1. **Register User** → Automatic ChirpStack tenant/application/profile creation
2. **Get Device Versions** → Choose appropriate device version
3. **Get Allowed Devices** → Choose available DevEUI
4. **Create Device** → Automatic ChirpStack device creation and activation
5. **Get My Devices** → View created devices with status
