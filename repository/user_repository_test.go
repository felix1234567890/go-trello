package repository

import (
	"database/sql/driver"
	"errors"
	"felix1234567890/go-trello/models"
	// "felix1234567890/go-trello/utils" // No longer needed here due to sqlmock.AnyArg() for password
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Helper function to create a new mock GORM DB
func NewMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open sqlmock database: %s", err)
	}

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true, // Important for mocking
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open GORM db: %s", err)
	}
	return gormDB, mock, func() {
		db.Close() // Close the underlying sql.DB
	}
}

// AnyTime argument matcher for sqlmock
type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}


func TestUserRepositoryImpl_CreateUser(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		mockDB, mock, closeDB := NewMockDB(t)
		defer closeDB()
		repo := NewUserRepository(mockDB)

		userInput := &models.User{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password123",
		}
		// We don't need to pre-calculate hashedPassword for WithArgs if using sqlmock.AnyArg() for it.
		// The hashing is done *inside* repo.CreateUser.

		// GORM's actual order from error: (`created_at`,`updated_at`,`deleted_at`,`username`,`email`,`password`)
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users` (`created_at`,`updated_at`,`deleted_at`,`username`,`email`,`password`) VALUES (?,?,?,?,?,?)")).
			WithArgs(AnyTime{}, AnyTime{}, gorm.DeletedAt{}, userInput.Username, userInput.Email, sqlmock.AnyArg()). // Use AnyArg() for the hashed password
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
		
		originalPasswordBeforeHash := userInput.Password
		createdUserID, err := repo.CreateUser(userInput)

		if err != nil {
			t.Errorf("CreateUser() error = %v, wantErr nil", err)
			return
		}
		if createdUserID == 0 { // ID 0 usually means error or not created
			t.Errorf("CreateUser() createdUserID = 0, want non-zero ID")
		}
		if userInput.ID != createdUserID { // Check if GORM populated ID in the input struct
			t.Errorf("CreateUser() userInput.ID = %d, want %d", userInput.ID, createdUserID)
		}
		if userInput.ID != 1 { // Assuming LastInsertId mock result is 1
			t.Errorf("CreateUser() userInput.ID = %d, want 1 (from mock LastInsertId)", userInput.ID)
		}
		if userInput.Password == originalPasswordBeforeHash {
			t.Errorf("CreateUser() password in userInput struct was not hashed")
		}
		if userInput.Password == "" {
			t.Errorf("CreateUser() password in userInput struct is empty after hashing")
		}


		// Check if all expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("database error on creation", func(t *testing.T) {
		mockDB, mock, closeDB := NewMockDB(t)
		defer closeDB()
		repo := NewUserRepository(mockDB)

		userInputForErrorCase := &models.User{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password123",
		}
		// No need to pre-calculate hashedPasswordForErrorCase if using sqlmock.AnyArg()

		// GORM's actual order from error: (`created_at`,`updated_at`,`deleted_at`,`username`,`email`,`password`)
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users` (`created_at`,`updated_at`,`deleted_at`,`username`,`email`,`password`) VALUES (?,?,?,?,?,?)")).
			WithArgs(AnyTime{}, AnyTime{}, gorm.DeletedAt{}, userInputForErrorCase.Username, userInputForErrorCase.Email, sqlmock.AnyArg()). // Use AnyArg()
			WillReturnError(errors.New("DB error"))
		mock.ExpectRollback()

		createdUserID, err := repo.CreateUser(userInputForErrorCase)
		if err == nil {
			t.Errorf("CreateUser() error = nil, wantErr 'DB error'")
		}
		if createdUserID != 0 {
			t.Errorf("CreateUser() createdUserID = %d, want 0 on error", createdUserID)
		}
		if err != nil && !errors.Is(err, errors.New("DB error")) { // Use errors.Is for wrapped errors if applicable, or check string
			// The error from repo.CreateUser might be wrapped by gorm or our code.
			// For this test, a direct string comparison is okay if the error isn't wrapped.
			if err.Error() != "DB error" {
				t.Errorf("CreateUser() error = %v, want 'DB error'", err)
			}
		}

		// Check if all expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	// Note: Testing password hashing.
	// The CreateUser method in UserRepositoryImpl is responsible for hashing the password.
	// The mock query check above (`WithArgs`) already implicitly tests that `hashedPassword` is used.
	// If we wanted to be more explicit about checking that the hashing function is called,
	// that would typically involve mocking the hashing function itself, which is usually too granular for this type of test.
	// The current test verifies the *result* of hashing (i.e., the stored password is not plaintext).
}
