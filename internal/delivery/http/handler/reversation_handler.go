package handler

import (
	"check-in/internal/delivery/http/response"
	"check-in/internal/domain"
	"check-in/internal/usecase"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ReservationHandler struct {
	reservationUsecase usecase.ReverationUsecase
}

func NewReservationHandler(usecase usecase.ReverationUsecase) *ReservationHandler {
	return &ReservationHandler{usecase}
}

func (h *ReservationHandler) CreateReservation(c *gin.Context) {
	var req domain.CreateReservationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, domain.ErrBodyInvalidMsg)
		return
	}

	data, err := h.reservationUsecase.CreateReveration(c, req)

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
