-- rollback partial unique index
DROP INDEX IF EXISTS uq_reservations_showtime_seat_confirmed;

-- restore original unique index
CREATE UNIQUE INDEX IF NOT EXISTS uq_reservations_showtime_seat
ON reservations(showtime_id, seat_id);