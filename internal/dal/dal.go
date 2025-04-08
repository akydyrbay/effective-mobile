package dal

import (
	"database/sql"
	"effective-mobile/models"
	"fmt"
	"strings"
)

type PersonRepository interface {
	GetAll(filter models.FilterParams, pagination models.PaginationParams) ([]models.Person, error)
	Exists(input models.PersonInput) (bool, error)
	Insert(person models.Person) error
	Update(person models.Person) error
	Delete(id int) error
}

type personRepo struct {
	DB *sql.DB
}

func NewPersonRepository(db *sql.DB) PersonRepository {
	return &personRepo{DB: db}
}

func (r *personRepo) GetAll(filter models.FilterParams, pagination models.PaginationParams) ([]models.Person, error) {
	query := "SELECT id, name, surname, patronymic, age, gender, nationality FROM person"
	var conditions []string
	var args []interface{}
	argID := 1

	if filter.Name != "" {
		conditions = append(conditions, fmt.Sprintf("name ILIKE $%d", argID))
		args = append(args, "%"+filter.Name+"%")
		argID++
	}
	if filter.Surname != "" {
		conditions = append(conditions, fmt.Sprintf("surname ILIKE $%d", argID))
		args = append(args, "%"+filter.Surname+"%")
		argID++
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	limit := pagination.Limit
	if limit <= 0 {
		limit = 10 // default limit
	}
	offset := (pagination.Page - 1) * limit
	if offset < 0 {
		offset = 0
	}
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argID, argID+1)
	args = append(args, limit, offset)

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var persons []models.Person
	for rows.Next() {
		var p models.Person
		err := rows.Scan(&p.ID, &p.Name, &p.Surname, &p.Patronymic, &p.Age, &p.Gender, &p.Nationality)
		if err != nil {
			return nil, err
		}
		persons = append(persons, p)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return persons, nil
}

func (r *personRepo) Exists(input models.PersonInput) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM person 
			WHERE name = $1 AND surname = $2 AND (patronymic = $3 OR ($3 = '' AND patronymic IS NULL))
		)
	`
	var patronymic string
	if input.Patronymic != "" {
		patronymic = input.Patronymic
	}
	var exists bool
	err := r.DB.QueryRow(query, input.Name, input.Surname, patronymic).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *personRepo) Insert(person models.Person) error {
	query := `
		INSERT INTO person (name, surname, patronymic, age, gender, nationality)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id;
	`
	err := r.DB.QueryRow(query, person.Name, person.Surname, person.Patronymic,
		person.Age, person.Gender, person.Nationality).Scan(&person.ID)
	return err
}

func (r *personRepo) Update(person models.Person) error {
	query := `
		UPDATE person
		SET name = $1, surname = $2, patronymic = $3, age = $4, gender = $5, nationality = $6
		WHERE id = $7
	`
	res, err := r.DB.Exec(query, person.Name, person.Surname, person.Patronymic,
		person.Age, person.Gender, person.Nationality, person.ID)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no rows updated with id %d", person.ID)
	}
	return nil
}

func (r *personRepo) Delete(id int) error {
	query := `DELETE FROM person WHERE id = $1`
	res, err := r.DB.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no rows deleted with id %d", id)
	}
	return nil
}
