package service

import (
	"effective-mobile/internal/model"
	"effective-mobile/internal/repository"
)

type PersonServiceInterface interface {
	CreatePerson(req model.CreatePersonRequest) (*model.Person, error)
	GetAllPersons() ([]model.Person, error)
	GetPersonByID(id uint) (*model.Person, error)
	UpdatePerson(id uint, req model.UpdatePersonRequest) (*model.Person, error)
	DeletePerson(id uint) error
}

type PersonService struct {
	repo repository.PersonRepositoryInterface
}

func NewPersonService(repo repository.PersonRepositoryInterface) *PersonService {
	return &PersonService{repo: repo}
}

func (s *PersonService) CreatePerson(req model.CreatePersonRequest) (*model.Person, error) {
	data, err := EnrichPerson(req.Name)
	if err != nil {
		return nil, err
	}

	person := &model.Person{
		Name:        req.Name,
		Surname:     req.Surname,
		Patronymic:  req.Patronymic,
		Gender:      data.Gender,
		Age:         data.Age,
		Nationality: data.Nationality,
	}

	if err := s.repo.Save(person); err != nil {
		return nil, err
	}

	return person, nil
}

func (s *PersonService) GetAllPersons() ([]model.Person, error) {
	return s.repo.FindAll()
}

func (s *PersonService) GetPersonByID(id uint) (*model.Person, error) {
	return s.repo.FindByID(id)
}

func (s *PersonService) UpdatePerson(id uint, update model.UpdatePersonRequest) (*model.Person, error) {
	p, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if update.Name != "" {
		p.Name = update.Name
	}
	if update.Surname != "" {
		p.Surname = update.Surname
	}
	if update.Patronymic != "" {
		p.Patronymic = update.Patronymic
	}
	if update.Gender != "" {
		p.Gender = update.Gender
	}
	if update.Age != 0 {
		p.Age = update.Age
	}
	if update.Nationality != "" {
		p.Nationality = update.Nationality
	}

	return s.repo.Update(p)
}

func (s *PersonService) DeletePerson(id uint) error {
	return s.repo.Delete(id)
}
