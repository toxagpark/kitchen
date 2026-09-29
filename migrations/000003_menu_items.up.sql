CREATE TABLE kitchen.menu_items (
    id UUID PRIMARY KEY,
    restaurant_id UUID NOT NULL,

    name TEXT NOT NULL,
    description TEXT NOT NULL,

    price BIGINT NOT NULL,
    is_available BOOLEAN NOT NULL DEFAULT TRUE,

    CONSTRAINT menu_items_restaurant_fk
        FOREIGN KEY (restaurant_id)
        REFERENCES kitchen.restaurants (id)
        ON DELETE CASCADE,

    CONSTRAINT menu_items_price_check
        CHECK (price >= 0)
);

CREATE INDEX idx_menu_items_restaurant_id
    ON kitchen.menu_items (restaurant_id);
