package cmd

import (
	"bufio"
	"flag"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"

	"effective-mobile/internal/dal"
	"effective-mobile/internal/handler"
	"effective-mobile/internal/service"
	"effective-mobile/postgres"
)

func ServeSwaggerFile(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "swagger.yaml")
}

func loadEnv() error {
	file, err := os.Open("config.env")
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			log.Printf("[WARN] Ignoring malformed line: %s", line)
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		value = strings.Trim(value, `"'`)
		os.Setenv(key, value)
	}

	return scanner.Err()
}

func StartServer() {
	err := loadEnv()
	if err != nil {
		log.Printf("[ERROR] Could not load .env file: %v", err)
	} else {
		log.Println("[INFO] Environment variables loaded from .env file successfully.")
	}

	db, err := postgres.CheckDB()
	if err != nil {
		slog.Error("Failed to start program", "CheckDB err:", err)
		log.Fatal(err)
	}
	defer db.Close()
	port := flag.Int("port", 8080, "The server port")

	personRepo := dal.NewPersonRepository(db)
	personService := service.NewPersonService(personRepo)
	personHandler := handler.NewPersonHandler(personService)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /persons", personHandler.PostPerson)
	mux.HandleFunc("GET /persons", personHandler.GetPersons)
	mux.HandleFunc("PUT /persons", personHandler.PostPerson)
	mux.HandleFunc("DELETE /persons", personHandler.DeletePerson)
	mux.HandleFunc("/swagger.yaml", ServeSwaggerFile)

	slog.Info("Starting server on port", "port", *port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), mux))
}
