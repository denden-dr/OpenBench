DROP INDEX IF EXISTS idx_tickets_id_text_prefix;

ALTER TABLE tickets
  DROP COLUMN customer_phone;
