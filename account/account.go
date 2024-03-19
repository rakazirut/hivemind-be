package account

import (
	"example/hivemind-be/db"
	"example/hivemind-be/token"
	"net/http"
	"net/mail"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Account struct {
	ID       int32       `json:"Id" gorm:"primaryKey:type:int32"`
	Username string      `json:"Username"`
	Email    string      `json:"Email"`
	Password string      `json:"Password"`
	UUID     string      `json:"Uuid"`
	Deleted  bool        `json:"Deleted"`
	Banned   bool        `json:"Banned"`
	Created  pq.NullTime `json:"Created"`
}

func CreateAccount(c *gin.Context) {
	var acc Account

	if err := c.BindJSON(&acc); err != nil {
		return
	}

	addr, err := mail.ParseAddress(strings.ToLower(acc.Email))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Error: Email address format is not valid. Please us a valid email address.",
		})
		return
	}

	hashedPassword, err := hashPassword(acc.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Error: A error occurred creating the account. Please try again.",
		})
		return
	}

	acc.Username = strings.ToLower(acc.Username)
	acc.Email = strings.ToLower(addr.Address)
	acc.Password = hashedPassword
	acc.UUID = uuid.NewString()
	acc.Deleted = false
	acc.Banned = false
	acc.Created = pq.NullTime{Time: time.Now(), Valid: true}

	if result := db.Db.Create(&acc); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Could not create account. Please try again.",
		})
		return
	}

	c.JSON(http.StatusCreated, acc)
}

// Hash password
func hashPassword(password string) (string, error) {
	// Convert password string to byte slice
	var passwordBytes = []byte(password)

	// Hash password with Bcrypt's min cost
	hashedPasswordBytes, err := bcrypt.
		GenerateFromPassword(passwordBytes, bcrypt.MinCost)

	return string(hashedPasswordBytes), err
}

// Check if two passwords match using Bcrypt's CompareHashAndPassword
// which return nil on success and an error on failure.
func doPasswordsMatch(hashedPassword, currPassword string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword), []byte(currPassword))
	return err == nil
}

func AccountLogin(c *gin.Context) {
	var acc Account
	var account Account

	if err := c.BindJSON(&acc); err != nil {
		return
	}

	if result := db.Db.Where("username = ?", acc.Username).First(&account); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Account not found. Please try again.",
		})
		return
	}

	if !doPasswordsMatch(account.Password, acc.Password) {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Password is incorrect. Please try again.",
		})
		return
	}

	token, err := token.CreateToken(account.Username, account.UUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Error: A error occurred creating the token. Please try again.",
		})
		return
	}

	c.SetCookie("Token", token, 86400, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"Token": token,
	})
}

func ValidateAccountToken(c *gin.Context) {
	cookie, err := c.Request.Cookie("Token")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "No token found in request.",
		})
		return
	}

	tokenString := cookie.Value
	if err := token.VerifyToken(tokenString); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Invalid token. Please try again.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Message": "Token is valid.",
	})
}
