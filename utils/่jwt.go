package utils

import (
	"errors"
	"fmt"
	"os"
	"github.com/golang-jwt/jwt/v4"
)

type (
	Claims struct {
		Id int `json:"id"`	
	}

	AuthMapClaims struct {
		*Claims
		jwt.RegisteredClaims
	}
)

func ParseToken(secertKey, tokenString string) (*AuthMapClaims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &AuthMapClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("error: unexpected signing method")
		}
		return []byte(secertKey), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, errors.New("error: token format is invalid")
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("error: token is expired")
		} else if errors.Is(err, jwt.ErrSignatureInvalid) {
			return nil, errors.New("error: token signature is invalid")
		} else {
			return nil, errors.New("error: token is invalid")
		}
	}

	if claims, ok := token.Claims.(*AuthMapClaims); ok {
		return claims, nil
	} else {
		return nil, errors.New("error: claims type is invalid")
	}
}

func ParseApiKey(apisecret, tokenString string) (*AuthMapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AuthMapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signing method is invalid")
		}
		return []byte(apisecret), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, fmt.Errorf("token format is invalid")
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token had expired")
		} else {
			return nil, fmt.Errorf("parse token failed: %v", err)
		}
	}

	if claims, ok := token.Claims.(*AuthMapClaims); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("claims type is invalid")
	}
}

func ValidApikey(accessToken string) error {

	if accessToken != os.Getenv("API_KEY") {
		return errors.New("error: access token is invalid")
	}

	return nil
}
