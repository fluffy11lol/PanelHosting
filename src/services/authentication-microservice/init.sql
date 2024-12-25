DO $$
BEGIN
   IF NOT EXISTS (
       SELECT FROM auth_db
       WHERE datname = 'auth_db'
   ) THEN
       CREATE DATABASE auth_db;
END IF;
END
$$;
