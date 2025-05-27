CREATE TABLE otp_verification (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    phone_number VARCHAR(20) UNIQUE NOT NULL,
    otp TEXT[] NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create an index on phone_number for faster lookups
CREATE INDEX idx_otp_verification_phone_number ON otp_verification(phone_number);

-- Create a function to automatically update the updated_at timestamp
CREATE OR REPLACE FUNCTION update_modified_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create a trigger to call the function before each update
CREATE TRIGGER update_otp_verification_timestamp
BEFORE UPDATE ON otp_verification
FOR EACH ROW
EXECUTE FUNCTION update_modified_column();

CREATE TABLE flyer_distribution_details (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID REFERENCES flyer_orders(id) ON DELETE CASCADE,

    target_area TEXT, -- GeoJSON, KML, or custom format depending on your mapping tech
    distribution_date DATE NOT NULL,
    pavement_required BOOLEAN DEFAULT FALSE,
    status VARCHAR(50) DEFAULT 'scheduled', -- scheduled, in_progress, completed

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE flyer_print_details (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID REFERENCES flyer_orders(id) ON DELETE CASCADE,

    upload_type VARCHAR(20) NOT NULL CHECK (upload_type IN ('upload_pdf', 'request_design')),
    design_file TEXT, -- Base64 or URL
    flyer_size VARCHAR(50), -- e.g., A4, A5
    paper_quality VARCHAR(50), -- e.g., Glossy, Matte
    flyer_quantity INT NOT NULL,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE flyer_orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    service_type VARCHAR(20) NOT NULL CHECK (service_type IN ('print', 'distribute', 'both')),

    order_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(50) DEFAULT 'pending', -- e.g., pending, in_progress, completed, cancelled
    total_cost DECIMAL(10,2), -- optional

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Step 1: Create ENUM types
CREATE TYPE account_type_enum AS ENUM ('delivery_partner', 'customer', 'admin');
CREATE TYPE gender_enum AS ENUM ('male', 'female', 'other');

-- Step 2: Create the users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),  -- requires pgcrypto extension
    account_type account_type_enum NOT NULL,

    phone_number VARCHAR(20) UNIQUE NOT NULL,
    email VARCHAR(255),
    name VARCHAR(255) NOT NULL,
    photo TEXT,

    -- Delivery Partner Fields
    gender gender_enum,
    aadhar_number VARCHAR(20),
    aadhar_photo_front TEXT,
    aadhar_photo_back TEXT,
    vehicle_type VARCHAR(100),
    age INT,

    -- Customer Fields
    gst_number VARCHAR(30),

    -- Admin Fields
    admin_role VARCHAR(50),

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE auth_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    auth_token TEXT NOT NULL,
    refresh_token TEXT NOT NULL,
    user_agent TEXT,
    ip_address TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    auth_token_expires_at TIMESTAMP,
    refresh_token_expires_at TIMESTAMP,
    auth_token_hash TEXT,
    revoked BOOLEAN DEFAULT FALSE
);

-- Populate the hash column (example using PostgreSQL's built-in functions)
UPDATE auth_tokens SET auth_token_hash = encode(digest(auth_token, 'sha256'), 'hex');

-- Create an index on the hash column
CREATE INDEX idx_auth_tokens_auth_token_hash ON auth_tokens (auth_token_hash);
