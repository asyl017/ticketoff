package repositories

import (
	"github.com/jinzhu/gorm"
	"ticketoff/models"
)

type FilmRepository interface {
	CreateFilm(film *models.Film) error
	GetFilmByID(id string) (*models.Film, error)
	UpdateFilm(updatedFilm *models.Film) (*models.Film, error)
	DeleteFilm(id string) error
	GetFilms(filter, sort string, limit, offset int) ([]models.Film, error)
}

type filmRepository struct {
	db *gorm.DB
}

func NewFilmRepository(db *gorm.DB) FilmRepository {
	return &filmRepository{db: db}
}

func (f filmRepository) CreateFilm(film *models.Film) error {
	return f.db.Create(film).Error
}

func (f filmRepository) GetFilmByID(id string) (*models.Film, error) {
	var film models.Film
	if err := f.db.First(&film, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &film, nil
}

func (f filmRepository) UpdateFilm(updatedFilm *models.Film) (*models.Film, error) {
	err := f.db.Save(updatedFilm).Error
	return updatedFilm, err
}

func (f filmRepository) DeleteFilm(id string) error {
	return f.db.Delete(&models.Film{}, "id = ?", id).Error
}

func (r *filmRepository) GetFilms(filter, sort string, limit, offset int) ([]models.Film, error) {
	var films []models.Film
	query := r.db.Limit(limit).Offset(offset)
	if filter != "" {
		query = query.Where("title ILIKE ?", "%"+filter+"%")
	}
	if sort != "" {
		query = query.Order(sort)
	}
	if err := query.Find(&films).Error; err != nil {
		return nil, err
	}
	return films, nil
}
