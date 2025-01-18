CREATE TYPE preference_enum AS ENUM ('CARDIO', 'WEIGHT');
CREATE TYPE weight_unit_enum AS ENUM ('KG', 'LBS');
CREATE TYPE height_unit_enum AS ENUM ('CM', 'INCH');

CREATE TABLE Users (
    id VARCHAR(255) NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    preference preference_enum,
    weight_unit weight_unit_enum,
    height_unit height_unit_enum,
    weight DECIMAL(5,2) CHECK (weight BETWEEN 10 AND 1000),
    height DECIMAL(5,2) CHECK (height BETWEEN 3 AND 250),
    name VARCHAR(60),
    image_uri TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE Activities (
    id VARCHAR(255) NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    done_at TIMESTAMP,
    duration_in_minutes INT CHECK (duration_in_minutes >= 1),
    calories_burned DECIMAL(10,2),
    activity_type VARCHAR(10),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES Users (id) ON DELETE CASCADE
);
