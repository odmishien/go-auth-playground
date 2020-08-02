package main

import (
	"net/http"

	"github.com/odmishien/go-auth-playground/auth"
	"github.com/odmishien/go-auth-playground/database"
	"github.com/odmishien/go-auth-playground/handlers"
)

func main() {
	// setting for database
	var databaseConfig = database.DatabaseConfig{
		User:     "root",
		Password: "",
		Host:     "localhost",
		Port:     "3306",
		Database: "go_auth_playground",
		Debug:    true,
	}

	// init database
	db, err := database.InitDatabase(databaseConfig)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	userHandler := handlers.UserHandler{
		Db: db,
	}

	// routing & serve
	http.Handle("/signup", http.HandlerFunc(userHandler.PreCreateUser))
	http.Handle("/verify", auth.JwtMiddleware.Handler(http.HandlerFunc(userHandler.CreateUser)))

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err.Error())
	}
}
