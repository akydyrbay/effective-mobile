package repository

import (
	"errors"

	"effective-mobile/database"
	"effective-mobile/internal/model"
)

type PersonRepositoryInterface interface {
	Save(p *model.Person) error
	FindAll() ([]model.Person, error)
	FindByID(id uint) (*model.Person, error)
	Update(p *model.Person) (*model.Person, error)
	Delete(id uint) error
}

type PersonRepository struct {
	db *database.DB
}

func NewPersonRepository(db *database.DB) *PersonRepository {
	return &PersonRepository{db: db}
}

func (r *PersonRepository) Save(p *model.Person) error {
	return r.db.Create(p).Error
}

func (r *PersonRepository) FindAll() ([]model.Person, error) {
	var people []model.Person
	err := r.db.Find(&people).Error
	return people, err
}

func (r *PersonRepository) FindByID(id uint) (*model.Person, error) {
	var p model.Person
	if err := r.db.First(&p, id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PersonRepository) Update(p *model.Person) (*model.Person, error) {
	err := r.db.Save(p).Error
	return p, err
}

func (r *PersonRepository) Delete(id uint) error {
	res := r.db.Delete(&model.Person{}, id)
	if res.RowsAffected == 0 {
		return errors.New("not found")
	}
	return res.Error
}
