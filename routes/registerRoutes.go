package routes

import (
	"database/sql"
	"net/http"
	"splitwise-app/controllers"

	"github.com/gorilla/mux"
)

func RegisterRoutes(router *mux.Router, db *sql.DB) {
	router.HandleFunc("/user", controllers.CreateUserHandler(db)).Methods(http.MethodPost)
	router.HandleFunc("/users/{id}", controllers.GetUserByIDHandler(db)).Methods(http.MethodGet)
	router.HandleFunc("/groups", controllers.CreateGroupHandler(db)).Methods(http.MethodPost)
	router.HandleFunc("/groups/{groupID}", controllers.GetGroupByIdHandler(db)).Methods(http.MethodGet)
	router.HandleFunc("/groups/{groupID}/members", controllers.AddGroupMemberHandler(db)).Methods(http.MethodPost)
	router.HandleFunc("/groups/{groupID}/members", controllers.GetGroupMembersHandler(db)).Methods(http.MethodGet)
}
