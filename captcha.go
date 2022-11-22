package wom

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

var (
	hcaptchaSiteKey   = flag.String("hcaptcha-site-key", "", "Site key to use for hCaptcha")
	hcaptchaSecretKey = flag.String("hcaptcha-secret-key", "", "Secret key to use for hCaptcha")
)

func checkCaptcha(token string) error {
	type Response struct {
		Success bool `json:"success"`
	}

	values := url.Values{
		"secret":   {*hcaptchaSecretKey},
		"sitekey":  {*hcaptchaSiteKey},
		"response": {token},
	}

	resp, err := http.PostForm("https://hcaptcha.com/siteverify", values)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
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
