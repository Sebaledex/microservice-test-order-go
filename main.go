package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"microservice-test-order-go/config"
	"microservice-test-order-go/handlers"
	"microservice-test-order-go/rabbitmq"
	"microservice-test-order-go/repositories"
	"microservice-test-order-go/services"
)

func main() {
	err := godotenv.Load(".env.development")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	rabbitmq.Connect()
	db := config.ConnectMongoDB()

	repo := repositories.NewOrderRepository(db)
	service := services.NewOrderService(repo)
	rabbitmq.StartOrderConsumer(repo)

	r := mux.NewRouter()
	handlers.RegisterOrderRoutes(r, service)

	port := os.Getenv("PORT")
	log.Printf("Server listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
