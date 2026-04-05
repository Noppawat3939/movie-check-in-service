package http

import (
	"check-in/internal/delivery/http/handler"
	"check-in/internal/delivery/http/response"
	"check-in/internal/infra/postgresl"
	"check-in/internal/infra/redis"
	"check-in/internal/usecase"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	client, _ := redis.NewClient()
	lockRepo := redis.NewLockRepository(client)
	reservationRepo := postgresl.NewReversationRepository(db)
	reservationUsecase := usecase.NewReservationUsecase(reservationRepo, *lockRepo)
	reservationHandler := handler.NewReservationHandler(reservationUsecase)

	movieRepo := postgresl.NewMovieRepository(db)
	cache := redis.NewCache(client)
	movieUsecase := usecase.NewMovieUsecase(movieRepo, cache)
	movieHandler := handler.NewMovieHandler(movieUsecase)

	// health
	r.GET("/health", func(c *gin.Context) {
		data := map[string]interface{}{
			"message": "service is running",
			"requested_at": time.Now().
				Format("2006-01-02 15:04:05"),
		}
		response.Success(c, data)
	})

	// routes
	api := r.Group("/api/v1")
	{
		// reservations
		api.POST("/reservation", reservationHandler.CreateReservation)

		// movies
		api.GET("/movies", movieHandler.FindAllMovies)
	}

	return r
}
