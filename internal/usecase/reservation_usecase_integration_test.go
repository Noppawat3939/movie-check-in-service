package usecase

import (
	"check-in/internal/domain"
	"check-in/internal/infra/postgresl"
	"check-in/internal/infra/redis"
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type testDeps struct {
	ctx             context.Context
	db              *gorm.DB
	client          *redis.Client
	lockRepo        *redis.LockRepository
	reservationRepo postgresl.ReservationRepository
	lockLogRepo     postgresl.ReservationLockLogRepository
	usecase         ReservationUsecase
}

func setupTestEnv(t *testing.T) {
	t.Setenv("DB_HOST", "localhost")
	t.Setenv("DB_EXTERNAL_PORT", "5434")
	t.Setenv("DB_USER", "movie_checkin")
	t.Setenv("DB_PASSWORD", "movie_checkin0")
	t.Setenv("DB_NAME", "movie_checkin")
	t.Setenv("DB_SSLMODE", "disable")

	t.Setenv("REDIS_HOST", "localhost")
	t.Setenv("REDIS_EXTERNAL_PORT", "6380")
}

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := postgresl.NewDB()
	assert.NoError(t, err)

	t.Cleanup(func() {
		postgresl.Close(db)
	})

	return db
}

func setupTestRedis(t *testing.T) *redis.Client {
	client, err := redis.NewClient()
	assert.NoError(t, err)
	t.Cleanup(func() {
		_ = client.Close()
	})

	return client
}

func cleanupReservationData(t *testing.T, db *gorm.DB, showtimeID uuid.UUID, seatID uuid.UUID) {
	err := db.Exec(`DELETE FROM reservations WHERE showtime_id = ? AND seat_id = ?`, showtimeID, seatID).Error

	assert.NoError(t, err)
}

func cleanupReservationByShowTimeID(t *testing.T, db *gorm.DB, showtimeID uuid.UUID) {
	err := db.Exec("DELETE FROM reservations WHERE showtime_id = ?", showtimeID).Error

	assert.NoError(t, err)
}

func setupReverationTest(t *testing.T) *testDeps {
	t.Helper()
	setupTestEnv(t)

	db := setupTestDB(t)
	client := setupTestRedis(t)
	reservationRepo := postgresl.NewReversationRepository(db)
	lockLogRepo := postgresl.NewReservationLockLogRepository(db)
	lockRepo := redis.NewLockRepository(client)

	uc := NewReservationUsecase(reservationRepo, *lockRepo, lockLogRepo)

	return &testDeps{
		ctx:             context.Background(),
		db:              db,
		client:          client,
		reservationRepo: reservationRepo,
		lockRepo:        lockRepo,
		lockLogRepo:     lockLogRepo,
		usecase:         uc,
	}
}

func toUUID(s string) uuid.UUID {
	parsed, _ := uuid.Parse(s)
	return parsed
}

func TestCreateReservation_ConcurrencyRequests(t *testing.T) {
	deps := setupReverationTest(t)
	reservationRepo := deps.reservationRepo
	uc := deps.usecase

	showtimeID := toUUID("b1000000-0000-0000-0000-000000000001")
	seatID := toUUID("c1000000-0000-0000-0000-000000000002")

	cleanupReservationData(t, deps.db, showtimeID, seatID)

	req := domain.CreateReservationRequest{ShowTimeID: showtimeID, SeatID: seatID}

	const numGoroutines = 100
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	var (
		successCount atomic.Int32
		failedCount  atomic.Int32
		wg           sync.WaitGroup
	)

	wg.Add(numGoroutines)

	for i := range numGoroutines {
		go func(userNum int) {
			defer wg.Done()

			_, err := uc.CreateReservation(ctx, req)
			if err == nil {
				successCount.Add(1)
			} else {
				failedCount.Add(1)
			}

		}(i)
	}

	wg.Wait()

	count, err := reservationRepo.CountByShowTimeAndSeat(ctx, showtimeID, seatID)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)
	assert.Equal(t, int32(1), successCount.Load())
	assert.Equal(t, int32(numGoroutines-1), failedCount.Load())
	assert.Equal(t, int32(numGoroutines), successCount.Load()+failedCount.Load())
}

