-- CREATE TABLE IF NOT EXISTS slots (
--   id SERIAL PRIMARY KEY,
--   user_id UUID NULL,
--   email TEXT NULL,
--   start_time TIMESTAMP NULL,
--   end_time TIMESTAMP NULL,
--   notified BOOLEAN DEFAULT FALSE
-- );

-- -- Seed a few empty slots
-- INSERT INTO slots (user_id, email, start_time, end_time, notified)
-- SELECT NULL, NULL, NULL, NULL, FALSE
-- WHERE NOT EXISTS (SELECT 1 FROM slots);

-- -- Optionally add more initial empty rows
-- INSERT INTO slots (user_id, email, start_time, end_time, notified)
-- SELECT NULL, NULL, NULL, NULL, FALSE
-- FROM generate_series(1, 5)
-- WHERE (SELECT COUNT(*) FROM slots) < 6;