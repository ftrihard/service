// Package plugindemo a demo plugin.
package service1

import (
	"time"
	"context"
	"net/http"
	"encoding/json"
	jwt "github.com/golang-jwt/jwt/v4"
)

const APPLICATION_NAME = "Service One"
const LOGIN_EXPIRATION_DURATION = time.Duration(1) * time.Hour
const JWT_SIGNATURE_KEY = "the secret of kalimdor"

// Config the plugin configuration.
type Config struct {}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{}
}

// a encoder plugin.
type Encoder struct {
	next     http.Handler
	name     string
}

type Claim struct {
	jwt.StandardClaims
	UserId string
}

// New created a new encoder plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {

	return &Encoder{
		next:     next,
		name:     name,
	}, nil
}

func (e *Encoder) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	claim := Claim {
		StandardClaims: jwt.StandardClaims{
			Issuer:    APPLICATION_NAME,
			ExpiresAt: time.Now().Add(LOGIN_EXPIRATION_DURATION).Unix(),
		},
		UserId: req.Header.Get("User-Id"),
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256,claim).SignedString([]byte(JWT_SIGNATURE_KEY))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	req.Header.Set("User-Id",token)

	e.next.ServeHTTP(rw, req)
}
