package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"effective-mobile/internal/handler"
	"effective-mobile/internal/model"
	"effective-mobile/pkg/logger"
)

func init() {
	logger.Init()
}

type mockPersonService struct{}

func (m *mockPersonService) CreatePerson(req model.CreatePersonRequest) (*model.Person, error) {
	return &model.Person{
		ID:          1,
		Name:        req.Name,
		Surname:     req.Surname,
		Gender:      "female",
		Age:         30,
		Nationality: "US",
	}, nil
}

func (m *mockPersonService) GetAllPersons() ([]model.Person, error) {
	return []model.Person{
		{ID: 1, Name: "Alice"},
	}, nil
}

func (m *mockPersonService) GetPersonByID(id uint) (*model.Person, error) {
	return &model.Person{ID: id, Name: "Alice"}, nil
}

func (m *mockPersonService) UpdatePerson(id uint, req model.UpdatePersonRequest) (*model.Person, error) {
	return &model.Person{ID: id, Name: req.Name}, nil
}

func (m *mockPersonService) DeletePerson(id uint) error {
	return nil
}

func TestCreatePersonHandler(t *testing.T) {
	svc := &mockPersonService{}
	h := handler.NewPersonHandler(svc)

	person := model.CreatePersonRequest{
		Name:    "Alice",
		Surname: "Smith",
	}

	body, _ := json.Marshal(person)

	req := httptest.NewRequest(http.MethodPost, "/person", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	h.CreatePerson(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("expected status 200 OK, got %d", res.StatusCode)
	}
}

func TestGetAllPersonsHandler(t *testing.T) {
	svc := &mockPersonService{}
	h := handler.NewPersonHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/person", nil)
	rec := httptest.NewRecorder()

	h.GetAllPersons(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200 OK, got %d", res.StatusCode)
	}
}

func TestGetPersonByIDHandler(t *testing.T) {
	svc := &mockPersonService{}
	h := handler.NewPersonHandler(svc)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /person/{id}", h.GetPersonByID)

	req := httptest.NewRequest(http.MethodGet, "/person/1", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Result().StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rec.Result().StatusCode)
	}
}

func TestUpdatePersonHandler(t *testing.T) {
	svc := &mockPersonService{}
	h := handler.NewPersonHandler(svc)

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /person/{id}", h.UpdatePerson)

	update := model.UpdatePersonRequest{Name: "UpdatedName"}
	body, _ := json.Marshal(update)

	req := httptest.NewRequest(http.MethodPut, "/person/1", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Result().StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rec.Result().StatusCode)
	}
}

func TestDeletePersonHandler(t *testing.T) {
	svc := &mockPersonService{}
	h := handler.NewPersonHandler(svc)

	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /person/{id}", h.DeletePerson)

	req := httptest.NewRequest(http.MethodDelete, "/person/1", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Result().StatusCode != http.StatusNoContent {
		t.Fatalf("expected 204 No Content, got %d", rec.Result().StatusCode)
	}
}
