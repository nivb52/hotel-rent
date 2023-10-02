package middleware

import (
	"fmt"
	"go/token"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthentication(c *fiber.Ctx) error {
	token, ok := c.GetReqHeaders()["X-Api-Token"]
	if !ok {
		return fmt.Errorf("Unauthorized")
	}

	err :=parseJWTToken(token)
	if err != nil {
		return err
	}

	fmt.Println("token", token)
	return c.Next()
}


func parseJWTToken(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token)) (any, error) {
	 	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("Unexpected signing method: /%v",token.Header["X-Api-Token"])
			return mil, fmt.Errorf("Invalid")
		}

		secret := os.Getenv("JWT_SECRET")
		return []byte( secret), nil
	}

	if if err != nil {
		fmt.Println("failed to parse JWT token", err)
		return nil, fmt.Errorf("Unauthorized")
	}
	
	

	 if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims)
	 }

	return token, err
}