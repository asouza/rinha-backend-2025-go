-- Create transactions table for Rinha Backend Go project
-- This script creates the necessary table structure for storing transaction data

CREATE TABLE IF NOT EXISTS transactions (
    id VARCHAR(36) PRIMARY KEY,
    valor_centavos BIGINT NOT NULL,
    instante TIMESTAMP WITH TIME ZONE NOT NULL,
    url VARCHAR(255) NOT NULL,
    is_priority BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create index on instante for efficient date range queries
CREATE INDEX IF NOT EXISTS idx_transactions_instante ON transactions(instante);

-- Create index on is_priority for efficient summary queries
CREATE INDEX IF NOT EXISTS idx_transactions_is_priority ON transactions(is_priority);

-- Create composite index for summary queries with date filtering
CREATE INDEX IF NOT EXISTS idx_transactions_summary ON transactions(instante, is_priority);