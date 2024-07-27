package routers

import (
	"net/http"
	"note-service/handlers"

	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	router := mux.NewRouter()

	// API routes
	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/notes", handlers.GetNotes).Methods("GET")
	apiRouter.HandleFunc("/notes", handlers.CreateNote).Methods("POST")
	apiRouter.HandleFunc("/notes/{id}", handlers.GetNote).Methods("GET")
	apiRouter.HandleFunc("/notes/{id}", handlers.UpdateNote).Methods("PUT")
	apiRouter.HandleFunc("/notes/{id}", handlers.DeleteNote).Methods("DELETE")

	// Serve static files from the frontend directory
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./frontend/"))))

	return router
}
