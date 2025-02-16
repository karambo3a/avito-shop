CREATE TABLE users
(
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    coins INTEGER DEFAULT 1000 CHECK (coins >= 0)
);

CREATE TABLE items
(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price INTEGER NOT NULL CHECK (price >= 0)
);

INSERT INTO items (name, price) VALUES ('t-shirt', 80);
INSERT INTO items (name, price) VALUES ('cup', 20);
INSERT INTO items (name, price) VALUES ('book', 50);
INSERT INTO items (name, price) VALUES ('pen', 10);
INSERT INTO items (name, price) VALUES ('powerbank', 200);
INSERT INTO items (name, price) VALUES ('hoody', 300);
INSERT INTO items (name, price) VALUES ('umbrella', 200);
INSERT INTO items (name, price) VALUES ('socks', 10);
INSERT INTO items (name, price) VALUES ('wallet', 50);
INSERT INTO items (name, price) VALUES ('pink-hoody', 500);

CREATE TABLE sales
(
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    item_id INTEGER REFERENCES items(id) ON DELETE CASCADE
);

CREATE TABLE transactions
(
    id SERIAL PRIMARY KEY,
    from_user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    to_user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    amount INTEGER NOT NULL CHECK (amount > 0)
);
