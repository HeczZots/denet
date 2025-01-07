package middleware

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthMiddleware struct {
	JWTSecret []byte
}

type Claims struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

func NewAuthMiddleware(jwtSecret []byte) *AuthMiddleware {
	return &AuthMiddleware{JWTSecret: jwtSecret}
}

func (mw *AuthMiddleware) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/login", "/registration":
			next.ServeHTTP(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		err := mw.validateToken(authHeader)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (mw *AuthMiddleware) GenerateJWT(username string) (tokenString string, err error) {
	c := Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "denet",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	tokenString, err = token.SignedString(mw.JWTSecret)
	if err != nil {
		log.Fatal("Failed to sign token:", err)
	}
	slog.Info("Generated ", "Token", tokenString, "User", username)

	return tokenString, nil
}

func (mw *AuthMiddleware) validateToken(signedToken string) error {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return mw.JWTSecret, nil
		},
	)
	if err != nil {
		return err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		err = fmt.Errorf("couldn't parse claims")
		return err
	}

	if time.Now().After(claims.ExpiresAt.Time) {
		return fmt.Errorf("token expired")
	}

	return nil
}
