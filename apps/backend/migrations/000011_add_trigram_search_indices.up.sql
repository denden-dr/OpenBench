DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_extension WHERE extname = 'pg_trgm') THEN
        BEGIN
            CREATE EXTENSION pg_trgm;
        EXCEPTION
            WHEN OTHERS THEN
                RAISE WARNING 'pg_trgm extension could not be created. Please ensure it is enabled: %', SQLERRM;
        END;
    END IF;
END $$;

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
