-- drop old unique index
DROP INDEX IF EXISTS uq_reservations_showtime_seat;

-- create partial unique index
-- allow only 1 active (confirmed) reservation per seat per showtime
CREATE UNIQUE INDEX IF NOT EXISTS uq_reservations_showtime_seat_confirmed
ON reservations(showtime_id, seat_id)
WHERE status = 'confirmed';