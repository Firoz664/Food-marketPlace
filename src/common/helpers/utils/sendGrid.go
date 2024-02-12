package utils

import (
	"fmt"
)

func SendOtpViaEmail(mobileNumber string, otpCode int) error {
	// Implement the logic to send the OTP to the user's mobile number/email
	fmt.Printf("Sending OTP %d to %s\n", otpCode, mobileNumber)
	return nil
}

// func SendOtpViaEmail(toEmail string, otpCode int) error {
// 	// from := mail.NewEmail("Virality Media", "your-email@example.com") // Sender email
// 	// subject := "Your OTP"
// 	// to := mail.NewEmail("Recipient Name", toEmail)
// 	// plainTextContent := fmt.Sprintf("Your OTP is: %d", otpCode)
// 	// htmlContent := fmt.Sprintf("<strong>Your OTP is: %d</strong>", otpCode)
// 	// message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

// 	// client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
// 	// response, err := client.Send(message)
// 	if err != nil {
// 		fmt.Println(err)
// 		return err
// 	} else {
// 		fmt.Printf("Email sent, status code: %d\n", response.StatusCode)
// 	}

// 	return nil
// }
