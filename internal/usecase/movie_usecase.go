package usecase

import (
	"check-in/internal/domain"
	"check-in/internal/infra/postgresl"
	"check-in/internal/infra/redis"
	"context"
	"time"
)

type MovieUsecase interface {
	FindAllMovies(ctx context.Context) ([]domain.Movie, error)
}

type movieUsecase struct {
	movieRepo postgresl.MovieRepository
	cache     *redis.Cache
}

const movieCacheKey = "movie:all"

func NewMovieUsecase(movieRepo postgresl.MovieRepository, cache *redis.Cache) MovieUsecase {
	return &movieUsecase{movieRepo, cache}
}

func (u *movieUsecase) FindAllMovies(ctx context.Context) ([]domain.Movie, error) {
	// get cached
	var data []domain.Movie
	if err := u.cache.Get(ctx, movieCacheKey, &data); err == nil {
		return data, nil
	}

	movies, err := u.movieRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	// update cache
	u.cache.Set(ctx, movieCacheKey, movies, 5*time.Minute)

	return movies, nil
}
