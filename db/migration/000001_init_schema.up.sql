CREATE TABLE UserTypes (
    user_type_id SERIAL PRIMARY KEY,
    type_name VARCHAR(50) NOT NULL UNIQUE
);

CREATE TABLE Users (
    user_id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    registration_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_login TIMESTAMP,
    user_type_id INT REFERENCES UserTypes(user_type_id)
);

CREATE INDEX idx_users_user_type_id ON Users(user_type_id);

CREATE TABLE Plans (
    plan_id SERIAL PRIMARY KEY,
    plan_name VARCHAR(50) NOT NULL UNIQUE,
    price DECIMAL(10, 2) NOT NULL
);

CREATE TABLE Campaigns (
    campaign_id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    max_users INT NOT NULL,
    voucher_duration INT NOT NULL DEFAULT 1440,  -- Default to 24 hours in minutes
    discount_percentage DECIMAL(5, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_campaigns_start_date ON Campaigns(start_date);
CREATE INDEX idx_campaigns_end_date ON Campaigns(end_date);

CREATE TABLE CampaignPlans (
    campaign_id INT REFERENCES Campaigns(campaign_id) ON DELETE CASCADE,
    plan_id INT REFERENCES Plans(plan_id) ON DELETE CASCADE,
    PRIMARY KEY (campaign_id, plan_id)
);

CREATE TABLE Vouchers (
    voucher_id SERIAL PRIMARY KEY,
    user_id INT REFERENCES Users(user_id) ON DELETE CASCADE,
    campaign_id INT REFERENCES Campaigns(campaign_id) ON DELETE CASCADE,
    code VARCHAR(50) NOT NULL UNIQUE,
    valid_until TIMESTAMP,
    is_redeemed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_vouchers_user_id ON Vouchers(user_id);
CREATE INDEX idx_vouchers_campaign_id ON Vouchers(campaign_id);
CREATE INDEX idx_vouchers_valid_until ON Vouchers(valid_until);

CREATE TABLE Purchases (
    purchase_id SERIAL PRIMARY KEY,
    user_id INT REFERENCES Users(user_id),
    voucher_id INT REFERENCES Vouchers(voucher_id),
    plan_id INT REFERENCES Plans(plan_id),
    amount DECIMAL(10, 2) NOT NULL,
    discount_applied BOOLEAN DEFAULT FALSE,
    purchase_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_purchases_user_id ON Purchases(user_id);
CREATE INDEX idx_purchases_voucher_id ON Purchases(voucher_id);
CREATE INDEX idx_purchases_plan_id ON Purchases(plan_id);
CREATE INDEX idx_purchases_purchase_date ON Purchases(purchase_date);