CREATE SCHEMA IF NOT EXISTS kitchen;

CREATE TABLE kitchen.restaurants (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    address TEXT NOT NULL,

    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,

    delivery_radius_m INTEGER NOT NULL,

    opening_time TIME NOT NULL,
    closing_time TIME NOT NULL,

    CONSTRAINT restaurants_latitude_check
        CHECK (latitude >= -90 AND latitude <= 90),

    CONSTRAINT restaurants_longitude_check
        CHECK (longitude >= -180 AND longitude <= 180),

    CONSTRAINT restaurants_delivery_radius_check
        CHECK (delivery_radius_m >= 0)
);
