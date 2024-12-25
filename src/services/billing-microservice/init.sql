DO $$
BEGIN
   IF NOT EXISTS (
       SELECT FROM billing_db
       WHERE datname = 'billing_db'
   ) THEN
       CREATE DATABASE billing_db;
END IF;
END
$$;
