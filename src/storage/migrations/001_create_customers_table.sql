CREATE TABLE IF NOT EXISTS customers (
		id UUID PRIMARY KEY,
		first_name VARCHAR(255),
		last_name VARCHAR(255),
		birthday DATE,
		email VARCHAR(255),
		phone VARCHAR(100),
		state VARCHAR(255),
		address TEXT,
		created_at TIMESTAMP WITH TIME ZONE
);