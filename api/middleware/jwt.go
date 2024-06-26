package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	e "github.com/nivb52/hotel-rent/api/errors"
)

var SECRET []byte

// function return a string in bytes which use as the secret for the JWT
func GetSecret() []byte {
	if SECRET != nil {
		return SECRET
	}
	secretString := os.Getenv("JWT_SECRET")
	if secretString == "" {
		secretString = "secret-key"
	}
	SECRET = []byte(secretString)
	return SECRET
}

func JWTAuthentication(c *fiber.Ctx) error {
	token, ok := c.GetReqHeaders()["X-Api-Token"]
	if !ok {
		return e.ErrUnAuthorized()
	}

	claims, err := validateToken(token)
	if err != nil {
		return err
	}

	isAdmin := IsAdminNormelize(claims["isAdmin"].(bool))
	c.Context().SetUserValue("userID", claims["id"].(string))
	c.Context().SetUserValue("userEmail", claims["email"].(string))
	c.Context().SetUserValue("isAdmin", isAdmin)

	return c.Next()
}

func validateToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			fmt.Printf("Unexpected signing method: /%v", token.Header["alg"])
			return nil, fmt.Errorf("Invalid")
		}

		secret := GetSecret()
		return secret, nil
	})

	if err != nil {
		fmt.Println("failed to parse JWT token", err)
		return nil, e.ErrUnAuthorized()
	}

	if !token.Valid {
		fmt.Println("invalid token", err)
		return nil, e.ErrUnAuthorized()
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Println("invalid claims", claims)
		return nil, e.ErrUnAuthorized()
	}

	exp := int64(claims["exp"].(float64))
	if time.Now().UTC().Unix() > exp {
		fmt.Printf("\n token expired, current value: %d, token value: %d", time.Now().UTC().Unix(), exp)
		return nil, e.ErrTokenExpired()
	}

	return claims, nil
}

// fn normelize admin data/miss data to boolean
func IsAdminNormelize(isAdmin any) bool {
	if isAdmin == nil {
		return false
	} else {
		isAdmin = isAdmin.(bool)
		if isAdmin == true {
			return true
		} else {
			return false
		}
	}
}
