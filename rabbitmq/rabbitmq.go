package rabbitmq

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/streadway/amqp"
	"microservice-test-order-go/models"
	"microservice-test-order-go/repositories"
)

var conn *amqp.Connection
var channel *amqp.Channel // Declaramos correctamente el canal

func Connect() {
	var err error

	// Obtener la URL de RabbitMQ desde el archivo .env
	rabbitMQURL := os.Getenv("AMQP_URL")
	if rabbitMQURL == "" {
		log.Fatal("AMQP_URL not set in environment variables")
	}

	// Establecer la conexión con RabbitMQ
	conn, err = amqp.Dial(rabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	// Crear un canal
	channel, err = conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}

	log.Println("Connected to RabbitMQ and channel created")
}

func StartOrderConsumer(repo *repositories.OrderRepository) {
	// Declarar la cola
	queueName := os.Getenv("ORDER_QUEUE")
	if queueName == "" {
		queueName = "orders" // Nombre por defecto si no está configurado
	}

	_, err := channel.QueueDeclare(
		queueName, // Nombre de la cola
		true,      // Durable
		false,     // Auto-deleted
		false,     // Exclusive
		false,     // No-wait
		nil,       // Argumentos adicionales
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	// Consumir mensajes de la cola
	msgs, err := channel.Consume(
		queueName, // Nombre de la cola
		"",        // Nombre del consumidor
		true,      // Auto-ack
		false,     // Exclusive
		false,     // No-local
		false,     // No-wait
		nil,       // Argumentos adicionales
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	go func() {
		ctx := context.Background() // Creamos un contexto vacío
		for d := range msgs {
			// Procesar cada mensaje recibido
			var order models.Order
			err := json.Unmarshal(d.Body, &order)
			if err != nil {
				log.Printf("Error decoding message: %v", err)
				continue
			}

			// Guardar la orden en la base de datos
			err = repo.Create(ctx, &order)
			if err != nil {
				log.Printf("Error saving order to DB: %v", err)
			} else {
				log.Printf("Order saved: %v", order)
			}
		}
	}()
}

func Close() {
	// Cerrar el canal y la conexión al salir
	if channel != nil {
		channel.Close()
	}
	if conn != nil {
		conn.Close()
	}
	log.Println("RabbitMQ connection and channel closed")
}
