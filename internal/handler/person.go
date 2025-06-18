package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"effective-mobile/internal/model"
	"effective-mobile/internal/service"
	"effective-mobile/pkg/logger"
	"effective-mobile/pkg/validator"

	"go.uber.org/zap"
)

type PersonHandler struct {
	service service.PersonServiceInterface
}

func NewPersonHandler(s service.PersonServiceInterface) *PersonHandler {
	return &PersonHandler{service: s}
}

// CreatePerson godoc
// @Summary Создание человека
// @Description Создаёт нового человека и обогащает его данными через внешние API
// @Tags persons
// @Accept json
// @Produce json
// @Param person body model.CreatePersonRequest true "Данные человека"
// @Success 201 {object} model.Person
// @Failure 400 {string} string "invalid JSON"
// @Failure 500 {string} string "failed to create person"
// @Router /person [post]
func (h *PersonHandler) CreatePerson(w http.ResponseWriter, r *http.Request) {
	logger.Log.Debug("POST /person - received request")

	var req model.CreatePersonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.Warn("failed to decode JSON", zap.Error(err))
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		http.Error(w, "validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	logger.Log.Info("creating person", zap.String("name", req.Name), zap.String("surname", req.Surname))
	person, err := h.service.CreatePerson(req)
	if err != nil {
		logger.Log.Error("failed to create person", zap.Error(err))
		http.Error(w, "failed to create person: "+err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Log.Info("person created", zap.Uint("id", person.ID))
	writeJSON(w, person, http.StatusCreated)
}

// GetAllPersons godoc
// @Summary Получение всех людей
// @Description Возвращает список всех сохранённых людей
// @Tags persons
// @Produce json
// @Success 200 {array} model.Person
// @Failure 500 {string} string "failed to get persons"
// @Router /person [get]
func (h *PersonHandler) GetAllPersons(w http.ResponseWriter, r *http.Request) {
	logger.Log.Debug("GET /person - listing all persons")

	people, err := h.service.GetAllPersons()
	if err != nil {
		logger.Log.Error("failed to get persons", zap.Error(err))
		http.Error(w, "failed to get persons: "+err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Log.Info("persons fetched", zap.Int("count", len(people)))
	writeJSON(w, people, http.StatusOK)
}

// GetPersonByID godoc
// @Summary Получение человека по ID
// @Description Возвращает данные конкретного человека
// @Tags persons
// @Produce json
// @Param id path int true "ID человека"
// @Success 200 {object} model.Person
// @Failure 400 {string} string "invalid ID"
// @Failure 404 {string} string "person not found"
// @Router /person/{id} [get]
func (h *PersonHandler) GetPersonByID(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := strconv.ParseUint(idString, 10, 32)
	if err != nil {
		logger.Log.Warn("invalid ID", zap.String("id", idString), zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	logger.Log.Debug("fetching person", zap.Uint("id", uint(id)))
	person, err := h.service.GetPersonByID(uint(id))
	if err != nil {
		logger.Log.Warn("person not found", zap.Uint("id", uint(id)))
		http.Error(w, "person not found", http.StatusNotFound)
		return
	}

	logger.Log.Info("person fetched", zap.Uint("id", uint(id)))
	writeJSON(w, person, http.StatusOK)
}

// UpdatePerson godoc
// @Summary Обновление человека
// @Description Обновляет существующего человека по ID
// @Tags persons
// @Accept json
// @Produce json
// @Param id path int true "ID человека"
// @Param person body model.UpdatePersonRequest true "Обновлённые данные"
// @Success 200 {object} model.Person
// @Failure 400 {string} string "invalid ID or JSON"
// @Failure 500 {string} string "failed to update"
// @Router /person/{id} [put]
func (h *PersonHandler) UpdatePerson(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := strconv.ParseUint(idString, 10, 32)
	if err != nil {
		logger.Log.Warn("invalid ID", zap.String("id", idString), zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var req model.UpdatePersonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.Warn("failed to decode update JSON", zap.Error(err))
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	logger.Log.Info("updating person", zap.Uint("id", uint(id)))
	person, err := h.service.UpdatePerson(uint(id), req)
	if err != nil {
		logger.Log.Error("failed to update person", zap.Error(err))
		http.Error(w, "failed to update: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, person, http.StatusOK)
}

// DeletePerson godoc
// @Summary Удаление человека
// @Description Удаляет человека по ID
// @Tags persons
// @Param id path int true "ID человека"
// @Success 204 {string} string "no content"
// @Failure 400 {string} string "invalid ID"
// @Failure 500 {string} string "failed to delete"
// @Router /person/{id} [delete]
func (h *PersonHandler) DeletePerson(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := strconv.ParseUint(idString, 10, 32)
	if err != nil {
		logger.Log.Warn("invalid ID", zap.String("id", idString), zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	logger.Log.Info("deleting person", zap.Uint("id", uint(id)))
	err = h.service.DeletePerson(uint(id))
	if err != nil {
		logger.Log.Error("failed to delete person", zap.Error(err))
		http.Error(w, "failed to delete: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	logger.Log.Info("person deleted", zap.Uint("id", uint(id)))
}

func writeJSON(w http.ResponseWriter, data any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
