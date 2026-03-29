-- 1. Create a test community
-- internal/database/seed.sql

-- 1. Communities
INSERT INTO communities (id, name, address) 
VALUES 
('550e8400-e29b-41d4-a716-446655440000', 'Green Valley Estate', '123 Tech Lane, Silicon Valley'),
('550e8400-e29b-41d4-a716-446655440001', 'Skyline Heights', '456 Metro Blvd, Downtown')
ON CONFLICT (id) DO NOTHING;

-- 2. Users (Password: 'password')
INSERT INTO users (id, email, password_hash, first_name, last_name)
VALUES 
('770e8400-e29b-41d4-a716-446655440000', 'admin@example.com', '$2a$10$vI8tmvshCiUCvLRNo.reBe.In6lxIn.Z.6.6.6.6.6.6.6.6.6.6.6', 'System', 'Admin'),
('770e8400-e29b-41d4-a716-446655440001', 'guard@example.com', '$2a$10$vI8tmvshCiUCvLRNo.reBe.In6lxIn.Z.6.6.6.6.6.6.6.6.6.6.6', 'Ellen', 'Ripley'),
('770e8400-e29b-41d4-a716-446655440002', 'sarah@example.com', '$2a$10$vI8tmvshCiUCvLRNo.reBe.In6lxIn.Z.6.6.6.6.6.6.6.6.6.6.6', 'Sarah', 'Connor'),
('770e8400-e29b-41d4-a716-446655440003', 'john@example.com', '$2a$10$vI8tmvshCiUCvLRNo.reBe.In6lxIn.Z.6.6.6.6.6.6.6.6.6.6.6', 'John', 'Wick')
ON CONFLICT (id) DO NOTHING;

-- 3. Link Users
INSERT INTO community_users (community_id, user_id, role)
VALUES 
('550e8400-e29b-41d4-a716-446655440000', '770e8400-e29b-41d4-a716-446655440000', 'ADMIN'),
('550e8400-e29b-41d4-a716-446655440000', '770e8400-e29b-41d4-a716-446655440001', 'SECURITY'),
('550e8400-e29b-41d4-a716-446655440000', '770e8400-e29b-41d4-a716-446655440002', 'RESIDENT')
ON CONFLICT (community_id, user_id) DO NOTHING;

-- 4. Units
INSERT INTO units (id, community_id, unit_number, block_floor, unit_type)
VALUES 
('660e8400-e29b-41d4-a716-446655440000', '550e8400-e29b-41d4-a716-446655440000', 'B-204', 'Block B, Floor 2', 'APARTMENT'),
('660e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440000', 'A-101', 'Block A, Floor 1', 'APARTMENT')
ON CONFLICT (id) DO NOTHING;

-- 5. Visitors
INSERT INTO visitors (id, community_id, unit_id, name, phone, expected_arrival, invite_code, purpose, status)
VALUES 
('a10e8400-e29b-41d4-a716-446655440000', '550e8400-e29b-41d4-a716-446655440000', '660e8400-e29b-41d4-a716-446655440001', 'Kyle Reese', '555-123', NOW() + INTERVAL '2 hours', 'GUEST-1', 'Personal', 'SCHEDULED')
ON CONFLICT (id) DO NOTHING;

-- 6. Maintenance
INSERT INTO maintenance_requests (id, community_id, unit_id, resident_id, title, description, category, status, priority)
VALUES 
('c10e8400-e29b-41d4-a716-446655440000', '550e8400-e29b-41d4-a716-446655440000', '660e8400-e29b-41d4-a716-446655440001', '990e8400-e29b-41d4-a716-446655440000', 'Burst Pipe', 'Kitchen is flooding', 'PLUMBING', 'OPEN', 'URGENT')
ON CONFLICT (id) DO NOTHING;