-- Create user and grant permissions

CREATE USER endian_bot WITH PASSWORD 'DEFAULT_PASSWORD';

GRANT CONNECT ON DATABASE endian_bot TO endian_bot;

GRANT EXECUTE ON FUNCTION get_user_timezone(varchar) TO endian_bot;
