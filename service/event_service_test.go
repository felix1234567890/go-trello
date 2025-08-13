package service

import (
	"errors"
	"felix1234567890/go-trello/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockEventRepository is a mock implementation of repository.EventRepository
type MockEventRepository struct {
	mock.Mock
}

func (m *MockEventRepository) CreateEvent(event *models.Event) (uint, error) {
	args := m.Called(event)
	return args.Get(0).(uint), args.Error(1)
}

func (m *MockEventRepository) GetEventByID(id uint) (*models.Event, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Event), args.Error(1)
}

func (m *MockEventRepository) GetAllEvents() ([]models.Event, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Event), args.Error(1)
}

func (m *MockEventRepository) UpdateEvent(id uint, eventData *models.UpdateEventRequest) (*models.Event, error) {
	args := m.Called(id, eventData)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Event), args.Error(1)
}

func (m *MockEventRepository) DeleteEvent(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockEventRepository) AddUserToEvent(eventID uint, userID uint) error {
	args := m.Called(eventID, userID)
	return args.Error(0)
}

func (m *MockEventRepository) RemoveUserFromEvent(eventID uint, userID uint) error {
	args := m.Called(eventID, userID)
	return args.Error(0)
}

func (m *MockEventRepository) AddGroupToEvent(eventID uint, groupID uint) error {
	args := m.Called(eventID, groupID)
	return args.Error(0)
}

func (m *MockEventRepository) RemoveGroupFromEvent(eventID uint, groupID uint) error {
	args := m.Called(eventID, groupID)
	return args.Error(0)
}

