BEGIN;

CREATE TABLE IF NOT EXISTS users
(
    userid          UUID PRIMARY KEY,
    name        TEXT,
    email       TEXT      NOT NULL,
    password    TEXT      NOT NULL,
    role        TEXT NOT NULL,
    createdby  UUID REFERENCES users (userid),
    createdat  TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updatedby  UUID REFERENCES users (userid),
    updatedat  TIMESTAMP WITH TIME ZONE,
    archivedat TIMESTAMP WITH TIME ZONE
);
CREATE UNIQUE INDEX IF NOT EXISTS unique_user ON users (email) WHERE archivedat IS NULL;


CREATE TABLE IF NOT EXISTS usersession
(
    sessionid          UUID PRIMARY KEY,
    userid     UUID REFERENCES users (userid) NOT NULL,
    createdat  TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    archivedat TIMESTAMP WITH TIME ZONE
);


CREATE TABLE IF NOT EXISTS address
(
    addressid          UUID PRIMARY KEY,
    addressline     TEXT                       NOT NULL,
    latitude    DOUBLE PRECISION,
    longitude   DOUBLE PRECISION,
    user_id     UUID REFERENCES users (userid) NOT NULL,
    createdat  TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    archivedat TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX IF NOT EXISTS unique_address
    ON address (user_id, addressline)
    WHERE archivedat IS NULL;


CREATE TABLE IF NOT EXISTS restaurants
(
    restaurantid          UUID PRIMARY KEY,
    name        TEXT                       NOT NULL,
    addressline     TEXT                       NOT NULL,
    latitude    DOUBLE PRECISION           ,
    longitude   DOUBLE PRECISION           ,
    createdby  UUID REFERENCES users (userid) NOT NULL,
    createdat  TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    archivedat TIMESTAMP WITH TIME ZONE
);
CREATE UNIQUE INDEX IF NOT EXISTS unique_restaurant ON restaurants (name, addressline) WHERE archivedat IS NULL;


CREATE TABLE IF NOT EXISTS dishes
(
    dishid            UUID PRIMARY KEY,
    name          TEXT                             NOT NULL,
    price         INTEGER                          NOT NULL,
    restaurantid UUID REFERENCES restaurants (restaurantid) NOT NULL,
    createdat    TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    archivedat   TIMESTAMP WITH TIME ZONE
);
CREATE UNIQUE INDEX IF NOT EXISTS unique_dish ON dishes (restaurantid, name) WHERE archivedat IS NULL;

COMMIT;