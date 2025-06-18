package service_test

import (
	"errors"
	"testing"

	"effective-mobile/internal/model"
	"effective-mobile/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ---- MOCK REPO ----

type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) Save(p *model.Person) error {
	args := m.Called(p)
	return args.Error(0)
}

func (m *mockRepo) FindAll() ([]model.Person, error) {
	args := m.Called()
	return args.Get(0).([]model.Person), args.Error(1)
}

func (m *mockRepo) FindByID(id uint) (*model.Person, error) {
	args := m.Called(id)
	return args.Get(0).(*model.Person), args.Error(1)
}

func (m *mockRepo) Update(p *model.Person) (*model.Person, error) {
	args := m.Called(p)
	return args.Get(0).(*model.Person), args.Error(1)
}

func (m *mockRepo) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

// ---- TESTS ----

func TestCreatePerson(t *testing.T) {
	mockRepo := new(mockRepo)
	svc := service.NewPersonService(mockRepo)

	req := model.CreatePersonRequest{
		Name:    "Alice",
		Surname: "Smith",
	}

	mockRepo.On("Save", mock.AnythingOfType("*model.Person")).Return(nil)

	result, err := svc.CreatePerson(req)

	assert.NoError(t, err)
	assert.Equal(t, "Alice", result.Name)
	assert.Equal(t, "Smith", result.Surname)

	mockRepo.AssertExpectations(t)
}

func TestUpdatePerson_OnlyUpdatesProvidedFields(t *testing.T) {
	mockRepo := new(mockRepo)
	svc := service.NewPersonService(mockRepo)

	existing := &model.Person{
		ID:      1,
		Name:    "OldName",
		Surname: "OldSurname",
		Age:     25,
	}

	mockRepo.On("FindByID", uint(1)).Return(existing, nil)

	mockRepo.On("Update", mock.MatchedBy(func(p *model.Person) bool {
		return p.Name == "NewName" && p.Surname == "OldSurname" && p.Age == 30
	})).Return(existing, nil)

	req := model.UpdatePersonRequest{
		Name: "NewName",
		Age:  30,
	}

	updated, err := svc.UpdatePerson(1, req)

	assert.NoError(t, err)
	assert.Equal(t, "NewName", updated.Name)
	assert.Equal(t, 30, updated.Age)

	mockRepo.AssertExpectations(t)
}

func TestGetPersonByID_NotFound(t *testing.T) {
	mockRepo := new(mockRepo)
	svc := service.NewPersonService(mockRepo)

	mockRepo.On("FindByID", uint(100)).Return(&model.Person{}, errors.New("not found"))

	_, err := svc.GetPersonByID(100)

	assert.Error(t, err)
}
