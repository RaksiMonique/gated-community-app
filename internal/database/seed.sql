-- 1. Create a test community
INSERT INTO communities (id, name, address) 
VALUES ('550e8400-e29b-41d4-a716-446655440000', 'Green Valley Estate', '123 Tech Lane')
ON CONFLICT (id) DO NOTHING;

-- 2. Create a test unit
INSERT INTO units (id, community_id, unit_number, unit_type)
VALUES ('660e8400-e29b-41d4-a716-446655440000', '550e8400-e29b-41d4-a716-446655440000', 'B-204', 'APARTMENT')
ON CONFLICT (id) DO NOTHING;

-- 3. Add an admin user (Password is 'password')
INSERT INTO users (id, email, password_hash, first_name, last_name)
VALUES ('770e8400-e29b-41d4-a716-446655440000', 'admin@example.com', '$2a$10$Z.vMq.Nb.v...placeholder...', 'Admin', 'User')
ON CONFLICT (id) DO NOTHING;

-- 4. Link user to community as ADMIN
INSERT INTO community_users (community_id, user_id, role)
VALUES ('550e8400-e29b-41d4-a716-446655440000', '770e8400-e29b-41d4-a716-446655440000', 'ADMIN')
ON CONFLICT (community_id, user_id) DO NOTHING;