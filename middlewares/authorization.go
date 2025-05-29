package middlewares

import (
	"felix1234567890/go-trello/models"
	"felix1234567890/go-trello/utils"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

// DeserializeUser is a middleware function that attempts to deserialize a user
// from a JWT token found in the "Authorization" header (Bearer token).
// It requires a *gorm.DB instance to fetch the user from the database.
// If successful, it retrieves the user, attaches the user object (models.User)
// to `c.Locals("user")`, and calls `c.Next()`.
// If token extraction, parsing, validation, or user retrieval fails,
// it returns an appropriate HTTP error response (401 Unauthorized or 403 Forbidden)
// and does not call `c.Next()`.
func DeserializeUser(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var tokenString string
		authorization := c.Get("Authorization")

		if strings.HasPrefix(authorization, "Bearer ") {
			tokenString = strings.TrimPrefix(authorization, "Bearer ")
		}
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": "You are not logged in"})
		}

		tokenByte, err := jwt.Parse(tokenString, func(jwtToken *jwt.Token) (interface{}, error) {
			if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %s", jwtToken.Header["alg"])
			}

			return []byte(utils.SECRET_KEY), nil
		})
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": fmt.Sprintf("invalidate token: %v", err)})
		}

		claims, ok := tokenByte.Claims.(jwt.MapClaims)
		if !ok || !tokenByte.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": "invalid token claim"})

		}

		var user models.User
		db.First(&user, fmt.Sprint(claims["id"]))
		if user.ID == 0 { // Check if user was not found
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "message": "the user belonging to this token no longer exists"})
		}
		// It's better to check if user.ID is 0 (or whatever the zero value for your ID is)
		// instead of comparing with claims["id"] again, as GORM wouldn't populate user.ID if not found.
		// However, the original logic was `user.ID != uint(claims["id"].(float64))`.
		// If `claims["id"]` is guaranteed to be a float64 that can be converted to uint, this might be okay.
		// For robustness, checking if user was found (e.g. user.ID != 0) is often preferred after a First() call.
		// I will keep the original comparison logic for now but add a comment.
		if user.ID != uint(claims["id"].(float64)) { // Original logic
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "message": "the user belonging to this token no longer exists"})
		}


		c.Locals("user", user)

		return c.Next()
	}
}
