-- +goose Up
-- --------------------
-- Create tables and seed data
-- --------------------
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS expenses (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID NOT NULL DEFAULT gen_random_uuid(),
    user_id BIGINT NOT NULL REFERENCES users(id),
    amount_idr BIGINT NOT NULL,
    description TEXT NOT NULL,
    receipt_url TEXT,
    status VARCHAR(50) NOT NULL,
    submitted_at TIMESTAMP NOT NULL,
    processed_at TIMESTAMP,
    requires_approval BOOLEAN NOT NULL,
    auto_approved BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS approvals (
    id BIGSERIAL PRIMARY KEY,
    expense_id BIGINT NOT NULL REFERENCES expenses(id),
    approver_id BIGINT NULL REFERENCES users(id),
    status VARCHAR(50) NOT NULL,
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_expenses_user_id ON expenses(user_id);
CREATE INDEX IF NOT EXISTS idx_expenses_status ON expenses(status);
CREATE INDEX IF NOT EXISTS idx_approvals_expense_id ON approvals(expense_id);

-- -----------------------
-- Seed additional users
-- -----------------------
INSERT INTO users (email, name, role, password_hash)
VALUES
('alice@manager.com', 'Alice Manager', 'manager', '$2a$10$FZxuzOEihHwkJr2TVv80zuRjeXnnMxaPx7da6He90FLMx2V72/LBm'),
('bob@user.com', 'Bob Employee', 'user', '$2a$10$FZxuzOEihHwkJr2TVv80zuRjeXnnMxaPx7da6He90FLMx2V72/LBm'),
('dave@user.com', 'Dave Employee', 'user', '$2a$10$FZxuzOEihHwkJr2TVv80zuRjeXnnMxaPx7da6He90FLMx2V72/LBm'),
('eve@user.com', 'Eve Employee', 'user', '$2a$10$FZxuzOEihHwkJr2TVv80zuRjeXnnMxaPx7da6He90FLMx2V72/LBm');

INSERT INTO expenses (user_id, amount_idr, description, receipt_url, status, submitted_at, requires_approval, auto_approved)
VALUES
-- Bob (4)
((SELECT id FROM users WHERE email='bob@user.com'), 500000, 'Lunch with client', '/receipt-placeholder.png', 'pending',  NOW(),                       true,  false),
((SELECT id FROM users WHERE email='bob@user.com'), 80000,  'Taxi fare',        '/receipt-placeholder.png', 'approved', NOW() - INTERVAL '2 days', false, true),
((SELECT id FROM users WHERE email='bob@user.com'), 1200000,'Hotel stay',       '/receipt-placeholder.png', 'pending',  NOW() - INTERVAL '1 day',  true,  false),
((SELECT id FROM users WHERE email='bob@user.com'), 45000,  'Coffee meeting',   '/receipt-placeholder.png', 'approved', NOW() - INTERVAL '5 days', false, true),

-- Dave (4)
((SELECT id FROM users WHERE email='dave@user.com'), 120000, 'Office supplies',  '/receipt-placeholder.png', 'pending',  NOW() - INTERVAL '1 day',  true,  false),
((SELECT id FROM users WHERE email='dave@user.com'), 450000, 'Team lunch',       '/receipt-placeholder.png', 'approved', NOW() - INTERVAL '3 days', false, true),
((SELECT id FROM users WHERE email='dave@user.com'), 2000000,'Client workshop',  '/receipt-placeholder.png', 'pending',  NOW(),                       true,  false),
((SELECT id FROM users WHERE email='dave@user.com'), 95000,  'Parking fee',      '/receipt-placeholder.png', 'approved', NOW() - INTERVAL '6 days', false, true),

-- Eve (4)
((SELECT id FROM users WHERE email='eve@user.com'), 300000, 'Project materials', '/receipt-placeholder.png', 'pending',  NOW(),                       true,  false),
((SELECT id FROM users WHERE email='eve@user.com'), 100000, 'Taxi fare',         '/receipt-placeholder.png', 'approved', NOW() - INTERVAL '4 days', false, true),
((SELECT id FROM users WHERE email='eve@user.com'), 1750000,'Design software',   '/receipt-placeholder.png', 'pending',  NOW() - INTERVAL '2 days', true,  false),
((SELECT id FROM users WHERE email='eve@user.com'), 60000,  'Snacks',            '/receipt-placeholder.png', 'approved', NOW() - INTERVAL '7 days', false, true);

-- Approvals for pending expenses
INSERT INTO approvals (expense_id, status, created_at)
SELECT
    e.id AS expense_id,
    'pending' AS status,
    NOW() AS created_at
FROM expenses e
JOIN users u ON e.user_id = u.id
WHERE e.status = 'pending'
  AND e.requires_approval = true;

-- +goose Down
-- --------------------
-- Drop tables (rollback)
-- --------------------
DROP TABLE IF EXISTS approvals;
DROP TABLE IF EXISTS expenses;
DROP TABLE IF EXISTS users;
