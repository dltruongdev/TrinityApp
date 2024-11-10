-- Create tables --
CREATE TABLE UserTypes (
    user_type_id SERIAL PRIMARY KEY,
    type_name VARCHAR(50) NOT NULL UNIQUE
);

CREATE TABLE Plans (
    plan_id SERIAL PRIMARY KEY,
    plan_name VARCHAR(50) NOT NULL UNIQUE,
    price DECIMAL(10, 2) NOT NULL
);

CREATE TABLE Users (
    user_id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    registration_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_login TIMESTAMP,
    user_type_id INT REFERENCES UserTypes(user_type_id),
    plan_id INT DEFAULT 1 REFERENCES Plans(plan_id) --default to basic plan (id=1)
);

CREATE TABLE Campaigns (
    campaign_id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    code VARCHAR(50) NOT NULL UNIQUE,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    max_vouchers INT NOT NULL,
    redeemed_vouchers INT DEFAULT 0,
    voucher_lifetime INT NOT NULL DEFAULT 30,
    discount_percentage INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_campaigns_redeemed_vouchers ON Campaigns (redeemed_vouchers);
CREATE INDEX idx_campaigns_start_date ON Campaigns (start_date);
CREATE INDEX idx_campaigns_end_date ON Campaigns (end_date);
CREATE INDEX idx_campaigns_code ON Campaigns (code);

CREATE TABLE CampaignPlans (
    campaign_id INT REFERENCES Campaigns(campaign_id) ON DELETE CASCADE,
    plan_id INT REFERENCES Plans(plan_id) ON DELETE CASCADE,
    PRIMARY KEY (campaign_id, plan_id)
);

CREATE TABLE Vouchers (
    voucher_id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES Users(user_id) ON DELETE CASCADE,
    campaign_id INT NOT NULL REFERENCES Campaigns(campaign_id) ON DELETE CASCADE,
    code VARCHAR(50) NOT NULL UNIQUE,
    valid_until TIMESTAMP NOT NULL,
    is_redeemed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_vouchers_valid_until ON Vouchers (valid_until);
CREATE INDEX idx_vouchers_is_redeemed ON Vouchers (is_redeemed);
CREATE INDEX idx_vouchers_user_id ON Vouchers (user_id);
CREATE INDEX idx_vouchers_campaign_id ON Vouchers (campaign_id);

CREATE TABLE Purchases (
    purchase_id SERIAL PRIMARY KEY,
    user_id INT REFERENCES Users(user_id),
    voucher_id INT REFERENCES Vouchers(voucher_id),
    plan_id INT REFERENCES Plans(plan_id),
    amount DECIMAL(10, 2) NOT NULL,
    discount_applied BOOLEAN DEFAULT FALSE,
    purchase_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_purchases_user_id ON Purchases (user_id);
CREATE INDEX idx_purchases_purchase_date ON Purchases (purchase_date);

-- Data seeding --
-- User Types
INSERT INTO UserTypes (type_name) VALUES
('Admin'),
('User');

-- Plans
INSERT INTO Plans (plan_name, price) VALUES
('Basic', 0.00),
('Silver', 19.99),
('Gold', 29.99),
('Platinum', 49.99);

-- Campaigns
INSERT INTO Campaigns (name, description, code, start_date, end_date, max_vouchers, voucher_lifetime, discount_percentage) VALUES
('Welcome New Users', 'Welcome campaign for only 100 users', 'WELCOME2024', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP + INTERVAL '7 days', 100, 30, 10);

-- CampaignPlans
INSERT INTO CampaignPlans (campaign_id, plan_id) VALUES
((SELECT campaign_id FROM Campaigns WHERE name = 'Welcome New Users'), (SELECT plan_id FROM Plans WHERE plan_name = 'Silver'));