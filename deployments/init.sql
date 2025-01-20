CREATE TABLE IF NOT EXISTS wallets
(
    id      SERIAL PRIMARY KEY,
    address UUID NOT NULL UNIQUE,
    balance NUMERIC NOT NULL
);

CREATE TABLE IF NOT EXISTS transactions
(
    id           SERIAL PRIMARY KEY,
    from_address UUID REFERENCES wallets(address),
    to_address   UUID REFERENCES wallets(address),
    timestamp    TIMESTAMP NOT NULL,
    amount       NUMERIC NOT NULL,
    successful   BOOLEAN NOT NULL
);
