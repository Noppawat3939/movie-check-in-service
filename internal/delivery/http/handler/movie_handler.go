package handler

import (
	"check-in/internal/delivery/http/response"
	"check-in/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
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
