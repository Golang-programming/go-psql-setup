package router

import (
	"go-psql-setup/handler"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter().PathPrefix("/api").Subrouter()

	router.HandleFunc("/stock/{id}", handler.GetStock).Methods("GET", "OPTIONS")
	router.HandleFunc("/stocks/{id}", handler.UpdateStock).Methods("UPDATE", "OPTIONS")
	router.HandleFunc("/stocks/{id}", handler.DeleteStock).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/stocks", handler.GetAllStocks).Methods("GET", "OPTIONS")
	router.HandleFunc("/stock", handler.CreateStock).Methods("POST", "OPTIONS")

	return router
}
