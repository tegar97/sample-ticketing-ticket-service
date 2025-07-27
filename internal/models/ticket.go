package models

import (
	"time"
)

type Ticket struct {
	ID         string    `json:"id" gorm:"primaryKey;type:varchar(255)"`
	BookingID  string    `json:"booking_id" gorm:"type:varchar(255);not null"`
	EventID    string    `json:"event_id" gorm:"type:varchar(255);not null"`
	UserID     string    `json:"user_id" gorm:"type:varchar(255);not null"`
	TicketCode string    `json:"ticket_code" gorm:"type:varchar(255);uniqueIndex;not null"`
	Status     string    `json:"status" gorm:"type:varchar(50);default:'valid'"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Event      Event     `json:"event,omitempty" gorm:"foreignKey:EventID;references:ID"`
}

type Event struct {
	ID        string    `json:"id" gorm:"primaryKey;type:varchar(255)"`
	Title     string    `json:"title" gorm:"type:varchar(255);not null"`
	Venue     string    `json:"venue" gorm:"type:varchar(255);not null"`
	EventDate time.Time `json:"event_date" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type GenerateTicketsRequest struct {
	BookingID string `json:"booking_id" binding:"required"`
	EventID   string `json:"event_id" binding:"required"`
	UserID    string `json:"user_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,min=1"`
}

type TicketWithEvent struct {
	Ticket
	EventTitle string    `json:"event_title"`
	EventVenue string    `json:"event_venue"`
	EventDate  time.Time `json:"event_date"`
}
