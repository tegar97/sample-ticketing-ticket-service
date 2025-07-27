package main

import (
	"log"
	"net/http"
	"os"

	"ticketing-service/internal/config"
	"ticketing-service/internal/handler"
	"ticketing-service/internal/repository"
	"ticketing-service/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	db, err := config.ConnectDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	ticketRepo := repository.NewTicketRepository(db)
	ticketService := service.NewTicketService(ticketRepo)
	ticketHandler := handler.NewTicketHandler(ticketService)

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	})

	api := r.Group("/api/v1")
	{
		api.POST("/tickets/generate", ticketHandler.GenerateTickets)
		api.GET("/tickets/user/:userId", ticketHandler.GetUserTickets)
		api.GET("/tickets/:ticketCode/validate", ticketHandler.ValidateTicket)
		api.PUT("/tickets/:ticketCode/use", ticketHandler.UseTicket)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8004"
	}

	log.Printf("Ticketing service starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
