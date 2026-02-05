package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func getSecretKey() string {
	secretKey := os.Getenv("JWT_SECRET") 
	if secretKey == "" {
		panic("JWT_SECRET is not set in environment")
	}
	return secretKey
}

// GenerateToken creates a JWT for the user
func GenerateToken(email, userId string) (string, error) {
	claims := jwt.MapClaims{
		"email":  email,
		"userId": userId, 
		"exp":    time.Now().Add(24 * time.Hour).Unix(), // 24 hours expiry
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(getSecretKey()))
	if err != nil {
		return "", errors.New("failed to generate token: " + err.Error())
	}

	return signedToken, nil
}

// VerifyToken validates JWT and returns userId (UUID string)
func VerifyToken(tokenStr string) (string, error) { // ✅ Return string, not int64
	parsedToken, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(getSecretKey()), nil
	})

	if err != nil {
		return "", errors.New("could not parse token: " + err.Error())
	}

	if !parsedToken.Valid {
		return "", errors.New("invalid token")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	// Check expiration (jwt.Parse already validates this, but explicit check is good practice)
	exp, ok := claims["exp"].(float64)
	if !ok {
		return "", errors.New("token missing expiration claim")
	}
	if int64(exp) < time.Now().Unix() {
		return "", errors.New("token has expired")
	}

	// Extract userId as string (UUID)
	userId, ok := claims["userId"].(string) 
	if !ok {
		return "", errors.New("token missing or invalid userId")
	}

	return userId, nil
}