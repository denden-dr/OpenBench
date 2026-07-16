-- Drop foreign key constraint on claims table
ALTER TABLE claims DROP CONSTRAINT IF EXISTS claims_warranty_id_fkey;

-- Revert column types to VARCHAR(100)
ALTER TABLE warranties ALTER COLUMN id TYPE VARCHAR(100);
ALTER TABLE claims ALTER COLUMN warranty_id TYPE VARCHAR(100);
ALTER TABLE claims ALTER COLUMN id TYPE VARCHAR(100);

-- Recreate foreign key constraint
ALTER TABLE claims ADD CONSTRAINT claims_warranty_id_fkey FOREIGN KEY (warranty_id) REFERENCES warranties(id) ON DELETE CASCADE;
