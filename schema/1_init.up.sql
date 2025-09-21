CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    telegram_username VARCHAR(255) NOT NULL UNIQUE,
    telegram_id VARCHAR(255) NOT NULL UNIQUE,
    coins_balance INT NOT NULL DEFAULT 0,
    coins_per_tap INT NOT NULL DEFAULT 1,
    level INT NOT NULL DEFAULT 1,
    energy INT NOT NULL,
    max_energy INT NOT NULL,
    last_energy_update TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    user_ref VARCHAR(255) NOT NULL,
    by_ref VARCHAR(255),
    passive_income INT NOT NULL DEFAULT 0,
    last_passive_income_update TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE notifications (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    text VARCHAR(255) NOT NULL
);

CREATE TABLE upgrades (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description VARCHAR(255) NOT NULL,
    coins_price INT NOT NULL,
    stars_price INT NOT NULL,
    new_max_energy INT NOT NULL,
    new_one_tap_coins INT NOT NULL,
    new_passive_income INT NOT NULL
);

CREATE TABLE user_upgrades (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    upgrade_id INT NOT NULL
);

ALTER TABLE notifications ADD FOREIGN KEY (user_id) REFERENCES users(id);
ALTER TABLE user_upgrades ADD FOREIGN KEY (user_id) REFERENCES users(id);
ALTER TABLE user_upgrades ADD FOREIGN KEY (upgrade_id) REFERENCES upgrades(id);