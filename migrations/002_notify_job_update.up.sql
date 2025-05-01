-- Drop old trigger if it exists
DROP TRIGGER IF EXISTS job_update_notify ON jobs;

-- Drop old function if it exists
DROP FUNCTION IF EXISTS notify_job_update();

-- Create the new function
CREATE OR REPLACE FUNCTION notify_job_update()
RETURNS trigger AS $$
BEGIN
    PERFORM pg_notify('job_updates', NEW.id::text);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create the new trigger
CREATE TRIGGER job_update_notify
AFTER INSERT OR UPDATE ON jobs
FOR EACH ROW
EXECUTE FUNCTION notify_job_update();
