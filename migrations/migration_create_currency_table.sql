-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create currencies table
CREATE TABLE currencies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(3) UNIQUE NOT NULL,
    description VARCHAR(255) NOT NULL,
    amount_display_format VARCHAR(50) DEFAULT '###,###.##',
    html_encoded_symbol VARCHAR(50),
    factor INTEGER DEFAULT 100,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_by UUID NOT NULL
);

-- Create indexes
CREATE INDEX idx_currencies_code ON currencies(code);
CREATE INDEX idx_currencies_created_at ON currencies(created_at);

-- Add comments
COMMENT ON TABLE currencies IS 'Currency master data with display formatting and symbols';
COMMENT ON COLUMN currencies.code IS 'ISO 4217 currency code';
COMMENT ON COLUMN currencies.factor IS 'Factor for decimal precision (100 = 2 decimal places, 1000 = 3 decimal places)';
COMMENT ON COLUMN currencies.html_encoded_symbol IS 'HTML encoded currency symbol for display';