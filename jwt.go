package wom

import (
	"flag"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var (
	jwtSecret          = flag.String("jwt-secret", "", "Secret to use to validate JWTs")
	subscriptionSecret = flag.String("subscription-secret", "", "Secret used to create subscription JWTs")
)

func emailAddress(token string) (string, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(*jwtSecret), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := t.Claims.(jwt.MapClaims); ok && t.Valid {
		return fmt.Sprintf("%s", claims["email"]), nil
	} else {
		return "", fmt.Errorf("invalid token")
	}
}

func getEmailFromJwt(r *http.Request) string {
	h := r.Header.Get("Authorization")
	if strings.HasPrefix(h, "Bearer") {
		email, _ := emailAddress(strings.TrimPrefix(h, "Bearer "))
		return email
	}
	return ""
}

func createSubscriptionJwt(email string) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "wom",
		Subject:   fmt.Sprintf("subscribe:%s", email),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		NotBefore: jwt.NewNumericDate(time.Now()),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}).SignedString([]byte(*subscriptionSecret))
}

func createUnsubscribeJwt(email string) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:   "wom",
		Subject:  fmt.Sprintf("unsubscribe:%s", email),
		IssuedAt: jwt.NewNumericDate(time.Now()),
	}).SignedString([]byte(*subscriptionSecret))
}

func validateSubscriptionJwt(kind, token string) (string, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(*subscriptionSecret), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := t.Claims.(jwt.MapClaims); ok && t.Valid {
		subject := fmt.Sprintf("%s", claims["sub"])
		if strings.HasPrefix(subject, kind+":") {
			return strings.TrimPrefix(subject, kind+":"), nil
		}
		return "", fmt.Errorf("not a %s token", kind)
	} else {
		return "", fmt.Errorf("invalid token")
	}
}
