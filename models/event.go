package models

import (
	"time"

	"gorm.io/gorm"
)

// Event represents an event
type Event struct {
	gorm.Model
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	Date        time.Time `json:"date" gorm:"not null"`
	Location    string    `json:"location"`
	Users       []User    `json:"users" gorm:"many2many:event_users;"`
	Groups      []Group   `json:"groups" gorm:"many2many:event_groups;"`
}

// CreateEventRequest represents the request body for creating an event
type CreateEventRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Date        string `json:"date" validate:"required" comment:"Use string for request, parse to time.Time in service"`
	Location    string `json:"location"`
}

// UpdateEventRequest represents the request body for updating an event
type UpdateEventRequest struct {
	Name        string `json:"name,omitempty" validate:"omitempty,min=1"`
	Description string `json:"description,omitempty"`
	Date        string `json:"date,omitempty" comment:"Use string for request, parse to time.Time in service"`
	Location    string `json:"location,omitempty"`
}

// ToEvent converts CreateEventRequest to Event model
func (r *CreateEventRequest) ToEvent() (*Event, error) {
	date, err := time.Parse(time.RFC3339, r.Date)
	if err != nil {
		// Attempt to parse with just date if RFC3339 fails
		date, err = time.Parse("2006-01-02", r.Date)
		if err != nil {
			return nil, err
		}
	}
	return &Event{
		Name:        r.Name,
		Description: r.Description,
		Date:        date,
		Location:    r.Location,
	}, nil
}

// ToEvent converts UpdateEventRequest to Event model
// Note: This method creates a new Event struct.
// In a real application, you'd likely fetch the existing event
// and update its fields.
func (r *UpdateEventRequest) ToEvent() (*Event, error) {
	event := &Event{}
	if r.Name != "" {
		event.Name = r.Name
	}
	if r.Description != "" {
		event.Description = r.Description
	}
	if r.Date != "" {
		date, err := time.Parse(time.RFC3339, r.Date)
		if err != nil {
			// Attempt to parse with just date if RFC3339 fails
			date, err = time.Parse("2006-01-02", r.Date)
			if err != nil {
				return nil, err
			}
		}
		event.Date = date
	}
	if r.Location != "" {
		event.Location = r.Location
	}
	return event, nil
}
