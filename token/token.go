package token

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte(os.Getenv("TOKEN_SECRET"))

type UserClaim struct {
	jwt.RegisteredClaims
	AccountUUID string
	Username    string
	Exp         int64
}

// CreateToken generates a JWT token with the provided username and UUID.
// The token is signed using the HS256 signing method and contains the following claims:
// - "username": the username of the account
// - "accountUuid": the UUID of the account
// - "exp": the expiration time of the token, set to 24 hours from the current time
// The function returns the generated token string and an error if any occurred.
func CreateToken(username string, uuid string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaim{
		RegisteredClaims: jwt.RegisteredClaims{},
		AccountUUID:      uuid,
		Username:         username,
		Exp:              time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// VerifyToken verifies the validity of a JWT token.
// It takes a token string as input and returns an error if the token is invalid.
func VerifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}

func ParseToken(tokenString string) (*UserClaim, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaim{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*UserClaim)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
