-- create enum
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_type WHERE typname = 'lock_status'
    ) THEN
        CREATE TYPE lock_status AS ENUM ('acquired','failed','released');
    END IF;
END $$;

-- create a table
CREATE TABLE IF NOT EXISTS reservation_lock_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    showtime_id UUID NOT NULL,
    seat_id UUID NOT NULL,
    lock_key VARCHAR(255) NOT NULL,
    status lock_status NOT NULL,
    acquired_at TIMESTAMPTZ,
    released_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);