func TestChangeReservation_Success(t *testing.T) {
	deps := setupReverationTest(t)

	reservationRepo := deps.reservationRepo

	showtimeID := toUUID("b1000000-0000-0000-0000-000000000001")
	oldSeatID := toUUID("c1000000-0000-0000-0000-000000000002")
	newSeatID := toUUID("c1000000-0000-0000-0000-000000000003")

	// cleanup data
	cleanupReservationByShowTimeID(t, deps.db, showtimeID)

	ctx := deps.ctx
	uc := deps.usecase

	// create existing reservation
	oldReserve := &domain.Reservation{
		ID:         uuid.New(),
		ShowTimeID: showtimeID,
		SeatID:     oldSeatID,
		Status:     domain.ReservationConfirmed,
		ReservedAt: time.Now(),
	}

	err := reservationRepo.Create(ctx, oldReserve)
	assert.NoError(t, err)

	// make request for change reservation
	req := domain.ChangeReservationRequest{
		ReservationID: oldReserve.ID,
		NewSeatID:     newSeatID,
	}

	resp, err := uc.ChangeReservation(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	// check response data
	assert.Equal(t, showtimeID, resp.ShowtimeID)
	assert.Equal(t, newSeatID, resp.SeatID)
	assert.Equal(t, domain.ReservationConfirmed, resp.Status)

	// verify old reservation cancelled
	oldData, err := reservationRepo.FindByID(ctx, oldReserve.ID)
	assert.NoError(t, err)
	assert.Equal(t, domain.ReservationCancelled, oldData.Status)

	// verify new reservation exists (expected count = 1)
	count, err := reservationRepo.CountByShowTimeAndSeat(ctx, showtimeID, newSeatID)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)
}

func TestChangeReservation_NewSeatAlreadyReserved(t *testing.T) {
	deps := setupReverationTest(t)
	ctx := deps.ctx
	reservationRepo := deps.reservationRepo
	uc := deps.usecase

	showtimeID := toUUID("b1000000-0000-0000-0000-000000000001")
	oldSeatID := toUUID("c1000000-0000-0000-0000-000000000002")
	newSeatID := toUUID("c1000000-0000-0000-0000-000000000003")

	// clean up
	cleanupReservationByShowTimeID(t, deps.db, showtimeID)

	// create 2 reserve (old A and old B , then change new seat same seat with old B)
	oldA := &domain.Reservation{
		ID:         uuid.New(),
		ShowTimeID: showtimeID,
		SeatID:     oldSeatID,
		Status:     domain.ReservationConfirmed,
		ReservedAt: time.Now(),
	}

	oldB := &domain.Reservation{
		ID:         uuid.New(),
		ShowTimeID: showtimeID,
		SeatID:     newSeatID, // same with new seat
		Status:     domain.ReservationConfirmed,
		ReservedAt: time.Now(),
	}

	assert.NoError(t, reservationRepo.Create(ctx, oldA), oldA)
	assert.NoError(t, reservationRepo.Create(ctx, oldB), oldB)

	// make request for change reservation (change seat from oldA to new A but seat new A is reserved by old B)
	req := domain.ChangeReservationRequest{
		NewSeatID:     newSeatID,
		ReservationID: oldA.ID,
	}

	resp, err := uc.ChangeReservation(ctx, req)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, domain.ErrSeatAlreadyReserved)
	assert.Nil(t, resp)
}

func TestChangeReservation_ReservationNotFound(t *testing.T) {
	deps := setupReverationTest(t)

	uc := deps.usecase
	newSeatID := toUUID("c1000000-0000-0000-0000-000000000003") // seat existing
	newReserveID := uuid.New()

	// make request for change reservation (reservation id not found in db)
	req := domain.ChangeReservationRequest{
		ReservationID: newReserveID,
		NewSeatID:     newSeatID,
	}

	resp, err := uc.ChangeReservation(deps.ctx, req)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, domain.ErrReservationNotFound)
	assert.Nil(t, resp)
}
