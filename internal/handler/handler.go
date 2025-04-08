package handler

import (
	"effective-mobile/internal/service"
	"effective-mobile/models"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type PersonHandler interface {
	PostPerson(w http.ResponseWriter, r *http.Request)
	GetPersons(w http.ResponseWriter, r *http.Request)
	PutPerson(w http.ResponseWriter, r *http.Request)
	DeletePerson(w http.ResponseWriter, r *http.Request)
}

type personHandler struct {
	Service service.PersonService
}

func NewPersonHandler(personService service.PersonService) *personHandler {
	return &personHandler{Service: personService}
}

func (h *personHandler) PostPerson(w http.ResponseWriter, r *http.Request) {
	var input models.PersonInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	err := h.Service.CreatePerson(input)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to create person: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "Person created successfully")
}

func (h *personHandler) PutPerson(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var input models.PersonInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	if err := h.Service.UpdatePerson(id, input); err != nil {
		http.Error(w, fmt.Sprintf("update failed: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Person updated successfully")
}

func (h *personHandler) DeletePerson(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	err = h.Service.DeletePerson(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("delete failed: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Person deleted successfully")
}

func (h *personHandler) GetPersons(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	surname := r.URL.Query().Get("surname")

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, _ := strconv.Atoi(pageStr)
	if page <= 0 {
		page = 1
	}
	limit, _ := strconv.Atoi(limitStr)
	if limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit

	filter := models.FilterParams{
		Name:    name,
		Surname: surname,
	}
	pagination := models.PaginationParams{
		Limit: limit,
		Page:  offset,
	}

	data, err := h.Service.GetPersons(filter, pagination)
	if err != nil {
		http.Error(w, fmt.Sprintf("fetch failed: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
