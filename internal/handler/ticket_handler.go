package handler

import (
	"net/http"

	"ticketing-service/internal/models"
	"ticketing-service/internal/service"

	"github.com/gin-gonic/gin"
)

type TicketHandler struct {
	ticketService *service.TicketService
}

func NewTicketHandler(ticketService *service.TicketService) *TicketHandler {
	return &TicketHandler{ticketService: ticketService}
}

func (h *TicketHandler) GenerateTickets(c *gin.Context) {
	var req models.GenerateTicketsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tickets, err := h.ticketService.GenerateTickets(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"tickets": tickets})
}

func (h *TicketHandler) GetUserTickets(c *gin.Context) {
	userID := c.Param("userId")

	tickets, err := h.ticketService.GetUserTickets(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tickets": tickets})
}

func (h *TicketHandler) ValidateTicket(c *gin.Context) {
	ticketCode := c.Param("ticketCode")

	ticket, err := h.ticketService.ValidateTicket(ticketCode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ticket": ticket, "valid": ticket.Status == "valid"})
}

func (h *TicketHandler) UseTicket(c *gin.Context) {
	ticketCode := c.Param("ticketCode")

	ticket, err := h.ticketService.UseTicket(ticketCode)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ticket": ticket, "message": "Ticket used successfully"})
}
