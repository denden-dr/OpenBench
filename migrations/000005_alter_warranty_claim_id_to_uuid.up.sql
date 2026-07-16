-- Clean up any invalid non-UUID records first to prevent cast failure
DELETE FROM claims WHERE id !~ '^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$';
DELETE FROM warranties WHERE id !~ '^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$';

-- Drop foreign key constraint on claims table
ALTER TABLE claims DROP CONSTRAINT IF EXISTS claims_warranty_id_fkey;

-- Alter warranties.id to UUID
ALTER TABLE warranties ALTER COLUMN id TYPE UUID USING id::uuid;

-- Alter claims.warranty_id to UUID
ALTER TABLE claims ALTER COLUMN warranty_id TYPE UUID USING warranty_id::uuid;

-- Alter claims.id to UUID
ALTER TABLE claims ALTER COLUMN id TYPE UUID USING id::uuid;

-- Recreate the foreign key constraint
ALTER TABLE claims ADD CONSTRAINT claims_warranty_id_fkey FOREIGN KEY (warranty_id) REFERENCES warranties(id) ON DELETE CASCADE;
