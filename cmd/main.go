package main

import (
	"fmt"
	"log"
	"net/http"

	"effective-mobile/database"
	"effective-mobile/internal/handler"
	"effective-mobile/internal/model"
	"effective-mobile/internal/repository"
	"effective-mobile/internal/service"
	"effective-mobile/pkg/logger"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title           People Info API
// @version         1.0
// @description     REST API для работы с информацией о людях (создание, чтение, обновление, удаление)
// @host            localhost:8080
// @BasePath        /
func main() {
	logger.Init()
	defer logger.Log.Sync()

	db := database.NewDB()
	db.AutoMigrate(&model.Person{})

	repo := repository.NewPersonRepository(db)
	svc := service.NewPersonService(repo)
	personHandler := handler.NewPersonHandler(svc)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /person", personHandler.CreatePerson)
	mux.HandleFunc("GET /person", personHandler.GetAllPersons)
	mux.HandleFunc("GET /person/{id}", personHandler.GetPersonByID)
	mux.HandleFunc("PUT /person/{id}", personHandler.UpdatePerson)
	mux.HandleFunc("DELETE /person/{id}", personHandler.DeletePerson)

	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	fmt.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
