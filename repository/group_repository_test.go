package repository

import (
	"testing"
	"felix1234567890/go-trello/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"github.com/stretchr/testify/assert"
)

func setupTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&models.User{}, &models.Group{})
	return db
}

func TestGroupCRUD(t *testing.T) {
	db := setupTestDB()
	repo := NewGroupRepository(db)
	// Create
	group := &models.Group{Name: "TestGroup"}
	err := repo.CreateGroup(group)
	assert.NoError(t, err)
	assert.NotZero(t, group.ID)
	// Read
	g, err := repo.GetGroupByID(group.ID)
	assert.NoError(t, err)
	assert.Equal(t, "TestGroup", g.Name)
	// Update
	group.Name = "UpdatedGroup"
	err = repo.UpdateGroup(group)
	assert.NoError(t, err)
	g, _ = repo.GetGroupByID(group.ID)
	assert.Equal(t, "UpdatedGroup", g.Name)
	// Delete
	err = repo.DeleteGroup(group.ID)
	assert.NoError(t, err)
	_, err = repo.GetGroupByID(group.ID)
	assert.Error(t, err)
}

func TestAddAndRemoveUserToGroup(t *testing.T) {
	db := setupTestDB()
	repo := NewGroupRepository(db)
	user := &models.User{Username: "user1", Email: "user1@email.com", Password: "pass"}
	group := &models.Group{Name: "Group1"}
	db.Create(user)
	repo.CreateGroup(group)
	// Add user to group
	err := repo.AddUserToGroup(group.ID, user.ID)
	assert.NoError(t, err)
	g, _ := repo.GetGroupByID(group.ID)
	assert.Equal(t, 1, len(g.Users))
	assert.Equal(t, user.ID, g.Users[0].ID)
	// Remove user from group
	err = repo.RemoveUserFromGroup(group.ID, user.ID)
	assert.NoError(t, err)
	g, _ = repo.GetGroupByID(group.ID)
	assert.Equal(t, 0, len(g.Users))
}
