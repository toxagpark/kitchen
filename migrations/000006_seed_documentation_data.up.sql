INSERT INTO kitchen.restaurants (
    id,
    name,
    address,
    latitude,
    longitude,
    delivery_radius_m,
    opening_time,
    closing_time
) VALUES (
    '83b417b6-e153-4a09-b26a-4f226830980d',
    'Avito Kitchen Demo',
    'ул. Документации, 1',
    55.7558,
    37.6173,
    5000,
    '09:00',
    '23:00'
);

INSERT INTO kitchen.menu_items (
    id,
    restaurant_id,
    name,
    description,
    price,
    is_available
) VALUES (
    '117f50a0-7719-47f0-b63c-80c28224fed3',
    '83b417b6-e153-4a09-b26a-4f226830980d',
    'Демо-блюдо',
    'Блюдо для проверки API',
    29900,
    TRUE
);
