CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY,
    sender_account_id UUID NOT NULL,
    receiver_account_id UUID NOT NULL,
    amount FLOAT,
    currency VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL
);
