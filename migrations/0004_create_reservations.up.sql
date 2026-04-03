-- create enum
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_type WHERE typname = 'reservation_status'
    ) THEN
        CREATE TYPE reservation_status AS ENUM (
            'confirmed',
            'cancelled',
            'expired'
            );
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS reservations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    showtime_id UUID NOT NULL REFERENCES showtimes(id),
    seat_id UUID NOT NULL REFERENCES seats(id),
    status reservation_status NOT NULL DEFAULT 'confirmed',
    reserved_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_reservations_showtime_seat
ON reservations(showtime_id, seat_id);