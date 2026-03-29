-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Trigger function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- 1. Communities (The Tenants)
CREATE TABLE IF NOT EXISTS communities (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    address TEXT,
    timezone VARCHAR(50) DEFAULT 'UTC',
    settings JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TRIGGER update_communities_modtime BEFORE UPDATE ON communities 
FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();

-- 2. Users and Roles
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    phone VARCHAR(20),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TRIGGER update_users_modtime BEFORE UPDATE ON users 
FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();

-- Junction table for Multi-tenancy (Users can belong to multiple communities)
CREATE TABLE IF NOT EXISTS community_users (
    community_id UUID REFERENCES communities(id),
    user_id UUID REFERENCES users(id),
    role VARCHAR(50) NOT NULL, -- e.g., 'ADMIN', 'SECURITY', 'RESIDENT', 'MANAGER'
    PRIMARY KEY (community_id, user_id)
);
CREATE INDEX IF NOT EXISTS idx_community_users_user_id ON community_users(user_id);

-- 3. Units
CREATE TABLE IF NOT EXISTS units (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    community_id UUID NOT NULL REFERENCES communities(id),
    unit_number VARCHAR(50) NOT NULL,
    block_floor VARCHAR(50),
    unit_type VARCHAR(50), -- e.g., 'APARTMENT', 'VILLA'
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    UNIQUE(community_id, unit_number)
);

CREATE TRIGGER update_units_modtime BEFORE UPDATE ON units 
FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();

-- 4. Households and Residents
CREATE TABLE IF NOT EXISTS households (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    community_id UUID NOT NULL REFERENCES communities(id),
    unit_id UUID NOT NULL REFERENCES units(id),
    name VARCHAR(255), -- e.g., "The Smith Household"
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER update_households_modtime BEFORE UPDATE ON households 
FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();

CREATE TABLE IF NOT EXISTS residents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    community_id UUID NOT NULL REFERENCES communities(id),
    unit_id UUID NOT NULL REFERENCES units(id),
    user_id UUID REFERENCES users(id),
    household_id UUID REFERENCES households(id),
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(255),
    phone VARCHAR(20),
    resident_type VARCHAR(50) NOT NULL, -- e.g., 'OWNER', 'TENANT', 'DEPENDENT'
    status VARCHAR(50) DEFAULT 'ACTIVE', -- 'ACTIVE', 'INACTIVE', 'PENDING'
    is_primary_contact BOOLEAN DEFAULT false,
    move_in_date DATE,
    move_out_date DATE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TRIGGER update_residents_modtime BEFORE UPDATE ON residents 
FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();

-- 5. Vendors
CREATE TABLE IF NOT EXISTS vendors (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    community_id UUID NOT NULL REFERENCES communities(id),
    company_name VARCHAR(255) NOT NULL,
    category VARCHAR(100), -- 'PLUMBING', 'ELECTRICAL'
    contact_name VARCHAR(100),
    email VARCHAR(255),
    phone VARCHAR(20),
    address TEXT,
    status VARCHAR(50) DEFAULT 'ACTIVE', -- 'ACTIVE', 'INACTIVE', 'PROBATION'
    rating DECIMAL(3,2) DEFAULT 0.00,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER update_vendors_modtime BEFORE UPDATE ON vendors 
FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();

-- 6. Visitors and Security Logs
CREATE TABLE IF NOT EXISTS visitors (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    community_id UUID NOT NULL REFERENCES communities(id),
    unit_id UUID NOT NULL REFERENCES units(id),
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    vehicle_number VARCHAR(50),
    expected_arrival TIMESTAMPTZ,
    invite_code VARCHAR(10) UNIQUE,
    purpose VARCHAR(255),
    status VARCHAR(50) DEFAULT 'SCHEDULED', -- 'SCHEDULED', 'ARRIVED', 'DEPARTED', 'EXPIRED'
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS security_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    community_id UUID NOT NULL REFERENCES communities(id),
    person_type VARCHAR(50) NOT NULL, -- 'RESIDENT', 'VISITOR', 'VENDOR'
    visitor_id UUID REFERENCES visitors(id),
    resident_id UUID REFERENCES residents(id),
    vendor_id UUID REFERENCES vendors(id),
    unit_id UUID REFERENCES units(id),
    time_in TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    time_out TIMESTAMPTZ,
    guard_id UUID REFERENCES users(id), -- The security guard
    remarks TEXT
);

-- 7. Maintenance
CREATE TABLE IF NOT EXISTS maintenance_requests (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    community_id UUID NOT NULL REFERENCES communities(id),
    unit_id UUID NOT NULL REFERENCES units(id),
    resident_id UUID NOT NULL REFERENCES residents(id),
    vendor_id UUID REFERENCES vendors(id),
    assigned_to_user_id UUID REFERENCES users(id),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(100), -- 'PLUMBING', 'ELECTRICAL', 'GENERAL', etc.
    status VARCHAR(50) DEFAULT 'OPEN', -- 'OPEN', 'IN_PROGRESS', 'COMPLETED', 'CANCELLED'
    priority VARCHAR(50) DEFAULT 'MEDIUM',
    scheduled_date TIMESTAMPTZ,
    attachments JSONB DEFAULT '[]',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER update_maintenance_requests_modtime BEFORE UPDATE ON maintenance_requests 
FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();

-- 7. Invoices and Payments
CREATE TABLE IF NOT EXISTS invoices (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    community_id UUID NOT NULL REFERENCES communities(id),
    unit_id UUID NOT NULL REFERENCES units(id),
    amount DECIMAL(12,2) NOT NULL,
    paid_amount DECIMAL(12,2) DEFAULT 0.00,
    due_date DATE NOT NULL,
    status VARCHAR(50) DEFAULT 'UNPAID', -- 'UNPAID', 'PARTIAL', 'PAID', 'OVERDUE'
    billing_period VARCHAR(20) NOT NULL, -- e.g., '2024-05'
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS payments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    community_id UUID NOT NULL REFERENCES communities(id),
    invoice_id UUID NOT NULL REFERENCES invoices(id),
    amount DECIMAL(12,2) NOT NULL,
    payment_method VARCHAR(50), -- 'STRIPE', 'CASH', 'BANK_TRANSFER'
    transaction_id VARCHAR(255),
    payment_date TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_payments_invoice_id ON payments(invoice_id);

-- 8. Announcements and Notifications
CREATE TABLE IF NOT EXISTS announcements (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    community_id UUID NOT NULL REFERENCES communities(id),
    author_id UUID REFERENCES users(id),
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    audience_type VARCHAR(50) DEFAULT 'ALL', -- 'ALL', 'OWNERS', 'TENANTS'
    priority BOOLEAN DEFAULT false,
    expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id),
    community_id UUID NOT NULL REFERENCES communities(id),
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    is_read BOOLEAN DEFAULT false,
    link_to_resource VARCHAR(255), -- e.g., '/invoices/uuid'
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 9. Recommended Indexes
CREATE INDEX IF NOT EXISTS idx_units_community_id ON units(community_id);
CREATE INDEX IF NOT EXISTS idx_residents_household_id ON residents(household_id);
CREATE INDEX IF NOT EXISTS idx_security_logs_community_id ON security_logs(community_id);
CREATE INDEX IF NOT EXISTS idx_maintenance_status ON maintenance_requests(status);
CREATE INDEX IF NOT EXISTS idx_invoices_unit_id ON invoices(unit_id);
CREATE INDEX IF NOT EXISTS idx_invoices_status ON invoices(status);
CREATE INDEX IF NOT EXISTS idx_notifications_user_unread ON notifications(user_id) WHERE is_read = false;
CREATE INDEX IF NOT EXISTS idx_maintenance_requests_resident_id ON maintenance_requests(resident_id);
CREATE INDEX IF NOT EXISTS idx_visitors_invite_code ON visitors(invite_code);

-- Soft delete partial index
CREATE INDEX IF NOT EXISTS idx_users_active ON users(id) WHERE deleted_at IS NULL;

-- Enforce one primary contact per household
CREATE UNIQUE INDEX IF NOT EXISTS idx_primary_contact_per_household ON residents(household_id) WHERE (is_primary_contact = true AND deleted_at IS NULL);

-- Example of Row-Level Security (RLS) setup
-- 1. Enable RLS on a table
ALTER TABLE announcements ENABLE ROW LEVEL SECURITY;

-- 2. Create a policy that restricts access based on a session variable
-- This variable 'app.current_community_id' will be set by the Go middleware per request
CREATE POLICY community_isolation_policy ON announcements
    USING (community_id = current_setting('app.current_community_id')::uuid);

-- Repeat for all community-scoped tables