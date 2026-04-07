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

func TestCreateReservation_ConcurrencyRequests(t *testing.T) {
	setupTestEnv(t)

	db := setupTestDB(t)
	client := setupTestRedis(t)

	reservationRepo := postgresl.NewReversationRepository(db)
	lockRepo := redis.NewLockRepository(client)

	uc := NewReservationUsecase(reservationRepo, *lockRepo)
	showtimeID, _ := uuid.Parse("b1000000-0000-0000-0000-000000000001")
	seatID, _ := uuid.Parse("c1000000-0000-0000-0000-000000000002")

	cleanupReservationData(t, db, showtimeID, seatID)

	req := domain.CreateReservationRequest{ShowTimeID: showtimeID, SeatID: seatID}

	const numGoroutines = 1000
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	var (
		successCount atomic.Int32
		failedCount  atomic.Int32
		wg           sync.WaitGroup
	)

	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
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
