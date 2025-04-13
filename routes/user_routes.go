package routes

import (
	"database/sql"
	"net/http"
	"splitwise-app/controllers"

	"github.com/gorilla/mux"
)

func RegisterUserRoutes(router *mux.Router, db *sql.DB) {
	router.HandleFunc("/users/{id}", controllers.GetUserByIDHandler(db)).Methods(http.MethodGet)
	router.HandleFunc("/user", controllers.CreateUserHandler(db)).Methods(http.MethodPost)
}
