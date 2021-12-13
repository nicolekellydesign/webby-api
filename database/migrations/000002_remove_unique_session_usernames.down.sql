ALTER TABLE sessions
ADD CONSTRAINT sessions_user_name_key UNIQUE USING INDEX sessions_user_name_key;
ALTER TABLE sessions
ADD CONSTRAINT sessions_user_id_key UNIQUE USING INDEX sessions_user_id_key;