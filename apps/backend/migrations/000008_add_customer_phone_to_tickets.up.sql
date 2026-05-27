ALTER TABLE tickets
  ADD COLUMN customer_phone TEXT NOT NULL DEFAULT '';

CREATE INDEX IF NOT EXISTS idx_tickets_id_text_prefix
  ON tickets ((left(id::text, 8)));
