CREATE TABLE IF NOT EXISTS accounts (
  id TEXT PRIMARY KEY,
  balance NUMERIC NOT NULL
);

CREATE TABLE IF NOT EXISTS transactions (
  id TEXT,
  end_to_end_id TEXT UNIQUE,
  sender_account TEXT,
  receiver_account TEXT,
  amount NUMERIC,
  status TEXT,
  created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_tx_end ON transactions(end_to_end_id);