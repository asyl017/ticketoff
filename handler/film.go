package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"ticketoff/repositories"
	"ticketoff/utils"
)

type FilmHandler interface {
	CreateFilm(w http.ResponseWriter, r *http.Request)
	GetFilmByID(w http.ResponseWriter, r *http.Request)
	UpdateFilm(w http.ResponseWriter, r *http.Request)
	DeleteFilm(w http.ResponseWriter, r *http.Request)
	GetFilms(w http.ResponseWriter, r *http.Request)
}

type filmHandler struct {
	filmRepo repositories.FilmRepository
}

func NewFilmHandler(filmRepo repositories.FilmRepository) *filmHandler {
	return &filmHandler{filmRepo: filmRepo}
}

/*
func (f *filmHandler) CreateFilm(w http.ResponseWriter, r *http.Request) {
	var film models.Film
	if err := json.NewDecoder(r.Body).Decode(&film); err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error": err,
		}).Error("Invalid request payload")
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if err := f.filmRepo.CreateFilm(&film); err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error": err,
		}).Error("Error creating film")
		http.Error(w, "Error creating film: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(film)
}
*/

func (f *filmHandler) GetFilmByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	film, err := f.filmRepo.GetFilmByID(id)
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error": err,
			"id":    id,
		}).Error("Film not found")
		http.Error(w, "Film not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(film)
}

/*
func (f *filmHandler) UpdateFilm(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var updatedFilm models.Film
	if err := json.NewDecoder(r.Body).Decode(&updatedFilm); err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error": err,
		}).Error("Invalid request payload")
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	parsedID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error": err,
			"id":    id,
		}).Error("Invalid ID format")
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}
	updatedFilm.ID = uint(parsedID)
	if _, err := f.filmRepo.UpdateFilm(&updatedFilm); err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error": err,
			"id":    id,
		}).Error("Error updating film")
		http.Error(w, "Error updating film: "+err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(updatedFilm)
}

func (f *filmHandler) DeleteFilm(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if err := f.filmRepo.DeleteFilm(id); err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error": err,
			"id":    id,
		}).Error("Error deleting film")
		http.Error(w, "Error deleting film: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
*/

func (f *filmHandler) GetFilms(w http.ResponseWriter, r *http.Request) {

	filter := r.URL.Query().Get("filter")
	sort := r.URL.Query().Get("sort")
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	// TODO: I changed it to 10. Change it to 17 later!
	limit := 10 // 17 films per page to divide 50 films into 3 pages
	offset := (page - 1) * limit

	fmt.Println(offset, " ", page)

	films, err := f.filmRepo.GetFilms(filter, sort, limit, offset)
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error": err,
		}).Error("Error fetching films")
		http.Error(w, "Error fetching films: "+err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(films)
}
