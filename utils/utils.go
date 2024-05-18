package utils

import (
	"example/hivemind-be/token"
	"net/http"
	"reflect"
	"regexp"
	"unicode"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// JsonDataHasKey checks if the given data has a field with the specified key.
// It returns the value of the field and a boolean indicating whether the field exists.
func JsonDataHasKey(data interface{}, key string) (interface{}, bool) {
	value := reflect.ValueOf(data)

	// Check if the data is a struct
	if value.Kind() != reflect.Struct {
		return "null", false
	}

	// Get the field by name
	fieldValue := value.FieldByName(key)

	// Check if the field exists
	if !fieldValue.IsValid() {
		return "null", false
	}

	// Return the field value
	return fieldValue.Interface(), true
}

func ValidateAuthentication(c *gin.Context, authToken string) (*token.UserClaim, bool) {
	validToken := token.CheckToken(c, authToken)
	expire := token.CheckTokenNotExpired(c, authToken)

	if !validToken {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Error": "Unauthorized.",
		})
		return nil, false
	}
	if !expire {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Error": "Token has expired.",
		})
		return nil, false
	}
	claims, err := token.ParseToken(authToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Error": "Unauthorized.",
		})
		return nil, false
	}
	return claims, true
}

func ValidateRefreshAuthentication(c *gin.Context, refreshToken string) (*token.RefreshUserClaim, bool) {
	validToken := token.CheckRefreshToken(c, refreshToken)
	expire := token.CheckRefreshTokenNotExpired(c, refreshToken)

	if !validToken {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Error": "Unauthorized.",
		})
		return nil, false
	}
	if !expire {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Error": "Token has expired.",
		})
		return nil, false
	}
	claims, err := token.ParseRefreshToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Error": "Unauthorized.",
		})
		return nil, false
	}
	return claims, true
}

func ValidateCommentMessage(message string) bool {
	if len(message) >= 1 && len(message) <= 2048 {
		return true
	} else {
		return false
	}
}

func ValidateContentTitle(title string) bool {
	if len(title) >= 1 && len(title) <= 256 {
		return true
	} else {
		return false
	}
}

func ValidateContentMessage(message string) bool {
	if len(message) >= 1 && len(message) <= 5000 {
		return true
	} else {
		return false
	}
}

func ValidateHiveName(name string) bool {
	namePattern := "^[a-zA-Z]{1,30}$"
	nameRegex, err := regexp.Compile(namePattern)
	if err != nil {
		return false
	}
	// Check if the test string matches the pattern
	if !nameRegex.MatchString(name) {
		return false
	}
	return true
}

func ValidateHiveDescription(description string) bool {
	if len(description) >= 1 && len(description) <= 256 {
		return true
	} else {
		return false
	}
}

// Hash password
func HashPassword(password string) (string, error) {
	// Convert password string to byte slice
	var passwordBytes = []byte(password)

	// Hash password with Bcrypt's min cost
	hashedPasswordBytes, err := bcrypt.
		GenerateFromPassword(passwordBytes, bcrypt.MinCost)

	return string(hashedPasswordBytes), err
}

// Check if two passwords match using Bcrypt's CompareHashAndPassword
// which return nil on success and an error on failure.
func DoPasswordsMatch(hashedPassword, currPassword string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword), []byte(currPassword))
	return err == nil
}

// ValidatePasswordComplexity checks if a password meets the required complexity criteria.
// It returns true if the password meets the criteria, and false otherwise.
func ValidatePasswordComplexity(s string) bool {
	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)
	if len(s) >= 12 {
		hasMinLen = true
	}
	for _, char := range s {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}
