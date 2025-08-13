package tests

import (
	"felix1234567890/go-trello/models"
	"felix1234567890/go-trello/utils"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestDatabase represents a test database instance
type TestDatabase struct {
	DB *gorm.DB
}

// SetupTestDB creates a new in-memory SQLite database for testing
func SetupTestDB(t *testing.T) *TestDatabase {
	// Setup test database using SQLite in memory
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err)

	// Migrate the schema
	err = db.AutoMigrate(&models.User{}, &models.Group{})
	require.NoError(t, err)

	return &TestDatabase{DB: db}
}

// TeardownTestDB closes the test database connection
func (td *TestDatabase) TeardownTestDB() {
	if td.DB != nil {
		sqlDB, _ := td.DB.DB()
		sqlDB.Close()
	}
}

// CleanupTestDB cleans all data from test database tables
func (td *TestDatabase) CleanupTestDB() {
	td.DB.Exec("DELETE FROM user_groups")
	td.DB.Exec("DELETE FROM users")
	td.DB.Exec("DELETE FROM groups")
}

// SeedTestUser creates a test user in the database
func (td *TestDatabase) SeedTestUser() *models.User {
	user := &models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword", // This would be hashed in real scenario
	}
	
	td.DB.Create(user)
	return user
}

// SeedTestGroup creates a test group in the database
func (td *TestDatabase) SeedTestGroup() *models.Group {
	group := &models.Group{
		Name: "Test Group",
	}
	
	td.DB.Create(group)
	return group
}

// SeedTestUserWithRealPassword creates a test user with a properly hashed password
func (td *TestDatabase) SeedTestUserWithRealPassword(email, plainPassword string) *models.User {
	hashedPassword, _ := utils.HashPassword(plainPassword)
	user := &models.User{
		Username: "testuser",
		Email:    email,
		Password: hashedPassword,
	}
	
	td.DB.Create(user)
	return user
}

// SetupTestEnvironment initializes test environment variables
func SetupTestEnvironment() {
	os.Setenv("JWT_SECRET", "test-secret-key-for-testing")
	os.Setenv("MYSQL_USER", "test")
	os.Setenv("MYSQL_PASSWORD", "test")
	os.Setenv("MYSQL_DATABASE", "test")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "3306")
	
	// Initialize secret key for JWT
	utils.InitSecretKey()
}

// TeardownTestEnvironment cleans up test environment
func TeardownTestEnvironment() {
	os.Unsetenv("JWT_SECRET")
	os.Unsetenv("MYSQL_USER")
	os.Unsetenv("MYSQL_PASSWORD")
	os.Unsetenv("MYSQL_DATABASE")
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_PORT")
}

// AssertUserExists checks if a user exists in the database
func (td *TestDatabase) AssertUserExists(t *testing.T, email string) *models.User {
	var user models.User
	err := td.DB.Where("email = ?", email).First(&user).Error
	require.NoError(t, err, "User should exist in database")
	return &user
}

// AssertUserNotExists checks if a user does not exist in the database
func (td *TestDatabase) AssertUserNotExists(t *testing.T, email string) {
	var user models.User
	err := td.DB.Where("email = ?", email).First(&user).Error
	require.Error(t, err, "User should not exist in database")
}

// AssertGroupExists checks if a group exists in the database
func (td *TestDatabase) AssertGroupExists(t *testing.T, name string) *models.Group {
	var group models.Group
	err := td.DB.Where("name = ?", name).First(&group).Error
	require.NoError(t, err, "Group should exist in database")
	return &group
}

// AssertGroupNotExists checks if a group does not exist in the database
func (td *TestDatabase) AssertGroupNotExists(t *testing.T, name string) {
	var group models.Group
	err := td.DB.Where("name = ?", name).First(&group).Error
	require.Error(t, err, "Group should not exist in database")
}

// AssertUserInGroup checks if a user is in a group
func (td *TestDatabase) AssertUserInGroup(t *testing.T, userID, groupID uint) {
	var group models.Group
	err := td.DB.Preload("Users").First(&group, groupID).Error
	require.NoError(t, err)
	
	found := false
	for _, user := range group.Users {
		if user.ID == userID {
			found = true
			break
		}
	}
	require.True(t, found, "User should be in group")
}

// AssertUserNotInGroup checks if a user is not in a group
func (td *TestDatabase) AssertUserNotInGroup(t *testing.T, userID, groupID uint) {
	var group models.Group
	err := td.DB.Preload("Users").First(&group, groupID).Error
	require.NoError(t, err)
	
	found := false
	for _, user := range group.Users {
		if user.ID == userID {
			found = true
			break
		}
	}
	require.False(t, found, "User should not be in group")
}

// CountUsers returns the number of users in the database
func (td *TestDatabase) CountUsers() int64 {
	var count int64
	td.DB.Model(&models.User{}).Count(&count)
	return count
}

// CountGroups returns the number of groups in the database
func (td *TestDatabase) CountGroups() int64 {
	var count int64
	td.DB.Model(&models.Group{}).Count(&count)
	return count
}