-- Drop new trigger and function
DROP TRIGGER IF EXISTS job_update_notify ON jobs;
DROP FUNCTION IF EXISTS notify_job_update();
