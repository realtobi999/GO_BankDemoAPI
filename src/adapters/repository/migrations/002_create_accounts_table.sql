CREATE TABLE IF NOT EXISTS accounts (
    id UUID PRIMARY KEY,
    customer_id UUID REFERENCES customers(id) NOT NULL,
    balance FLOAT,
    account_type VARCHAR(255),
    currency VARCHAR(255),
    status BOOLEAN,
    opening_date TIMESTAMP WITH TIME ZONE,
    last_transaction_date TIMESTAMP WITH TIME ZONE,
    interest_rate FLOAT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL
);