func TestEventService_CreateEvent(t *testing.T) {
	mockRepo := new(MockEventRepository)
	service := NewEventService(mockRepo)
	now := time.Now()

	createRequest := &models.CreateEventRequest{
		Name:        "Test Event",
		Description: "A cool event",
		Date:        now.Format(time.RFC3339),
		Location:    "Online",
	}

	expectedEventModel, _ := createRequest.ToEvent() // Use the same conversion logic
	expectedEventModel.ID = 1 // Assume repo assigns ID 1

	t.Run("success", func(t *testing.T) {
		// We need to mock ToEvent's behavior for the input event
		// then the repo.CreateEvent, then repo.GetEventByID
		eventFromRequest, _ := createRequest.ToEvent() // This is what service will pass to CreateEvent

		mockRepo.On("CreateEvent", eventFromRequest).Return(uint(1), nil).Once()
		mockRepo.On("GetEventByID", uint(1)).Return(expectedEventModel, nil).Once()

		createdEvent, err := service.CreateEvent(createRequest)

		assert.NoError(t, err)
		assert.NotNil(t, createdEvent)
		assert.Equal(t, expectedEventModel.Name, createdEvent.Name)
		assert.Equal(t, uint(1), createdEvent.ID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid date format in request", func(t *testing.T) {
		invalidDateRequest := &models.CreateEventRequest{
			Name: "Event with Bad Date",
			Date: "not-a-date",
		}
		// No mocks needed as ToEvent() should fail before repo calls

		_, err := service.CreateEvent(invalidDateRequest)
		assert.Error(t, err)
		// Check if the error is from time.Parse or a custom error message from ToEvent
		assert.Contains(t, err.Error(), "cannot parse") // time.Parse error
		mockRepo.AssertNotCalled(t, "CreateEvent", mock.Anything)
		mockRepo.AssertNotCalled(t, "GetEventByID", mock.Anything)
	})

	t.Run("repository CreateEvent fails", func(t *testing.T) {
		eventFromRequest, _ := createRequest.ToEvent()
		mockRepo.On("CreateEvent", eventFromRequest).Return(uint(0), errors.New("db error")).Once()

		_, err := service.CreateEvent(createRequest)
		assert.Error(t, err)
		assert.Equal(t, "db error", err.Error())
		mockRepo.AssertExpectations(t)
		mockRepo.AssertNotCalled(t, "GetEventByID", mock.Anything)
	})

	t.Run("repository GetEventByID fails after create", func(t *testing.T) {
		eventFromRequest, _ := createRequest.ToEvent()
		mockRepo.On("CreateEvent", eventFromRequest).Return(uint(1), nil).Once()
		mockRepo.On("GetEventByID", uint(1)).Return(nil, errors.New("get by id error")).Once()

		_, err := service.CreateEvent(createRequest)
		assert.Error(t, err)
		assert.Equal(t, "get by id error", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestEventService_GetEventByID(t *testing.T) {
	mockRepo := new(MockEventRepository)
	service := NewEventService(mockRepo)
	expectedEvent := &models.Event{ID: 1, Name: "Found Event"}

	t.Run("found", func(t *testing.T) {
		mockRepo.On("GetEventByID", uint(1)).Return(expectedEvent, nil).Once()
		event, err := service.GetEventByID(uint(1))
		assert.NoError(t, err)
		assert.Equal(t, expectedEvent, event)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo.On("GetEventByID", uint(2)).Return(nil, gorm.ErrRecordNotFound).Once()
		_, err := service.GetEventByID(uint(2))
		assert.Error(t, err)
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
		mockRepo.AssertExpectations(t)
	})
}

func TestEventService_GetAllEvents(t *testing.T) {
	mockRepo := new(MockEventRepository)
	service := NewEventService(mockRepo)
	expectedEvents := []models.Event{{ID: 1, Name: "Event 1"}, {ID: 2, Name: "Event 2"}}

	mockRepo.On("GetAllEvents").Return(expectedEvents, nil).Once()
	events, err := service.GetAllEvents()

	assert.NoError(t, err)
	assert.Equal(t, expectedEvents, events)
	mockRepo.AssertExpectations(t)
}

func TestEventService_UpdateEvent(t *testing.T) {
	mockRepo := new(MockEventRepository)
	service := NewEventService(mockRepo)
	nowStr := time.Now().Format(time.RFC3339)

	updateRequest := &models.UpdateEventRequest{
		Name: "Updated Name",
		Date: nowStr,
	}
	// The service's UpdateEvent passes UpdateEventRequest directly to repository's UpdateEvent.
	// The repository's UpdateEvent is responsible for parsing the date string.
	// So, we don't test ToEvent() directly here for date parsing,
	// that's handled by the repository test or if ToEvent() was called in service.
	// In this implementation, service.UpdateEvent calls repo.UpdateEvent(id, request)
	// So, no ToEvent() call in the service layer for UpdateEvent.

	expectedUpdatedEvent := &models.Event{ID: 1, Name: "Updated Name"} // Repo would return this

	t.Run("success", func(t *testing.T) {
		mockRepo.On("UpdateEvent", uint(1), updateRequest).Return(expectedUpdatedEvent, nil).Once()
		updatedEvent, err := service.UpdateEvent(uint(1), updateRequest)
		assert.NoError(t, err)
		assert.Equal(t, expectedUpdatedEvent, updatedEvent)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository UpdateEvent fails", func(t *testing.T) {
		mockRepo.On("UpdateEvent", uint(1), updateRequest).Return(nil, errors.New("db update error")).Once()
		_, err := service.UpdateEvent(uint(1), updateRequest)
		assert.Error(t, err)
		assert.Equal(t, "db update error", err.Error())
		mockRepo.AssertExpectations(t)
	})

    t.Run("repository UpdateEvent returns GormErrRecordNotFound", func(t *testing.T) {
		mockRepo.On("UpdateEvent", uint(99), updateRequest).Return(nil, gorm.ErrRecordNotFound).Once()
		_, err := service.UpdateEvent(uint(99), updateRequest)
		assert.Error(t, err)
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
		mockRepo.AssertExpectations(t)
	})
}


func TestEventService_DeleteEvent(t *testing.T) {
	mockRepo := new(MockEventRepository)
	service := NewEventService(mockRepo)

	t.Run("success", func(t *testing.T) {
		mockRepo.On("DeleteEvent", uint(1)).Return(nil).Once()
		err := service.DeleteEvent(uint(1))
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository DeleteEvent fails", func(t *testing.T) {
		mockRepo.On("DeleteEvent", uint(1)).Return(errors.New("db delete error")).Once()
		err := service.DeleteEvent(uint(1))
		assert.Error(t, err)
		assert.Equal(t, "db delete error", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestEventService_AssociationMethods(t *testing.T) {
	mockRepo := new(MockEventRepository)
	service := NewEventService(mockRepo)
	eventID, userID, groupID := uint(1), uint(10), uint(100)

	// AddUserToEvent
	mockRepo.On("AddUserToEvent", eventID, userID).Return(nil).Once()
	err := service.AddUserToEvent(eventID, userID)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t) // Reset for next call

	mockRepo.On("AddUserToEvent", eventID, uint(11)).Return(errors.New("add user failed")).Once()
	err = service.AddUserToEvent(eventID, uint(11))
	assert.Error(t, err)
	assert.Equal(t, "add user failed", err.Error())
	mockRepo.AssertExpectations(t)

	// RemoveUserFromEvent
	mockRepo.On("RemoveUserFromEvent", eventID, userID).Return(nil).Once()
	err = service.RemoveUserFromEvent(eventID, userID)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)

	// AddGroupToEvent
	mockRepo.On("AddGroupToEvent", eventID, groupID).Return(nil).Once()
	err = service.AddGroupToEvent(eventID, groupID)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)

	// RemoveGroupFromEvent
	mockRepo.On("RemoveGroupFromEvent", eventID, groupID).Return(nil).Once()
	err = service.RemoveGroupFromEvent(eventID, groupID)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
