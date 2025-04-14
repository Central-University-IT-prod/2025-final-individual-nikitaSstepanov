-- +goose Up
-- +goose StatementBegin
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'gender') THEN
        CREATE TYPE gender AS ENUM (
            'MALE',
            'FEMALE',
            'ALL'
        );
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS clients (
    id UUID PRIMARY KEY,
    login VARCHAR(255),
    age INT,
    location VARCHAR(255),
    gender gender
);

CREATE TABLE IF NOT EXISTS advertisers (
    id UUID PRIMARY KEY,
    name VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS ml_scores (
    client_id UUID,
    advertiser_id UUID,
    score INT,
    FOREIGN KEY (client_id) REFERENCES clients (id) ON DELETE CASCADE,
    FOREIGN KEY (advertiser_id) REFERENCES advertisers (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS campaigns (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    title VARCHAR(255),
    text TEXT,
    start_date INT,
    end_date INT,
    advertiser_id UUID,
    image VARCHAR(255),
    FOREIGN KEY (advertiser_id) REFERENCES advertisers (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS targeting (
    campaign_id UUID PRIMARY KEY,
    gender gender,
    age_from INT,
    age_to INT,
    location VARCHAR(255),
    FOREIGN KEY (campaign_id) REFERENCES campaigns (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS billing (
    campaign_id UUID PRIMARY KEY,
    impressions_limit INT,
    clicks_limit INT,
    cost_per_impression REAL,
    cost_per_click REAL,
    impressions_count INT DEFAULT 0,
    clicks_count INT DEFAULT 0,
    spent_impressions REAL DEFAULT 0,
    spent_clicks REAL DEFAULT 0,
    FOREIGN KEY (campaign_id) REFERENCES campaigns (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS daily_billing (
    campaign_id UUID,
    date INT,
    impressions_count INT DEFAULT 0,
    clicks_count INT DEFAULT 0,
    spent_impressions REAL DEFAULT 0,
    spent_clicks REAL DEFAULT 0,
    FOREIGN KEY (campaign_id) REFERENCES campaigns (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS impressions (
    client_id UUID,
    campaign_id UUID,
    FOREIGN KEY (client_id) REFERENCES clients (id) ON DELETE CASCADE,
    FOREIGN KEY (campaign_id) REFERENCES campaigns (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS clicks (
    client_id UUID,
    campaign_id UUID,
    FOREIGN KEY (client_id) REFERENCES clients (id) ON DELETE CASCADE,
    FOREIGN KEY (campaign_id) REFERENCES campaigns (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS clients;

DROP TABLE IF EXISTS advertisers;

DROP TABLE IF EXISTS ml_scores;

DROP TABLE IF NOT EXISTS campaigns;

DROP TABLE IF EXISTS targeting;

DROP TABLE IF EXISTS billing;

DROP TABLE IF EXISTS daily_billing;

DROP TYPE IF EXISTS gender;
-- +goose StatementEnd
