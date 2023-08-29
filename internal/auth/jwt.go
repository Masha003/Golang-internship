package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(userId string, lifespan time.Duration, secret string) (string, string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = userId
	claims["exp"] = time.Now().Add(lifespan).Unix()
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := t.SignedString([]byte(secret))
	if err != nil {
		return "", "", err
	}

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["user_id"] = userId
	rtClaims["exp"] = time.Now().Add(lifespan).Unix()
	rt, err := refreshToken.SignedString([]byte("secret"))
	if err != nil {
		return "", "", err
	}

	return token, rt, nil
}

func Validate(tokenString string, secret string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

func ExtractId(token *jwt.Token) (string, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		uid, ok := claims["user_id"].(string)
		if !ok {
			return "0", fmt.Errorf("invalid user_id: %v", claims["user_id"])
		}

		return uid, nil
	}

	return "0", nil
}
