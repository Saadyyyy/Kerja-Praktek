package emails

import (
	"os"
	"strconv"

	"github.com/go-gomail/gomail"
)

func SendLoginNotification(userEmail string, name string) error {
	// Mengambil nilai dari environment variables
	smtpServer := os.Getenv("SMTP_SERVER")
	smtpPortStr := os.Getenv("SMTP_PORT")
	smtpUsername := os.Getenv("SMTP_USERNAME")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	// Konfigurasi pengiriman email
	sender := smtpUsername
	recipient := userEmail
	subject := "Successful Login Notification"
	emailBody := `
    <html>
    <head>
        <style>
            body {
                font-family: Arial, sans-serif;
                background-color: #f5f5f5;
            }p {
    font-size: 18px;
    margin-top: 15px;
    line-height: 1.6;
}

h1 {
    font-size: 24px;
}
a {
    text-decoration: none;
    color: #007BFF;
    font-weight: bold;
}

a:hover {
    color: #0056b3;
}
h1 {
    text-align: center;
    color: #007BFF;
}
.footer {
    text-align: center;
    margin-top: 20px;
    color: #666;
    padding: 10px;
    background-color: #f5f5f5;
    border-top: 1px solid #ddd;
}

        </style>
    </head>
    <body>
        <div class="container">
            <h1>Login Successful</h1>
            <div class="message">
                <p>Hello, <strong>` + name + `</strong>,</p>
                <p>You have successfully logged into your account.</p>
                <p>If this was not you, please contact our support team immediately, thank you.</p>
                <p><strong>Support Team:</strong> <a href="laodesaady12345@gmail.com">Laodesaady12345@gmail.com</a></p>
            </div>
            <div class="footer">
                <p>&copy; 2023 laode saady. All rights reserved.</p>
            </div>
        </div>
    </body>
    </html>
    `

	// Convert the SMTP port from string to integer
	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", sender)
	m.SetHeader("To", recipient)
	m.SetHeader("Subject", subject)
	// Set pesan dalam format HTML
	m.SetBody("text/html", emailBody)

	d := gomail.NewDialer(smtpServer, smtpPort, smtpUsername, smtpPassword)

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
