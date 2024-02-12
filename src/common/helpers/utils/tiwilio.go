package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func SendOtpViaMobile(mobileNumber string, otpCode int) error {
	// Set account keys & information
	accountSid := os.Getenv("SID")
	authToken := os.Getenv("AUTH_TOKEN")
	from := os.Getenv("TWILIO_NO")

	// Correctly encode the OTP in the message body
	body := fmt.Sprintf("Your OTP is: %d", otpCode)

	urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/Messages.json"

	// Use url.Values for form data
	msgData := url.Values{}
	msgData.Set("To", mobileNumber)
	msgData.Set("From", from)
	msgData.Set("Body", body)
	msgDataReader := strings.NewReader(msgData.Encode())

	client := &http.Client{}
	req, err := http.NewRequest("POST", urlStr, msgDataReader)
	if err != nil {
		log.Println("Error creating request: ", err)
		return err
	}

	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending request: ", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var data map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&data)
		if err == nil {
			fmt.Println(data["sid"])
		}
		log.Println("------------Sent SMS successfully------------")
	} else {
		log.Println("Failed to send SMS, status code: ", resp.Status)
		return fmt.Errorf("failed to send SMS, status code: %s", resp.Status)
	}

	return nil
}
