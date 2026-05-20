CREATE TABLE transactions (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
  amount bigint NOT NULL CHECK (amount > 0),
  currency varchar(3) NOT NULL DEFAULT 'NGN',
  status varchar(20) NOT NULL DEFAULT 'pending',
  type varchar(20) NOT NULL,
  reference varchar(255) NOT NULL UNIQUE,
  description text,
  provider varchar(50),
  metadata jsonb,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT transactions_status_check CHECK (
    status IN ('pending', 'completed', 'failed', 'cancelled')
  ),
  CONSTRAINT transactions_type_check CHECK (
    type IN ('credit', 'debit', 'transfer', 'payment')
  )
);

CREATE INDEX idx_transactions_user_id ON transactions (user_id);
CREATE INDEX idx_transactions_status ON transactions (status);
CREATE INDEX idx_transactions_created_at ON transactions (created_at DESC);

COMMENT ON TABLE transactions IS 'User payment and wallet transactions';
COMMENT ON COLUMN transactions.amount IS 'Amount in smallest currency unit (e.g. kobo for NGN)';
COMMENT ON COLUMN transactions.reference IS 'Unique external or internal payment reference';
COMMENT ON COLUMN transactions.metadata IS 'Optional provider-specific payload';
