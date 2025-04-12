package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"effective-mobile/internal/dal"
	"effective-mobile/models"
)

type PersonService interface {
	CreatePerson(input models.PersonInput) error
	UpdatePerson(id int, input models.PersonInput) error
	DeletePerson(id int) error
	GetPersons(filter models.FilterParams, pagination models.PaginationParams) ([]models.Person, error)
	EnrichPerson(name string) (age int, gender string, nationality string, err error)
}

type personService struct {
	repo   dal.PersonRepository
	client *http.Client
}

func NewPersonService(repo dal.PersonRepository) PersonService {
	return &personService{
		repo: repo,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (s *personService) EnrichPerson(name string) (int, string, string, error) {
	age, err := s.getAge(name)
	if err != nil {
		return 0, "", "", err
	}
	gender, err := s.getGender(name)
	if err != nil {
		return 0, "", "", err
	}
	nationality, err := s.getNationality(name)
	if err != nil {
		return 0, "", "", err
	}
	return age, gender, nationality, nil
}

func (s *personService) getAge(name string) (int, error) {
	url := fmt.Sprintf("https://api.agify.io/?name=%s", name)
	resp, err := s.client.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var result struct {
		Age int `json:"age"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	return result.Age, err
}

func (s *personService) getGender(name string) (string, error) {
	url := fmt.Sprintf("https://api.genderize.io/?name=%s", name)
	resp, err := s.client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Gender string `json:"gender"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	return result.Gender, err
}

func (s *personService) getNationality(name string) (string, error) {
	url := fmt.Sprintf("https://api.nationalize.io/?name=%s", name)
	resp, err := s.client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Country []struct {
			CountryID   string  `json:"country_id"`
			Probability float64 `json:"probability"`
		} `json:"country"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}
	if len(result.Country) > 0 {
		return result.Country[0].CountryID, nil
	}
	return "", nil
}

func (s *personService) CreatePerson(input models.PersonInput) error {
	exists, err := s.repo.Exists(input)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("person already exists")
	}

	age, gender, nationality, err := s.EnrichPerson(input.Name)
	if err != nil {
		return err
	}

	person := models.Person{
		Name:        input.Name,
		Surname:     input.Surname,
		Patronymic:  input.Patronymic,
		Age:         age,
		Gender:      gender,
		Nationality: nationality,
	}

	return s.repo.Insert(person)
}

func (s *personService) UpdatePerson(id int, input models.PersonInput) error {
	age, gender, nationality, err := s.EnrichPerson(input.Name)
	if err != nil {
		return err
	}

	person := models.Person{
		ID:          id,
		Name:        input.Name,
		Surname:     input.Surname,
		Patronymic:  input.Patronymic,
		Age:         age,
		Gender:      gender,
		Nationality: nationality,
	}

	return s.repo.Update(person)
}

func (s *personService) DeletePerson(id int) error {
	return s.repo.Delete(id)
}

func (s *personService) GetPersons(filter models.FilterParams, pagination models.PaginationParams) ([]models.Person, error) {
	return s.repo.GetAll(filter, pagination)
}
