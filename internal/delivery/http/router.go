package http

import (
	"check-in/internal/delivery/http/handler"
	"check-in/internal/infra/postgresl"
	"check-in/internal/usecase"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	reservationRepo := postgresl.NewReversationRepository(db)
	reservationUsecase := usecase.NewReverationUsecase(reservationRepo)
	reservationHandler := handler.NewReservationHandler(reservationUsecase)

	// health
	r.GET("/health", func(c *gin.Context) {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{"message": "service is running", "requested_at": time.Now().Format("2006-01-02 15:04:05")})
	})

	// routes
	api := r.Group("/api/v1")
	{
		// reservations
		api.POST("/reservation", reservationHandler.CreateReservation)

		// movies
	}

	return r
}
