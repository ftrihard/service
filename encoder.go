package service

import (
	"fmt"
	"time"
	"context"
	"net/http"
	jwt "github.com/golang-jwt/jwt/v4"
)

const APPLICATION_NAME = "JWT Service"
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
	UserId string `json:"UserId"`
}

// New created a new encoder plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {

	return &Encoder{
		next:     next,
		name:     name,
	}, nil
}

func (e *Encoder) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if req.Header.Get("X-Jwt") == "" {
		claim := Claim {
			StandardClaims: jwt.StandardClaims{
				Issuer:    APPLICATION_NAME,
			},
			UserId: req.Header.Get("X-User-Id"),
		}
		req.Header.Del("X-User-Id")
		token, err := jwt.NewWithClaims(jwt.SigningMethodHS256,claim).SignedString([]byte(JWT_SIGNATURE_KEY))
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
		req.Header.Add("X-Jwt",token)
	} else {
		token, err := jwt.Parse(req.Header.Get("X-Jwt"), func(token *jwt.Token) (interface{}, error) {
			if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Signing method invalid")
			} else if method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("Signing method invalid")
			}
		
			return []byte(JWT_SIGNATURE_KEY), nil
		})
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
		req.Header.Del("X-Jwt")
		req.Header.Add("X-User-Id",claims["UserId"].(string))
	}
	

	e.next.ServeHTTP(rw, req)
}
