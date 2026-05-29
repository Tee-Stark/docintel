-- Migration: add_name_to_users (down)

ALTER TABLE users DROP COLUMN IF EXISTS name;
