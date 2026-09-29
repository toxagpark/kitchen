CREATE TABLE kitchen.orders (
    id UUID PRIMARY KEY,
    restaurant_id UUID NOT NULL,

    status TEXT NOT NULL,
    total_price BIGINT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT orders_restaurant_fk
        FOREIGN KEY (restaurant_id)
        REFERENCES kitchen.restaurants (id),

    CONSTRAINT orders_status_check
        CHECK (status IN ('pending', 'accepted', 'rejected')),

    CONSTRAINT orders_total_price_check
        CHECK (total_price >= 0)
);

CREATE TABLE kitchen.order_items (
    id UUID PRIMARY KEY,
    order_id UUID NOT NULL,
    menu_item_id UUID NOT NULL,

    name TEXT NOT NULL,

    quantity INTEGER NOT NULL,
    unit_price BIGINT NOT NULL,
    total_price BIGINT NOT NULL,

    CONSTRAINT order_items_order_fk
        FOREIGN KEY (order_id)
        REFERENCES kitchen.orders (id)
        ON DELETE CASCADE,

    CONSTRAINT order_items_menu_item_fk
        FOREIGN KEY (menu_item_id)
        REFERENCES kitchen.menu_items (id),

    CONSTRAINT order_items_quantity_check
        CHECK (quantity > 0),

    CONSTRAINT order_items_unit_price_check
        CHECK (unit_price >= 0),

    CONSTRAINT order_items_total_price_check
        CHECK (total_price >= 0)
);

CREATE INDEX idx_orders_restaurant_id
    ON kitchen.orders (restaurant_id);

CREATE INDEX idx_order_items_order_id
    ON kitchen.order_items (order_id);
