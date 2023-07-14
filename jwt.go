package wom

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func createSubscriptionJwt(secret, email string) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "wom",
		Subject:   fmt.Sprintf("subscribe:%s", email),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		NotBefore: jwt.NewNumericDate(time.Now()),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}).SignedString([]byte(secret))
}

func createUnsubscribeJwt(secret, email string) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:   "wom",
		Subject:  fmt.Sprintf("unsubscribe:%s", email),
		IssuedAt: jwt.NewNumericDate(time.Now()),
	}).SignedString([]byte(secret))
}

func validateSubscriptionJWT(secret, kind, token string) (string, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
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

func checkCaptcha(siteKey, secretKey, token string) error {
	type Response struct {
		Success bool `json:"success"`
	}

	values := url.Values{
		"secret":   {secretKey},
		"sitekey":  {siteKey},
		"response": {token},
	}

	resp, err := http.PostForm("https://hcaptcha.com/siteverify", values)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		return err
	}

	var response = Response{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}

	if response.Success {
		return nil
	} else {
		return fmt.Errorf("captcha verification failed")
	}
}
