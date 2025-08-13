package repository

import (
	"felix1234567890/go-trello/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupEventTestDB initializes an in-memory SQLite database for event repository tests.
func setupEventTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	// AutoMigrate the schema
	err = db.AutoMigrate(&models.User{}, &models.Group{}, &models.Event{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}
	return db
}

func TestEventRepository_CreateEvent(t *testing.T) {
	db := setupEventTestDB(t)
	repo := NewEventRepository(db)

	event := &models.Event{
		Name:        "Test Event 1",
		Description: "Description for Test Event 1",
		Date:        time.Now().Add(24 * time.Hour),
		Location:    "Location 1",
	}

	id, err := repo.CreateEvent(event)
	assert.NoError(t, err)
	assert.NotZero(t, id)
	assert.Equal(t, event.ID, id) // GORM should populate the ID in the original struct

	var foundEvent models.Event
	err = db.First(&foundEvent, id).Error
	assert.NoError(t, err)
	assert.Equal(t, event.Name, foundEvent.Name)
}

func TestEventRepository_GetEventByID(t *testing.T) {
	db := setupEventTestDB(t)
	repo := NewEventRepository(db)

	// Setup: Create an event with users and groups
	user1 := models.User{Username: "user1", Email: "u1@e.com", Password: "p1"}
	user2 := models.User{Username: "user2", Email: "u2@e.com", Password: "p2"}
	db.Create(&user1)
	db.Create(&user2)

	group1 := models.Group{Name: "group1"}
	group2 := models.Group{Name: "group2"}
	db.Create(&group1)
	db.Create(&group2)

	eventDate := time.Now().Add(48 * time.Hour)
	event := &models.Event{
		Name:        "Event With Associations",
		Description: "Test associations",
		Date:        eventDate,
		Location:    "Assoc Location",
		Users:       []models.User{user1},
		Groups:      []models.Group{group1},
	}
	db.Create(event) // Create event directly with associations

	t.Run("found", func(t *testing.T) {
		foundEvent, err := repo.GetEventByID(event.ID)
		assert.NoError(t, err)
		assert.NotNil(t, foundEvent)
		assert.Equal(t, event.Name, foundEvent.Name)
		assert.WithinDuration(t, eventDate, foundEvent.Date, time.Second, "Event date doesn't match")


		// Check preloaded Users
		assert.Len(t, foundEvent.Users, 1, "Should have 1 user preloaded")
		if len(foundEvent.Users) > 0 {
			assert.Equal(t, user1.ID, foundEvent.Users[0].ID)
		}

		// Check preloaded Groups
		assert.Len(t, foundEvent.Groups, 1, "Should have 1 group preloaded")
		if len(foundEvent.Groups) > 0 {
			assert.Equal(t, group1.ID, foundEvent.Groups[0].ID)
		}
	})

	t.Run("not found", func(t *testing.T) {
		_, err := repo.GetEventByID(99999) // Non-existent ID
		assert.Error(t, err)
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	})
}

func TestEventRepository_GetAllEvents(t *testing.T) {
	db := setupEventTestDB(t)
	repo := NewEventRepository(db)

	// Setup: Create some events
	event1 := models.Event{Name: "Event A", Date: time.Now()}
	event2 := models.Event{Name: "Event B", Date: time.Now().Add(24 * time.Hour), Users: []models.User{{Username: "u1"}}}
	db.Create(&event1)
	db.Create(&event2)

	events, err := repo.GetAllEvents()
	assert.NoError(t, err)
	assert.Len(t, events, 2)

	// Verify that associations are preloaded (e.g., Users for event2)
	for _, ev := range events {
		if ev.Name == "Event B" {
			assert.Len(t, ev.Users, 1, "Users should be preloaded for Event B")
		}
	}
}

func TestEventRepository_UpdateEvent(t *testing.T) {
	db := setupEventTestDB(t)
	repo := NewEventRepository(db)

	event := models.Event{Name: "Old Name", Date: time.Now(), Location: "Old Location"}
	db.Create(&event)

	updateData := &models.UpdateEventRequest{
		Name:     "New Name",
		Location: "New Location",
		Date:     time.Now().Add(72 * time.Hour).Format(time.RFC3339),
	}

	t.Run("existing event", func(t *testing.T) {
		updatedEvent, err := repo.UpdateEvent(event.ID, updateData)
		assert.NoError(t, err)
		assert.NotNil(t, updatedEvent)
		assert.Equal(t, "New Name", updatedEvent.Name)
		assert.Equal(t, "New Location", updatedEvent.Location)

		parsedExpectedDate, _ := time.Parse(time.RFC3339, updateData.Date)
		assert.WithinDuration(t, parsedExpectedDate, updatedEvent.Date, time.Second)

		// Verify in DB
		var dbEvent models.Event
		db.First(&dbEvent, event.ID)
		assert.Equal(t, "New Name", dbEvent.Name)
	})

	t.Run("non-existent event", func(t *testing.T) {
		_, err := repo.UpdateEvent(99999, updateData)
		assert.Error(t, err)
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	})

	t.Run("update with only some fields", func(t *testing.T) {
		eventPartial := models.Event{Name: "Partial Original", Date: time.Now(), Description: "Original Desc"}
		db.Create(&eventPartial)

		partialUpdateData := &models.UpdateEventRequest{
			Name: "Partial New Name",
		}
		updatedEvent, err := repo.UpdateEvent(eventPartial.ID, partialUpdateData)
		assert.NoError(t, err)
		assert.Equal(t, "Partial New Name", updatedEvent.Name)
		assert.Equal(t, "Original Desc", updatedEvent.Description) // Description should not change
	})
}

func TestEventRepository_DeleteEvent(t *testing.T) {
	db := setupEventTestDB(t)
	repo := NewEventRepository(db)

	event := models.Event{Name: "To Be Deleted", Date: time.Now()}
	db.Create(&event)

	t.Run("existing event", func(t *testing.T) {
		err := repo.DeleteEvent(event.ID)
		assert.NoError(t, err)

		_, err = repo.GetEventByID(event.ID) // Should not be found
		assert.Error(t, err)
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	})

	t.Run("non-existent event", func(t *testing.T) {
		err := repo.DeleteEvent(99999)
		assert.Error(t, err)
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	})
}

func TestEventRepository_UserEventAssociations(t *testing.T) {
	db := setupEventTestDB(t)
	repo := NewEventRepository(db)

	user := models.User{Username: "eventUser", Email: "eu@e.com", Password: "p"}
	db.Create(&user)
	event := models.Event{Name: "User Assoc Event", Date: time.Now()}
	db.Create(&event)

	t.Run("AddUserToEvent", func(t *testing.T) {
		err := repo.AddUserToEvent(event.ID, user.ID)
		assert.NoError(t, err)

		foundEvent, _ := repo.GetEventByID(event.ID)
		assert.Len(t, foundEvent.Users, 1)
		if len(foundEvent.Users) > 0 {
			assert.Equal(t, user.ID, foundEvent.Users[0].ID)
		}

		// Test adding non-existent event
		err = repo.AddUserToEvent(999, user.ID)
		assert.Error(t, err)
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)


		// Test adding non-existent user
		err = repo.AddUserToEvent(event.ID, 888)
		assert.Error(t, err)
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	})

	t.Run("RemoveUserFromEvent", func(t *testing.T) {
		// Ensure user is added first (though previous test does this, good to be explicit if tests run in parallel or are separated)
		db.Model(&event).Association("Users").Append(&user)

		err := repo.RemoveUserFromEvent(event.ID, user.ID)
		assert.NoError(t, err)

		foundEvent, _ := repo.GetEventByID(event.ID)
		assert.Len(t, foundEvent.Users, 0)

		// Test removing non-existent event (or user from it)
		err = repo.RemoveUserFromEvent(999, user.ID)
		assert.Error(t, err)
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)

		// Test removing non-existent user from event
		err = repo.RemoveUserFromEvent(event.ID, 888)
		assert.Error(t, err)
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)

		// Test removing user not associated with event (should not error, GORM handles this gracefully)
		anotherUser := models.User{Username: "anotherU", Email: "au@e.com", Password: "p"}
		db.Create(&anotherUser)
		err = repo.RemoveUserFromEvent(event.ID, anotherUser.ID)
		assert.NoError(t, err)

	})
}

