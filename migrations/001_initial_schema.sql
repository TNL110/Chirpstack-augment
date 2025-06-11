-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create users table with UUID primary key
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    tenant_id VARCHAR(255),
    application_id VARCHAR(255),
    device_profile_id VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create device versions table
CREATE TABLE IF NOT EXISTS device_versions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    version VARCHAR(100) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(name, version)
);

-- Create allowed devices table
CREATE TABLE IF NOT EXISTS allowed_devices (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    dev_eui VARCHAR(16) UNIQUE NOT NULL,
    nwk_key VARCHAR(32) NOT NULL,
    app_key VARCHAR(32) NOT NULL,
    addr_key VARCHAR(8) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create devices table
CREATE TABLE IF NOT EXISTS devices (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    version_id UUID NOT NULL REFERENCES device_versions(id) ON DELETE RESTRICT,
    name VARCHAR(255) NOT NULL,
    dev_eui VARCHAR(16) NOT NULL,
    description TEXT,
    chirpstack_device_created BOOLEAN DEFAULT FALSE,
    chirpstack_device_activated BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (dev_eui) REFERENCES allowed_devices(dev_eui)
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_tenant_id ON users(tenant_id);
CREATE INDEX IF NOT EXISTS idx_devices_user_id ON devices(user_id);
CREATE INDEX IF NOT EXISTS idx_devices_dev_eui ON devices(dev_eui);
CREATE INDEX IF NOT EXISTS idx_devices_version_id ON devices(version_id);
CREATE INDEX IF NOT EXISTS idx_allowed_devices_dev_eui ON allowed_devices(dev_eui);

-- Insert sample device versions
INSERT INTO device_versions (name, version, description) VALUES 
('RAK7200', 'v1.0', 'RAK7200 LoRaWAN Tracker v1.0'),
('RAK7200', 'v1.1', 'RAK7200 LoRaWAN Tracker v1.1'),
('RAK4631', 'v1.0', 'RAK4631 WisBlock Core v1.0'),
('RAK11200', 'v1.0', 'RAK11200 WiFi Module v1.0')
ON CONFLICT (name, version) DO NOTHING;

-- Insert sample allowed devices
INSERT INTO allowed_devices (dev_eui, nwk_key, app_key, addr_key, description) VALUES 
('C5EABC521E8304EE', 'C518B15AB390B01762E4A3730E8C5F1C', '97784F3B7F2A57EECF19F10E625081E0', '2F972E56', 'Test device 1'),
('A1B2C3D4E5F60708', 'A1B2C3D4E5F607081234567890ABCDEF', '1234567890ABCDEFA1B2C3D4E5F60708', '12345678', 'Test device 2'),
('1122334455667788', '1122334455667788AABBCCDDEEFF0011', 'AABBCCDDEEFF00111122334455667788', 'AABBCCDD', 'Test device 3')
ON CONFLICT (dev_eui) DO NOTHING;
