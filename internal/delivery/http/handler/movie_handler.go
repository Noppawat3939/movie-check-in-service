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

type MoviewHandler struct {
	movieUsecase usecase.MovieUsecase
}

func NewMovieHandler(movieUsecase usecase.MovieUsecase) *MoviewHandler {
	return &MoviewHandler{movieUsecase}
}

func (h *MoviewHandler) FindAllMovies(c *gin.Context) {
	data, err := h.movieUsecase.FindAllMovies(c)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}

	response.Success(c, data)
}

func (h *MoviewHandler) FindByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	data, err := h.movieUsecase.FindMovieByID(c, id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, domain.ErrMovieNotFound) {
			statusCode = http.StatusNotFound
		}
		response.Error(c, statusCode, err.Error(), id)
		return
	}

	response.Success(c, data)
}
