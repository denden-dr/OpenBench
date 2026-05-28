CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX IF NOT EXISTS tickets_search_trgm_idx ON tickets USING gin (
  (lower(
    COALESCE(id::text, '') || ' ' ||
    COALESCE(customer_name, '') || ' ' ||
    COALESCE(customer_phone, '') || ' ' ||
    COALESCE(brand, '') || ' ' ||
    COALESCE(model, '') || ' ' ||
    COALESCE(issue, '')
  )) gin_trgm_ops
);
