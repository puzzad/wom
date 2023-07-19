package main

import (
	"flag"
	"fmt"
	"github.com/go-resty/resty/v2"
	"log"
	"os"
	"time"

	"github.com/csmith/envflag"
)

var (
	filename = flag.String("filename", "adventures.zip", "path to an adventures zip file")
	apiurl   = flag.String("apiurl", "http://localhost:8090/", "URL of the WOM instance")
	email    = flag.String("email", "", "Admin email address")
	password = flag.String("password", "", "Admin password")
)

func main() {
	envflag.Parse()

	if err := uploadAdventures(*filename, *apiurl, *email, *password); err != nil {
		log.Fatalf("Failed to upload: %v", err)
	}
}

type authRequest struct {
	Identity string `json:"identity"`
	Password string `json:"password"`
}

type authResponse struct {
	Token string `json:"token"`
}

func uploadAdventures(fileName, apiurl string, email string, password string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("unable to open zip: %s", err)
	}
	defer func() {
		_ = file.Close()
	}()

	var response authResponse
	//TODO: Should remove resty really, its overkill now I'm just uploading a file
	client := resty.New()
	client.SetRetryCount(3).
		SetRetryWaitTime(3*time.Second).
		SetRetryMaxWaitTime(10*time.Second).
		SetBaseURL(apiurl).
		SetHeader("Content-Type", "application/json")

	resp, err := client.R().
		SetBody(authRequest{
			Identity: email,
			Password: password,
		}).
		SetResult(&response).
		Post("api/admins/auth-with-password")
	if err != nil {
		return fmt.Errorf("unable to login")
	}
	if resp.StatusCode() != 200 {
		return fmt.Errorf("unable to login: Status %d", resp.StatusCode())
	}
	client.SetAuthToken(response.Token)
	resp, err = client.R().
		SetContentLength(true).
		SetFileReader("adventures.zip", "adventures.zip", file).
		Post("/wom/importzip")
	if err != nil {
		return fmt.Errorf("unable to upload zip: %s", err)
	}
	if resp.StatusCode() != 200 {
		return fmt.Errorf("unable to upload zip: Status %d", resp.StatusCode())
	}
	return nil
}
