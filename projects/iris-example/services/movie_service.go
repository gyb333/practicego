package services

import (
	"iris-example/datamodels"
	"iris-example/repositories"
)

// MovieService处理电影数据模型的一些CRUID操作。
//因为我们可能需要在将来更改或尝试实验性的不同域逻辑。
type MovieService interface {
	GetAll() []datamodels.Movie
	GetByID(id int64) (datamodels.Movie, bool)
	DeleteByID(id int64) bool
	UpdatePosterAndGenreByID(id int64, poster string, genre string) (datamodels.Movie, error)
}
// NewMovieService返回默认 movie service.
func NewMovieService(repo repositories.MovieRepository) MovieService {
	return &movieService{
		repo: repo,
	}
}


type movieService struct {
	repo repositories.MovieRepository
}
// GetAll 获取所有的movie.
func (s *movieService) GetAll() []datamodels.Movie {
	return s.repo.SelectMany(func(_ datamodels.Movie) bool {
		return true
	}, -1)
}
// GetByID 根据其ID返回一行。
func (s *movieService) GetByID(id int64) (datamodels.Movie, bool) {
	return s.repo.Select(func(m datamodels.Movie) bool {
		return m.ID == id
	})
}
// UpdatePosterAndGenreByID更新电影的海报和流派。
func (s *movieService) UpdatePosterAndGenreByID(id int64, poster string, genre string) (datamodels.Movie, error) {
	// update the movie and return it.
	return s.repo.InsertOrUpdate(datamodels.Movie{
		ID:     id,
		Poster: poster,
		Genre:  genre,
	})
}
// DeleteByID按ID删除电影。
//
//如果删除则返回true，否则返回false。
func (s *movieService) DeleteByID(id int64) bool {
	return s.repo.Delete(func(m datamodels.Movie) bool {
		return m.ID == id
	}, 1)
}
