package repository

import (
	"errors"
	"felix1234567890/go-trello/models"
	"time"

	"gorm.io/gorm"
)

// EventRepository defines the interface for event data operations
type EventRepository interface {
	CreateEvent(event *models.Event) (uint, error)
	GetEventByID(id uint) (*models.Event, error)
	GetAllEvents() ([]models.Event, error)
	UpdateEvent(id uint, eventData *models.UpdateEventRequest) (*models.Event, error)
	DeleteEvent(id uint) error
	AddUserToEvent(eventID uint, userID uint) error
	RemoveUserFromEvent(eventID uint, userID uint) error
	AddGroupToEvent(eventID uint, groupID uint) error
	RemoveGroupFromEvent(eventID uint, groupID uint) error
}

type eventRepositoryImpl struct {
	db *gorm.DB
}

// NewEventRepository creates a new instance of EventRepository
func NewEventRepository(db *gorm.DB) EventRepository {
	return &eventRepositoryImpl{db: db}
}

// CreateEvent creates a new event record in the database
func (r *eventRepositoryImpl) CreateEvent(event *models.Event) (uint, error) {
	result := r.db.Create(event)
	if result.Error != nil {
		return 0, result.Error
	}
	return event.ID, nil
}

// GetEventByID retrieves an event by its ID, preloading Users and Groups
func (r *eventRepositoryImpl) GetEventByID(id uint) (*models.Event, error) {
	var event models.Event
	err := r.db.Preload("Users").Preload("Groups").First(&event, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &event, nil
}

// GetAllEvents retrieves all events, preloading Users and Groups
func (r *eventRepositoryImpl) GetAllEvents() ([]models.Event, error) {
	var events []models.Event
	err := r.db.Preload("Users").Preload("Groups").Find(&events).Error
	if err != nil {
		return nil, err
	}
	return events, nil
}

// UpdateEvent updates an existing event's information
func (r *eventRepositoryImpl) UpdateEvent(id uint, eventData *models.UpdateEventRequest) (*models.Event, error) {
	var event models.Event
	if err := r.db.First(&event, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}

	// Update fields from eventData
	if eventData.Name != "" {
		event.Name = eventData.Name
	}
	if eventData.Description != "" {
		event.Description = eventData.Description
	}
	if eventData.Location != "" {
		event.Location = eventData.Location
	}
	if eventData.Date != "" {
		parsedDate, err := time.Parse(time.RFC3339, eventData.Date)
		if err != nil {
			parsedDate, err = time.Parse("2006-01-02", eventData.Date)
			if err != nil {
				return nil, errors.New("invalid date format for update")
			}
		}
		event.Date = parsedDate
	}

	if err := r.db.Save(&event).Error; err != nil {
		return nil, err
	}
	return &event, nil
}

// DeleteEvent deletes an event by its ID
func (r *eventRepositoryImpl) DeleteEvent(id uint) error {
	result := r.db.Delete(&models.Event{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// AddUserToEvent associates a user with an event
func (r *eventRepositoryImpl) AddUserToEvent(eventID uint, userID uint) error {
	var event models.Event
	if err := r.db.First(&event, eventID).Error; err != nil {
		return err // Handles gorm.ErrRecordNotFound
	}

	var user models.User
	if err := r.db.First(&user, userID).Error; err != nil {
		return err // Handles gorm.ErrRecordNotFound
	}

	return r.db.Model(&event).Association("Users").Append(&user)
}

// RemoveUserFromEvent dissociates a user from an event
func (r *eventRepositoryImpl) RemoveUserFromEvent(eventID uint, userID uint) error {
	var event models.Event
	if err := r.db.First(&event, eventID).Error; err != nil {
		return err
	}

	var user models.User
	if err := r.db.First(&user, userID).Error; err != nil {
		return err
	}

	return r.db.Model(&event).Association("Users").Delete(&user)
}

// AddGroupToEvent associates a group with an event
func (r *eventRepositoryImpl) AddGroupToEvent(eventID uint, groupID uint) error {
	var event models.Event
	if err := r.db.First(&event, eventID).Error; err != nil {
		return err
	}

	var group models.Group
	if err := r.db.First(&group, groupID).Error; err != nil {
		return err
	}

	return r.db.Model(&event).Association("Groups").Append(&group)
}

// RemoveGroupFromEvent dissociates a group from an event
func (r *eventRepositoryImpl) RemoveGroupFromEvent(eventID uint, groupID uint) error {
	var event models.Event
	if err := r.db.First(&event, eventID).Error; err != nil {
		return err
	}

	var group models.Group
	if err := r.db.First(&group, groupID).Error; err != nil {
		return err
	}

	return r.db.Model(&event).Association("Groups").Delete(&group)
}
