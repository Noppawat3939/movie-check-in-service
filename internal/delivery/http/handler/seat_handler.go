package handler

import (
	"check-in/internal/delivery/http/response"
	"check-in/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SeatHandler struct {
	seatUsecase usecase.SeatUsecase
}

func NewSeatHandler(seatUsecase usecase.SeatUsecase) *SeatHandler {
	return &SeatHandler{seatUsecase}
}

func (h *SeatHandler) FindAllSeatAvailable(c *gin.Context) {
	showtimeID, err := uuid.Parse(c.Param("showtimeID"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	data, err := h.seatUsecase.GetSeatAvailable(c, showtimeID)
	if err != nil {
		response.Error(c, http.StatusConflict, err.Error())
		return
	}
	response.Success(c, data)
}
