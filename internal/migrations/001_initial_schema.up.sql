-- Create enum type to check and constrain device status
CREATE TYPE device_state AS ENUM ('available', 'in-use', 'inactive');


CREATE TABLE devices (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    brand TEXT NOT NULL,
    state device_state NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_devices_brand ON devices (brand);

CREATE INDEX idx_devices_state ON devices (state);
