package service

import (
	"felix1234567890/go-trello/models"
	"felix1234567890/go-trello/repository"
)

// EventService defines the interface for event business logic
type EventService interface {
	CreateEvent(request *models.CreateEventRequest) (*models.Event, error)
	GetEventByID(id uint) (*models.Event, error)
	GetAllEvents() ([]models.Event, error)
	UpdateEvent(id uint, request *models.UpdateEventRequest) (*models.Event, error)
	DeleteEvent(id uint) error
	AddUserToEvent(eventID uint, userID uint) error
	RemoveUserFromEvent(eventID uint, userID uint) error
	AddGroupToEvent(eventID uint, groupID uint) error
	RemoveGroupFromEvent(eventID uint, groupID uint) error
}

type eventServiceImpl struct {
	eventRepo repository.EventRepository
}

// NewEventService creates a new instance of EventService
func NewEventService(eventRepo repository.EventRepository) EventService {
	return &eventServiceImpl{eventRepo: eventRepo}
}

// CreateEvent handles the business logic for creating a new event
func (s *eventServiceImpl) CreateEvent(request *models.CreateEventRequest) (*models.Event, error) {
	event, err := request.ToEvent()
	if err != nil {
		return nil, err // Error from date parsing in ToEvent
	}

	eventID, err := s.eventRepo.CreateEvent(event)
	if err != nil {
		return nil, err
	}

	// Retrieve the full event details to return the complete object
	return s.eventRepo.GetEventByID(eventID)
}

// GetEventByID retrieves an event by its ID
func (s *eventServiceImpl) GetEventByID(id uint) (*models.Event, error) {
	return s.eventRepo.GetEventByID(id)
}

// GetAllEvents retrieves all events
func (s *eventServiceImpl) GetAllEvents() ([]models.Event, error) {
	return s.eventRepo.GetAllEvents()
}

// UpdateEvent handles the business logic for updating an event
func (s *eventServiceImpl) UpdateEvent(id uint, request *models.UpdateEventRequest) (*models.Event, error) {
	// The repository's UpdateEvent method is designed to take UpdateEventRequest directly
	// and handles fetching and updating the event.
	return s.eventRepo.UpdateEvent(id, request)
}

// DeleteEvent handles the business logic for deleting an event
func (s *eventServiceImpl) DeleteEvent(id uint) error {
	return s.eventRepo.DeleteEvent(id)
}

// AddUserToEvent handles adding a user to an event
func (s *eventServiceImpl) AddUserToEvent(eventID uint, userID uint) error {
	// Future enhancement: check if user or event exists, or if user is already in event.
	return s.eventRepo.AddUserToEvent(eventID, userID)
}

// RemoveUserFromEvent handles removing a user from an event
func (s *eventServiceImpl) RemoveUserFromEvent(eventID uint, userID uint) error {
	return s.eventRepo.RemoveUserFromEvent(eventID, userID)
}

// AddGroupToEvent handles adding a group to an event
func (s *eventServiceImpl) AddGroupToEvent(eventID uint, groupID uint) error {
	// Future enhancement: check if group or event exists, or if group is already in event.
	return s.eventRepo.AddGroupToEvent(eventID, groupID)
}

// RemoveGroupFromEvent handles removing a group from an event
func (s *eventServiceImpl) RemoveGroupFromEvent(eventID uint, groupID uint) error {
	return s.eventRepo.RemoveGroupFromEvent(eventID, groupID)
}