func TestEventRepository_GroupEventAssociations(t *testing.T) {
	db := setupEventTestDB(t)
	repo := NewEventRepository(db)

	group := models.Group{Name: "eventGroup"}
	db.Create(&group)
	event := models.Event{Name: "Group Assoc Event", Date: time.Now()}
	db.Create(&event)

	t.Run("AddGroupToEvent", func(t *testing.T) {
		err := repo.AddGroupToEvent(event.ID, group.ID)
		assert.NoError(t, err)

		foundEvent, _ := repo.GetEventByID(event.ID)
		assert.Len(t, foundEvent.Groups, 1)
		if len(foundEvent.Groups) > 0 {
			assert.Equal(t, group.ID, foundEvent.Groups[0].ID)
		}

		// Test adding non-existent event
		err = repo.AddGroupToEvent(999, group.ID)
		assert.Error(t, err)
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)

		// Test adding non-existent group
		err = repo.AddGroupToEvent(event.ID, 888)
		assert.Error(t, err)
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	})

	t.Run("RemoveGroupFromEvent", func(t *testing.T) {
		db.Model(&event).Association("Groups").Append(&group)

		err := repo.RemoveGroupFromEvent(event.ID, group.ID)
		assert.NoError(t, err)

		foundEvent, _ := repo.GetEventByID(event.ID)
		assert.Len(t, foundEvent.Groups, 0)

		// Test removing non-existent event
		err = repo.RemoveGroupFromEvent(999, group.ID)
		assert.Error(t, err)
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)

		// Test removing non-existent group
		err = repo.RemoveGroupFromEvent(event.ID, 888)
		assert.Error(t, err)
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)

		// Test removing group not associated (should not error)
		anotherGroup := models.Group{Name: "anotherG"}
		db.Create(&anotherGroup)
		err = repo.RemoveGroupFromEvent(event.ID, anotherGroup.ID)
		assert.NoError(t, err)
	})
}

func TestEventRepository_UpdateEvent_InvalidDateFormat(t *testing.T) {
	db := setupEventTestDB(t)
	repo := NewEventRepository(db)

	event := models.Event{Name: "Date Event", Date: time.Now()}
	db.Create(&event)

	updateData := &models.UpdateEventRequest{
		Date: "invalid-date-format", // Not RFC3339 or YYYY-MM-DD
	}

	_, err := repo.UpdateEvent(event.ID, updateData)
	assert.Error(t, err) // Expect an error due to date parsing
	assert.Contains(t, err.Error(), "invalid date format for update")
}
