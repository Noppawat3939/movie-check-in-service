package usecase

import (
	"check-in/internal/domain"
	"check-in/internal/infra/postgresl"
	"check-in/internal/infra/redis"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type MovieUsecase interface {
	FindAllMovies(ctx context.Context) ([]domain.Movie, error)
	FindMovieByID(ctx context.Context, id uuid.UUID) (*domain.Movie, error)
}

type movieUsecase struct {
	movieRepo postgresl.MovieRepository
	cache     *redis.Cache
}

func NewMovieUsecase(movieRepo postgresl.MovieRepository, cache *redis.Cache) MovieUsecase {
	return &movieUsecase{movieRepo, cache}
}

func (u *movieUsecase) FindAllMovies(ctx context.Context) ([]domain.Movie, error) {
	// get cached
	const movieCacheKey = "movie:all"
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

func (u *movieUsecase) FindMovieByID(ctx context.Context, id uuid.UUID) (*domain.Movie, error) {
	// get cached
	var movie *domain.Movie
	cacheKey := fmt.Sprintf("movie:%s", id)
	if err := u.cache.Get(ctx, cacheKey, movie); err == nil {
		return movie, err
	}

	movie, err := u.movieRepo.FindMovieByID(ctx, id)
	if err != nil {
		return nil, domain.ErrMovieNotFound
	}

	// set cache
	_ = u.cache.Set(ctx, cacheKey, movie, 3*time.Hour)

	return movie, nil
}
