CREATE TABLE Users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    preference ENUM('CARDIO', 'WEIGHT') NOT NULL,
    weight_unit ENUM('KG', 'LBS') NOT NULL,
    height_unit ENUM('CM', 'INCH') NOT NULL,
    weight DECIMAL(5,2) NOT NULL CHECK (weight BETWEEN 10 AND 1000),
    height DECIMAL(5,2) NOT NULL CHECK (height BETWEEN 3 AND 250),
    name VARCHAR(60),
    image_uri TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE ActivityTypes (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    calories_per_minute DECIMAL(5,2) NOT NULL
);

CREATE TABLE Activities (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    activity_type_id INT NOT NULL,
    done_at TIMESTAMP NOT NULL,
    duration_in_minutes INT NOT NULL CHECK (duration_in_minutes >= 1),
    calories_burned DECIMAL(10,2) GENERATED ALWAYS AS (
        duration_in_minutes * (SELECT calories_per_minute FROM ActivityTypes WHERE id = activity_type_id)
    ) STORED,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES Users (id) ON DELETE CASCADE,
    FOREIGN KEY (activity_type_id) REFERENCES ActivityTypes (id) ON DELETE CASCADE
);
