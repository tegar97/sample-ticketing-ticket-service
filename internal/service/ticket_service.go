package service

import (
	"errors"
	"ticketing-service/internal/models"
	"ticketing-service/internal/repository"
)

type TicketService struct {
	ticketRepo *repository.TicketRepository
}

func NewTicketService(ticketRepo *repository.TicketRepository) *TicketService {
	return &TicketService{ticketRepo: ticketRepo}
}

func (s *TicketService) GenerateTickets(req *models.GenerateTicketsRequest) ([]*models.Ticket, error) {
	if req.Quantity <= 0 {
		return nil, errors.New("quantity must be greater than 0")
	}

	tickets := make([]*models.Ticket, req.Quantity)
	for i := 0; i < req.Quantity; i++ {
		tickets[i] = &models.Ticket{
			BookingID: req.BookingID,
			EventID:   req.EventID,
			UserID:    req.UserID,
		}
	}

	err := s.ticketRepo.CreateTickets(tickets)
	if err != nil {
		return nil, err
	}

	return tickets, nil
}

func (s *TicketService) GetUserTickets(userID string) ([]*models.TicketWithEvent, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	return s.ticketRepo.GetUserTickets(userID)
}

func (s *TicketService) ValidateTicket(ticketCode string) (*models.TicketWithEvent, error) {
	if ticketCode == "" {
		return nil, errors.New("ticket code is required")
	}

	ticket, err := s.ticketRepo.GetByTicketCode(ticketCode)
	if err != nil {
		return nil, errors.New("ticket not found")
	}

	return ticket, nil
}

func (s *TicketService) UseTicket(ticketCode string) (*models.TicketWithEvent, error) {
	if ticketCode == "" {
		return nil, errors.New("ticket code is required")
	}

	ticket, err := s.ticketRepo.GetByTicketCode(ticketCode)
	if err != nil {
		return nil, errors.New("ticket not found")
	}

	if ticket.Status != "valid" {
		return nil, errors.New("ticket is not valid")
	}

	err = s.ticketRepo.UpdateStatus(ticketCode, "used")
	if err != nil {
		return nil, err
	}

	ticket.Status = "used"
	return ticket, nil
}
