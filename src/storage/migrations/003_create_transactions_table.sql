CREATE IF NOT EXISTS transacations (
    id UUID PRIMARY KEY,
    sender_account_id UUID NOT NULL,
    receiver_account_id UUID NOT NULL,
    amount FLOAT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL
);