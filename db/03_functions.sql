-- Create functions

CREATE FUNCTION get_user_timezone(discord_id_param VARCHAR(20))
    RETURNS VARCHAR(32)
    LANGUAGE plpgsql
    SECURITY DEFINER
AS $$
DECLARE
    timezone VARCHAR(32);
BEGIN
    SELECT timezones.name
    INTO timezone
    FROM users
    JOIN user_timezones ON users.id = user_timezones.user_id
    JOIN timezones ON user_timezones.timezone_id = timezones.id
    WHERE users.discord_id = discord_id_param;
   
    RETURN timezone;
END;$$;
