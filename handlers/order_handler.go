package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"microservice-test-order-go/models"
	"microservice-test-order-go/services"
)

func RegisterOrderRoutes(r *mux.Router, service *services.OrderService) {
	// Crear una orden
	r.HandleFunc("/api/v2/order", func(w http.ResponseWriter, r *http.Request) {
		var order models.Order
		if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}
		if err := service.CreateOrder(r.Context(), &order); err != nil {
			http.Error(w, "Failed to create order: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(order)
	}).Methods("POST")

	// Listar todas las órdenes
	r.HandleFunc("/api/v2/order", func(w http.ResponseWriter, r *http.Request) {
		orders, err := service.GetAllOrders(r.Context())
		if err != nil {
			http.Error(w, "Failed to fetch orders", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(orders)
	}).Methods("GET")

	// Obtener una orden por ID
	r.HandleFunc("/api/v2/order/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		order, err := service.GetOrderById(r.Context(), id)
		if err != nil {
			http.Error(w, "Order not found: "+err.Error(), http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(order)
	}).Methods("GET")

	// Actualizar una orden por ID
	r.HandleFunc("/api/v2/order/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		var order models.Order
		if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}
		if err := service.UpdateOrder(r.Context(), id, &order); err != nil {
			http.Error(w, "Failed to update order: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}).Methods("PUT")

	// Eliminar una orden por ID
	r.HandleFunc("/api/v2/order/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		if err := service.DeleteOrder(r.Context(), id); err != nil {
			http.Error(w, "Failed to delete order: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}).Methods("DELETE")
}
