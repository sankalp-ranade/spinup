package config

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/golang-jwt/jwt"
)

const (
	DefaultNetwork = string("spinup_services")
)

type Configuration struct {
	Common struct {
		Architecture string `yaml:"architecture"`
		ProjectDir   string `yaml:"projectDir"`
		Ports        []int  `yaml:"ports"`
		ClientID     string `yaml:"client_id"`
		ClientSecret string `yaml:"client_secret"`
		ApiKey       string `yaml:"api_key"`
	} `yaml:"common"`
	VerifyKey *rsa.PublicKey
	SignKey   *rsa.PrivateKey
	UserID    string
}

var Cfg Configuration

func ValidateUser(authHeader string, apiKeyHeader string) (string, error) {
	if authHeader == "" && apiKeyHeader == "" {
		return "", errors.New("no authorization keys found")
	}
	if apiKeyHeader == "" {
		userId, err := ValidateToken(authHeader)
		if err != nil {
			log.Printf("error validating token %v", authHeader)
			return "", errors.New("error validating token")
		}
		return userId, nil
	}

	if authHeader == "" {
		err := ValidateApiKey(apiKeyHeader)
		if err != nil {
			log.Printf("error validating api-key %v", apiKeyHeader)
			return "", errors.New("error validating api-key")
		}
	}
	return "testuser", nil
}

func ValidateApiKey(apiKeyHeader string) error {
	if apiKeyHeader != Cfg.Common.ApiKey {
		return errors.New("invalid api key")
	}
	return nil
}

func ValidateToken(authHeader string) (string, error) {
	splitToken := strings.Split(authHeader, "Bearer ")
	if len(splitToken) < 2 {
		return "", fmt.Errorf("cannot validate empty token")
	}
	reqToken := splitToken[1]
	userID, err := JWTToString(reqToken)
	if err != nil {
		return "", err
	}

	if userID == "" {
		return "", errors.New("user ID cannot be blank")
	}
	return userID, nil
}

func JWTToString(tokenString string) (string, error) {
	keyFunc := func(t *jwt.Token) (interface{}, error) {
		return Cfg.VerifyKey, nil
	}
	claims := &Claims{}
	log.Println("JWT to string:", tokenString)
	token, err := jwt.ParseWithClaims(tokenString, claims, keyFunc)
	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", errors.New("invalid token")
	}
	log.Println("claims:", claims.Text)
	return claims.Text, nil
}

// Create a struct that will be encoded to a JWT.
// We add jwt.StandardClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	Text string `json:"text"`
	jwt.StandardClaims
}
