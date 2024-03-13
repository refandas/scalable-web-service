CREATE TABLE IF NOT EXISTS public.items (
    id BIGSERIAL,
    code VARCHAR(10),
    description VARCHAR(50),
    quantity BIGINT,
    order_id BIGINT,
    CONSTRAINT items_pkey PRIMARY KEY (id),
    CONSTRAINT items_orders_fkey FOREIGN KEY (order_id) REFERENCES orders(id)
);

CREATE TABLE IF NOT EXISTS public.orders (
    id BIGSERIAL,
    customer_name VARCHAR(50),
    ordered_at TIMESTAMP,
    CONSTRAINT orders_pkey PRIMARY KEY (id)
);
