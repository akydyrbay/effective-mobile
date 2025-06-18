package model

type CreatePersonRequest struct {
	Name       string `json:"name" validate:"required"`
	Surname    string `json:"surname" validate:"required"`
	Patronymic string `json:"patronymic"`
}

type UpdatePersonRequest struct {
	Name        string `json:"name,omitempty"`
	Surname     string `json:"surname,omitempty"`
	Patronymic  string `json:"patronymic,omitempty"`
	Gender      string `json:"gender,omitempty"`
	Age         int    `json:"age,omitempty"`
	Nationality string `json:"nationality,omitempty"`
}

// type PersonResponse struct {
// 	ID          uint   `json:"id"`
// 	Name        string `json:"name"`
// 	Surname     string `json:"surname"`
// 	Patronymic  string `json:"patronymic,omitempty"`
// 	Gender      string `json:"gender,omitempty"`
// 	Age         int    `json:"age,omitempty"`
// 	Nationality string `json:"nationality,omitempty"`
// 	CreatedAt   string `json:"created_at"`
// 	UpdatedAt   string `json:"updated_at"`
// }
