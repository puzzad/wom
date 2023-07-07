package wom

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/pocketbase/pocketbase"
	"github.com/spf13/cobra"
	"os"
	"time"
)

func NewImportCmd(app *pocketbase.PocketBase) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import",
		Short: "Imports adventures",
		Run: func(cmd *cobra.Command, args []string) {
			filename, _ := cmd.Flags().GetString("filename")
			apiurl, _ := cmd.Flags().GetString("apiurl")
			email, _ := cmd.Flags().GetString("email")
			password, _ := cmd.Flags().GetString("password")
			cobra.CheckErr(uploadAdventures(filename, apiurl, email, password))
		},
	}
	cmd.Flags().String("filename", "adventures.zip", "path to an adventures zip file")
	cmd.Flags().Bool("production", false, "Whether this is an upload to production or dev")
	cmd.Flags().String("apiurl", "http://localhost:8090/", "URL of the WOM instance")
	cmd.Flags().String("email", "", "Admin email address")
	cmd.Flags().String("password", "", "Admin password")
	_ = cmd.MarkFlagRequired("email")
	_ = cmd.MarkFlagRequired("password")
	return cmd
}

func uploadAdventures(fileName, apiurl string, email string, password string) error {
	data := struct {
		Token string
	}{}
	file, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("unable to open zip: %s", err)
	}
	defer func() {
		_ = file.Close()
	}()
	//TODO: Should remove resty really, its overkill now I'm just uploading a file
	client := resty.New()
	client.SetRetryCount(3).
		SetRetryWaitTime(3*time.Second).
		SetRetryMaxWaitTime(10*time.Second).
		SetBaseURL(apiurl).
		SetHeader("Content-Type", "application/json")
	resp, err := client.R().
		SetBody(`{"identity":"` + email + `", "password":"` + password + `"}`).
		SetResult(&data).
		Post("api/admins/auth-with-password")
	if err != nil {
		return fmt.Errorf("unable to login")
	}
	if resp.StatusCode() != 200 {
		return fmt.Errorf("unable to login: Status %d", resp.StatusCode())
	}
	client.SetAuthToken(data.Token)
	resp, err = client.R().
		SetContentLength(true).
		SetFileReader("adventures.zip", "adventures.zip", file).
		Post("/import/zip")
	if err != nil {
		return fmt.Errorf("unable to upload zip: %s", err)
	}
	if resp.StatusCode() != 200 {
		return fmt.Errorf("unable to upload zip: Status %d", resp.StatusCode())
	}
	return nil
}
