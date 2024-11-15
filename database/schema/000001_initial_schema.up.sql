
CREATE TABLE IF NOT EXISTS users
(
    Id          UUID PRIMARY KEY,
    name        TEXT,
    email       TEXT NOT NULL,
    password    TEXT NOT NULL,
    role        TEXT NOT NULL,
    createdBy   UUID REFERENCES users (Id),
    createdAt   TIMESTAMP DEFAULT NOW(),
    updatedBy   UUID REFERENCES users (Id),
    updatedAt   TIMESTAMP DEFAULT NOW(),
    archivedAt  TIMESTAMP DEFAULT NULL
    );
CREATE UNIQUE INDEX IF NOT EXISTS unique_user ON users (email) WHERE archivedAt IS NULL;


CREATE TABLE IF NOT EXISTS user_session
(
    Id         UUID PRIMARY KEY,
    userId     UUID REFERENCES users (Id) NOT NULL,
    createdAt  TIMESTAMP DEFAULT NOW(),
    archivedAt TIMESTAMP DEFAULT NULL
);


CREATE TABLE IF NOT EXISTS addresses
(
    Id          UUID PRIMARY KEY,
    address     TEXT NOT NULL,
    latitude    DOUBLE PRECISION,
    longitude   DOUBLE PRECISION,
    userId      UUID REFERENCES users (Id) NOT NULL,
    createdAt   TIMESTAMP DEFAULT NOW(),
    archivedAt  TIMESTAMP DEFAULT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS unique_addresses
    ON addresses (userId, address)
    WHERE archivedAt IS NULL;


CREATE TABLE IF NOT EXISTS restaurants
(
    Id            UUID PRIMARY KEY,
    name          TEXT NOT NULL,
    address       TEXT NOT NULL,
    latitude      DOUBLE PRECISION,
    longitude     DOUBLE PRECISION,
    createdBy     UUID REFERENCES users (Id) NOT NULL,
    createdAt     TIMESTAMP DEFAULT NOW(),
    archivedAt    TIMESTAMP DEFAULT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS unique_restaurant ON restaurants (name, address) WHERE archivedAt IS NULL;


CREATE TABLE IF NOT EXISTS dishes
(
    Id            UUID PRIMARY KEY,
    name          TEXT NOT NULL,
    price         INTEGER NOT NULL,
    restaurantId  UUID REFERENCES restaurants (Id) NOT NULL,
    createdAt     TIMESTAMP DEFAULT NOW(),
    archivedAt    TIMESTAMP DEFAULT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS unique_dish ON dishes (restaurantId, name) WHERE archivedAt IS NULL;
