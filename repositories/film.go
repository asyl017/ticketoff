package repositories

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"ticketoff/models"
)

type FilmRepository interface {
	CreateFilm(film *models.Film) error
	GetFilmByID(id string) (*models.Film, error)
	GetFilms(filter string, sort string, limit int, offset int) ([]models.Film, error)
	UpdateFilm(updatedFilm *models.Film) (*models.Film, error)
	DeleteFilm(id string) error
}

type filmRepository struct {
	db *mongo.Database
}

func NewFilmRepository(db *mongo.Database) FilmRepository {
	return &filmRepository{db: db}
}

func (f filmRepository) CreateFilm(film *models.Film) error {
	collection := f.db.Collection("movies")
	_, err := collection.InsertOne(context.Background(), film)
	return err
}

func (f filmRepository) GetFilmByID(id string) (*models.Film, error) {
	collection := f.db.Collection("movies")
	var film models.Film
	err := collection.FindOne(context.Background(), bson.M{"id": id}).Decode(&film)
	return &film, err
}

func (f filmRepository) GetFilms(filter string, sort string, limit int, offset int) ([]models.Film, error) {
	collection := f.db.Collection("movies")

	// Create filter and sort options
	filterOptions := bson.M{}
	if filter != "" {
		filterOptions = bson.M{"title": bson.M{"$regex": filter, "$options": "i"}}
	}

	sortOptions := bson.D{}
	if sort != "" {
		sortOptions = bson.D{{Key: sort, Value: 1}}
	}

	findOptions := options.Find()
	findOptions.SetSort(sortOptions)
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(offset))

	cursor, err := collection.Find(context.Background(), filterOptions, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var films []models.Film
	for cursor.Next(context.Background()) {
		var film models.Film
		if err := cursor.Decode(&film); err != nil {
			return nil, err
		}
		films = append(films, film)
	}
	return films, nil
}

func (f filmRepository) UpdateFilm(updatedFilm *models.Film) (*models.Film, error) {
	collection := f.db.Collection("movies")
	_, err := collection.UpdateOne(context.Background(), bson.M{"id": updatedFilm.ID}, bson.M{"$set": updatedFilm})
	return updatedFilm, err
}

func (f filmRepository) DeleteFilm(id string) error {
	collection := f.db.Collection("movies")
	_, err := collection.DeleteOne(context.Background(), bson.M{"id": id})
	return err
}
