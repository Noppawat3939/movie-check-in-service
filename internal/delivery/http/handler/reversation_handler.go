package handler

import (
	"check-in/internal/delivery/http/response"
	"check-in/internal/domain"
	"check-in/internal/usecase"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ReservationHandler struct {
	reservationUsecase usecase.ReservationUsecase
}

func NewReservationHandler(usecase usecase.ReservationUsecase) *ReservationHandler {
	return &ReservationHandler{usecase}
}

func (h *ReservationHandler) CreateReservation(c *gin.Context) {
	var req domain.CreateReservationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, domain.ErrBodyInvalidMsg)
		return
	}

	data, err := h.reservationUsecase.CreateReservation(c, req)

	if err != nil {
		statusCode := http.StatusInternalServerError

		if errors.Is(err, domain.ErrSeatAlreadyReserved) {
			statusCode = http.StatusConflict
		}

		response.Error(c, statusCode, err.Error(), req)
		return
	}

	response.Success(c, data)
}

func (h *ReservationHandler) ListReservation(c *gin.Context) {
	showtimeID, err := uuid.Parse(c.Param("showtimeID"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id format")
		return
	}

	data, err := h.reservationUsecase.ListReservation(c, showtimeID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, data)
}
