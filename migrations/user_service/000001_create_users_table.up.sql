CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    name text NOT NULL,
    email text NOT NULL,
    password text NOT NULL,
    age integer NOT NULL,
    version integer NOT NULL DEFAULT 1
);

-- migrate create -seq -ext=".sql" -dir="./migrations/user_service" create_users_table

-- migrate -path="./migrations/user_service" -database="postgres://postgres:olzhas4@localhost/sneakershop?sslmode=disable" down