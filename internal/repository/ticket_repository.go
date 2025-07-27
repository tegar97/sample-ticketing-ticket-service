package repository

import (
	"errors"
	"fmt"
	"ticketing-service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TicketRepository struct {
	db *gorm.DB
}

func NewTicketRepository(db *gorm.DB) *TicketRepository {
	return &TicketRepository{db: db}
}

func (r *TicketRepository) CreateTickets(tickets []*models.Ticket) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, ticket := range tickets {
			ticket.ID = uuid.New().String()
			ticket.TicketCode = generateTicketCode()
			ticket.Status = "valid"

			if err := tx.Create(ticket).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *TicketRepository) GetUserTickets(userID string) ([]*models.TicketWithEvent, error) {
	var tickets []models.Ticket
	err := r.db.Preload("Event").Where("user_id = ?", userID).Order("created_at DESC").Find(&tickets).Error
	if err != nil {
		return nil, err
	}

	var ticketsWithEvent []*models.TicketWithEvent
	for _, ticket := range tickets {
		ticketWithEvent := &models.TicketWithEvent{
			Ticket:     ticket,
			EventTitle: ticket.Event.Title,
			EventVenue: ticket.Event.Venue,
			EventDate:  ticket.Event.EventDate,
		}
		ticketsWithEvent = append(ticketsWithEvent, ticketWithEvent)
	}

	return ticketsWithEvent, nil
}

func (r *TicketRepository) GetByTicketCode(ticketCode string) (*models.TicketWithEvent, error) {
	var ticket models.Ticket
	err := r.db.Preload("Event").Where("ticket_code = ?", ticketCode).First(&ticket).Error
	if err != nil {
		return nil, err
	}

	ticketWithEvent := &models.TicketWithEvent{
		Ticket:     ticket,
		EventTitle: ticket.Event.Title,
		EventVenue: ticket.Event.Venue,
		EventDate:  ticket.Event.EventDate,
	}

	return ticketWithEvent, nil
}

func (r *TicketRepository) UpdateStatus(ticketCode, status string) error {
	result := r.db.Model(&models.Ticket{}).Where("ticket_code = ?", ticketCode).Update("status", status)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("ticket not found")
	}

	return nil
}

func generateTicketCode() string {
	return fmt.Sprintf("TKT-%s", uuid.New().String()[:8])
}
