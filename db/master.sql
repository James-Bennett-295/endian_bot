-- Create and connect to database "endian_bot"
CREATE DATABASE endian_bot;
\connect endian_bot;

-- Execute SQL files
\i 01_tables.sql;
\i 02_data.sql;
\i 03_functions.sql;
\i 04_user.sql;
